package atomiser

import(
	"math/big"
	"fmt"
	"io"
	"strconv"
)

type Lisper struct {
	*Atomiser
}

func NewLisper(f io.Reader) (l *Lisper) {
	return &Lisper{ NewAtomiser(f) }
}

func (l Lisper) ReadDigits(radix int) (r string) {
	var s	[]rune
	for ; l.IsRadix(radix); l.Next() {
		s = append(s, l.Peek())
	}
	r = string(s)
	return
}

func (l Lisper) ReadInteger(radix int) int64 {
	d := l.ReadDigits(radix)
	if i, e := strconv.ParseInt(d, radix, 64); e == nil {
		return i
	}
	i := big.NewInt(0)
	if _, ok := i.SetString(d, radix); !ok {
		panic(fmt.Sprintf("Invalid number: %v (\"%v\") for base: %v", i, d, radix))
	}
	return 0
}

func (l Lisper) ReadDecimalPlaces(r int) string {
	m := []rune{}
	e := []rune{}
	l.Next()
	for parsed := false; !parsed && !l.IsEOF() ; l.Next() {
		switch c := l.Peek(); c {
		case 'e', 'E':		if len(e) == 0 {
								e = append(e, 'E')
							} else {
								parsed = true
							}
		case '+', '-':		if len(e) == 1 {
								e = append(e, c)
							} else {
								parsed = true
							}
		default:			if l.IsRadix(r) {
								switch len(e) {
								case 0:		m = append(m, c)

								case 1:		e = append(e, '+', c)

								default:	e = append(e, c)
								}
							} else {
								parsed = true
							}
		}
	}
	return string(append(m, e...))
}

func (l Lisper) ReadSymbol() (i interface{}) {
	switch c := l.Peek(); c {
	case '#':		l.Next()
					switch radix := int(l.ReadInteger(10)); l.Peek() {
					case 'r', 'R':			l.Next()
											i = l.ReadInteger(radix)

					default:				panic("R missing from radix-format integer")
					}

	case '0':		l.Next()
					switch c = l.Peek(); c {
					case '.':				i, _ = strconv.ParseFloat("0."+l.ReadDecimalPlaces(10), 64)

					case 'x', 'X':			l.Next()
											i = l.ReadInteger(16)

					case 'b', 'B':			l.Next()
											i = l.ReadInteger(2)

					default:				if l.IsRadix(8) {
												i = l.ReadInteger(8)
											} else {
												i = int64(0)
											}
					}

	default:		if l.IsRadix(10) {
						if d := l.ReadDigits(10); l.Peek() == '.' {
							i, _ = strconv.ParseFloat(d+"."+l.ReadDecimalPlaces(10), 64)
						} else {
							i, _ = strconv.ParseInt(d, 10, 64)
						}
					} else {
						i = l.Atomiser.ReadSymbol()
					}
	}
	return
}