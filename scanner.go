package atomiser

import (
	"fmt"
	"io"
	"scanner"
	"strconv"
)

type Scanner struct {
	*scanner.Scanner
}

func NewScanner(f io.Reader) (s Scanner) {
	s = Scanner{new(scanner.Scanner)}
	s.Init(f)
	s.Mode &^= scanner.GoTokens
	return
}

func (s Scanner) IsEOF() bool {
	return s.Peek() == scanner.EOF
}

func (s Scanner) IsLineBreak() bool {
	c := s.Peek()
	return c == '\n' || c == '\r'
}

func (s Scanner) IsWhitespace() bool {
	c := s.Peek()
	return c == ' ' || s.IsLineBreak()
}

func (s Scanner) IsDelimiter(d Delimiter) bool {
	return d == Delimiter(s.Peek())
}

func (s Scanner) IsPrint() bool {
	c := s.Peek()
	return ' ' <= c && c <= '~'
}

func (s Scanner) IsAlpha() bool {
	c := s.Peek()
	return 'A' <= c && c <= 'Z' || 'a' <= c && c <= 'z'
}

func (s Scanner) IsRadix(r int) (ok bool) {
	c := s.Peek()
	switch {
	case r <= 10:	ok = '0' <= c && c < ('0' + r)
	case r <= 36:	if ok = '0' <= c && c <= '9'; !ok {
						r = r - 11
						ok = 'A' <= c && c <= ('A' + r) || 'a' <= c && c <= ('a' + r)
					}
	}
	return
}

func (s Scanner) DigitValue() (r int) {
	switch c := s.Next(); {
	case '0' <= c && c <= '9':		r = c - '0'
	case 'A' <= c && c <= 'Z':		r = c - 'A' + 10
	case 'a' <= c && c <= 'z':		r = c - 'a' + 10
	default:						panic("illegal digit in character escape")
	}
	return
}

func (s Scanner) ReadChar() (r int) {
	if c := s.Next(); c != '\\' {
		r = c
	} else {
		switch c = s.Next(); c {
		case 'a':			r = '\a'
		case 'b':			r = '\b'
		case 'f':			r = '\f'
		case 'n':			r = '\n'
		case 'r':			r = '\r'
		case 't':			r = '\t'
		case 'v':			r = '\v'
		case '\'':			r = '\''

		case 'u':			r = (s.DigitValue() << 24) + (s.DigitValue() << 16) + (s.DigitValue() << 8) + s.DigitValue()

		case 'x':			if s.IsRadix(16) {
								if r = s.DigitValue(); s.IsRadix(16) {
									r = r * 16 + s.DigitValue()
								}
							}

		case '0':			fallthrough
		case '1':			fallthrough
		case '2':			fallthrough
		case '3':			fallthrough
		case '4':			fallthrough
		case '5':			fallthrough
		case '6':			fallthrough
		case '7':			r = s.DigitValue()
							if s.IsRadix(8) {
								if r = r * 8 + s.DigitValue(); s.IsRadix(8) {
									r = r * 8 + s.DigitValue()
								}
							}

		default:			if s.IsAlpha() || s.IsRadix(10) {
								panic(fmt.Sprintf("illegal character escape: \\%c", c))
							}
							r = c
		}
	}
	return
}

func (s Scanner) ReadDigits(radix int) (r string) {
	for ; s.IsRadix(radix); s.Next() {
		r = fmt.Sprintf("%v%c", r, s.Peek())
	}
	if len(r) == 0 {
		panic(fmt.Sprintf("Invalid number: %c for base: %v", s.Peek(), radix))
	}
	return
}

func (s Scanner) ReadInteger(radix int) (i int64) {
	i, _ = strconv.Btoi64(s.ReadDigits(radix), radix)
	return
}

func (s Scanner) ReadDecimalPlaces() (r string) {
	r = "."
	e := ""
	s.Next()
	for parsed := false; !parsed && !s.IsEOF() ; s.Next() {
		switch c := s.Peek(); c {
		case 'e', 'E':		if len(e) == 0 {
								e = "E"
							} else {
								parsed = true
							}
		case '+', '-':		if len(e) == 1 {
								e = fmt.Sprintf("%v%c", e, c)
							} else {
								parsed = true
							}
		default:			if s.IsRadix(10) {
								switch len(e) {
								case 0:		r = fmt.Sprintf("%v%c", r, c)
								case 1:		e = fmt.Sprintf("E+%c", c)
								default:	e = fmt.Sprintf("%v%c", e, c)
								}
							} else {
								parsed = true
							}
		}
	}
	r += e
	return
}

func (s Scanner) ReadNumber() (i interface{}) {
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