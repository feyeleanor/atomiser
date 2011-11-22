package atomiser

import (
	"fmt"
	"github.com/feyeleanor/chain"
	"github.com/feyeleanor/slices"
	"io"
	"text/scanner"
)

type Symbol		string
type String		string

type Delimiter	rune

type Delimiters struct {
	Start, End	Delimiter
}

type Reader		func(*Atomiser) interface{}

type Atomiser struct {
	*scanner.Scanner
	String		Delimiters
	List		Delimiters
	Array		Delimiters
}

func NewAtomiser(f io.Reader) (a *Atomiser) {
	a = &Atomiser{
		Scanner:		new(scanner.Scanner),
		String: 		Delimiters{'"', '"'},
		List:			Delimiters{'(', ')'},
		Array:			Delimiters{'[', ']'},
	}
	a.Init(f)
	a.Mode &^= scanner.GoTokens
	return
}

func (a Atomiser) IsEOF() bool {
	return a.Peek() == scanner.EOF
}

func (a Atomiser) IsLineBreak() bool {
	c := a.Peek()
	return c == '\n' || c == '\r'
}

func (a Atomiser) IsWhitespace() bool {
	c := a.Peek()
	return c == ' ' || c == '\t' || a.IsLineBreak()
}

func (a Atomiser) SkipWhitespace() {
	for ; a.IsWhitespace(); a.Next() {}
}

func (a Atomiser) NextToken() {
	a.Next()
	a.SkipWhitespace()
}

func (a Atomiser) IsDelimiter(d Delimiter) bool {
	return d == Delimiter(a.Peek())
}

func (a Atomiser) IsStringStart() bool {
	return Delimiter(a.Peek()) == a.String.Start
}

func (a Atomiser) IsStringEnd() bool {
	return Delimiter(a.Peek()) == a.String.End
}

func (a Atomiser) IsListStart() bool {
	return Delimiter(a.Peek()) == a.List.Start
}

func (a Atomiser) IsListEnd() bool {
	return Delimiter(a.Peek()) == a.List.End
}

func (a Atomiser) IsArrayStart() bool {
	return Delimiter(a.Peek()) == a.Array.Start
}

func (a Atomiser) IsArrayEnd() bool {
	return Delimiter(a.Peek()) == a.Array.End
}

func (a Atomiser) IsValidSymbol() bool {
	return !a.IsEOF() && !a.IsWhitespace() && !a.IsListStart() && !a.IsListEnd() && !a.IsArrayStart() && !a.IsArrayEnd() && !a.IsStringStart() && !a.IsStringEnd()
}

func (a Atomiser) IsPrint() bool {
	c := a.Peek()
	return ' ' <= c && c <= '~'
}

func (a Atomiser) IsAlpha() bool {
	c := a.Peek()
	return 'A' <= c && c <= 'Z' || 'a' <= c && c <= 'z'
}

func (a Atomiser) IsRadix(r int) (ok bool) {
	c := a.Peek()
	switch {
	case r <= 10:	ok = '0' <= c && c < ('0' + r)
	case r <= 36:	if ok = '0' <= c && c <= '9'; !ok {
						r = r - 11
						ok = 'A' <= c && c <= ('A' + r) || 'a' <= c && c <= ('a' + r)
					}
	}
	return
}

func (a Atomiser) DigitValue() (r int) {
	switch c := a.Next(); {
	case '0' <= c && c <= '9':		r = c - '0'
	case 'A' <= c && c <= 'Z':		r = c - 'A' + 10
	case 'a' <= c && c <= 'z':		r = c - 'a' + 10
	default:						panic("illegal digit in character escape")
	}
	return
}

func (a Atomiser) ReadChar() (r rune) {
	if c := a.Next(); c != '\\' {
		r = c
	} else {
		switch c = a.Next(); c {
		case 'a':			r = '\a'
		case 'b':			r = '\b'
		case 'f':			r = '\f'
		case 'n':			r = '\n'
		case 'r':			r = '\r'
		case 't':			r = '\t'
		case 'v':			r = '\v'
		case '\'':			r = '\''

		case 'u':			r = (a.DigitValue() << 24) + (a.DigitValue() << 16) + (a.DigitValue() << 8) + a.DigitValue()

		case 'x':			if a.IsRadix(16) {
								if r = a.DigitValue(); a.IsRadix(16) {
									r = r * 16 + a.DigitValue()
								}
							}

		default:			switch {
							case a.IsRadix(8):		r = a.DigitValue()
													if a.IsRadix(8) {
														if r = r * 8 + a.DigitValue(); a.IsRadix(8) {
															r = r * 8 + a.DigitValue()
														}
													}
							case a.IsRadix(10):		fallthrough
							case a.IsAlpha():		panic(fmt.Sprintf("illegal character escape: \\%c", c))
							}
							r = c
		}
	}
	return
}

func (a Atomiser) ReadSymbol() (r Symbol) {
	var s	[]rune
	for ; a.IsValidSymbol(); a.Next() {
		s = append(s, a.Peek())
	}
	return Symbol(s)
}

func (a Atomiser) ReadString() (r String) {
	if !a.IsStringStart() {
		panic("Not a string")
	}
	s := []rune{}
	for a.Next(); !a.IsStringEnd(); a.Next() {
		if a.IsEOF() {
			panic("EOF in string literal")
		}
		s = append(s, a.Peek())
	}
	a.NextToken()
	return String(s)
}

func (a *Atomiser) ReadList() (c *chain.Cell) {
	var tail	*chain.Cell

	for a.NextToken(); !a.IsListEnd(); {
		if a.IsEOF() {
			panic("Unexpected EOF in list literal")
		}
		if c == nil {
			c = &chain.Cell{ Head: a.Read() }
			tail = c
		} else {
			tail.Append(a.Read())
			tail = tail.Tail
		}
	}
	a.NextToken()
	return
}

func (a *Atomiser) ReadArray() (s slices.Slice) {
	for a.NextToken(); !a.IsArrayEnd(); {
		if a.IsEOF() {
			panic("Unexpected EOF in array literal")
		}
		s = append(s, a.Read())
	}
	a.NextToken()
	return
}

func (a Atomiser) Read() (r interface{}) {
	switch a.SkipWhitespace(); {
	case a.IsListStart():		r = a.ReadList()

	case a.IsArrayStart():		r = a.ReadArray()

	case a.IsStringStart():		r = a.ReadString()

	case a.IsListEnd():			panic("Unmatched list terminator")

	case a.IsArrayEnd():		panic("Unmatched array terminator")

	default:					r = a.ReadSymbol()
	}
	return
}