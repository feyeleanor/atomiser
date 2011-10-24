package atomiser

import(
	"fmt"
	. "github.com/feyeleanor/chain"
	"github.com/feyeleanor/slices"
	"scanner"
	"strconv"
)

type Delimiter	int

type Reader		func(Scanner) interface{}

func ReadPair(s Scanner, read Reader) (r interface{}) {
	r = read(s)
	if _, ok := r.(Delimiter); ok || r == nil {
		panic("missing item after .")
	}
	obj := read(s)
	if _, ok := obj.(Delimiter); !ok && obj != nil {
		panic("extra item after .")
	}
	return
}

func ReadString(s Scanner, read Reader) (r string) {
	for c := s.Next(); c != '"'; c = s.Next() {
		if c == scanner.EOF {
			panic("EOF in string literal")
		}
		r = fmt.Sprintf("%v%c", r, c)
	}
	return r
}

func ReadList(s Scanner, read Reader) (head *Cell) {
	head = new(Cell)
	for tail := head; tail != nil ; tail = tail.Tail {
		switch obj := read(s); {
		case obj == Delimiter(')'):	break
		case obj == Delimiter('.'):	tail.Head = obj
									tail.Tail = &Cell{ Head: ReadPair(s, read) }
									tail = tail.Tail
									break
		default:					tail.Head = obj
		}
	}
	return
}

func ReadArray(s Scanner, read Reader) (array slices.Slice) {
	for obj := read(s); obj != Delimiter(']') ; obj = read(s) {
		if obj == Delimiter('.') {
			obj = ReadPair(s, read)
			array = append(array, obj)
			break
		} else {
			array = append(array, obj)
		}
	}
	return
}

func ReadNumber(s Scanner) (i interface{}) {
	switch c := s.Peek(); c {
	case '#':		defer func() {
						if recover() != nil {
							panic("Invalid radix format")
						}
					}()
					s.Next()
					switch radix := int(s.ReadInteger(10)); s.Peek() {
					case 'r', 'R':			s.Next()
											i = s.ReadInteger(radix)
					default:				panic("R missing from radix-format integer")
					}

	case '0':		s.Next()
					switch c = s.Peek(); c {
					case '.':				i, _ = strconv.Atof64("0" + s.ReadDecimalPlaces())

					case 'x', 'X':			s.Next()
											i = s.ReadInteger(16)
					case 'b', 'B':			s.Next()
											i = s.ReadInteger(2)

					default:				defer func() {
												if recover() != nil {
													i = int64(0)
												}
											}()
											i = s.ReadInteger(8)
					}

	default:		if d := s.ReadDigits(10); s.Peek() == '.' {
						i, _ = strconv.Atof64(d + s.ReadDecimalPlaces())
					} else {
						i, _ = strconv.Btoi64(d, 10)
					}
	}
	return
}