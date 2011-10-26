package atomiser

import (
	"strings"
	"testing"
)

func dummyReader(s *Scanner) interface{} {
	return nil
}

func TestIsLineBreak(t *testing.T) {
	ConfirmLineBreak := func(s string) {
		scanner := NewScanner(strings.NewReader(s), dummyReader)
		for ; !scanner.IsEOF(); scanner.Next() {
			if !scanner.IsLineBreak() {
				t.Fatalf("IsLineBreak() at %v should be true", scanner.Pos())
			}
		}
	}
	RefuteLineBreak := func(s string) {
		scanner := NewScanner(strings.NewReader(s), dummyReader)
		if scanner.IsLineBreak() {
			t.Fatalf("IsLineBreak() at %v should be false", scanner.Pos())
		}
	}

	ConfirmLineBreak("\n")
	ConfirmLineBreak("\r")
	ConfirmLineBreak("\n\r")
	ConfirmLineBreak("\r\n")
	ConfirmLineBreak("\n\r\n")
	ConfirmLineBreak("\r\n\r")

	RefuteLineBreak(" ")
}

func TestSkipWhitespace(t *testing.T) {
	ConfirmSkipWhitespace := func(s string, r int) {
		scanner := NewScanner(strings.NewReader(s), dummyReader)
		if scanner.SkipWhitespace(); scanner.Peek() != r {
			t.Fatalf("%v.SkipWhitespace() should be %v but is %v", s, r, scanner.Peek())
		}
	}
	ConfirmSkipWhitespace("", -1)
	ConfirmSkipWhitespace("A", 'A')
	ConfirmSkipWhitespace(" A", 'A')
	ConfirmSkipWhitespace("  A", 'A')
	ConfirmSkipWhitespace("   A", 'A')
	ConfirmSkipWhitespace("B   A", 'B')
	ConfirmSkipWhitespace(" \tA", 'A')
	ConfirmSkipWhitespace("  \tA", 'A')
	ConfirmSkipWhitespace("   \tA", 'A')
	ConfirmSkipWhitespace("B   \tA", 'B')
	ConfirmSkipWhitespace("\tB   A", 'B')
}

func TestIsDelimiter(t *testing.T) {
	ConfirmDelimiter := func(s string, d Delimiter) {
		scanner := NewScanner(strings.NewReader(s), dummyReader)
		for ; !scanner.IsEOF(); scanner.Next() {
			if !scanner.IsDelimiter(d) {
				t.Fatalf("IsDelimiter() at %v should be true", scanner.Pos())
			}
		}
	}
	RefuteDelimiter := func(s string, d Delimiter) {
		scanner := NewScanner(strings.NewReader(s), dummyReader)
		for ; !scanner.IsEOF(); scanner.Next() {
			if scanner.IsDelimiter(d) {
				t.Fatalf("IsDelimiter() at %v should be false", scanner.Pos())
			}
		}
	}

	ConfirmDelimiter(")", ')')
	ConfirmDelimiter("))", ')')
	RefuteDelimiter("(", ')')
}

func TestReadString(t *testing.T) {
	ConfirmReadString := func(s string, r interface{}) {
		if x := NewScanner(strings.NewReader(s), dummyReader).ReadString(); x != r {
			t.Fatalf("%v.ReadString() should be %v but is %v", s, r, x)
		}
	}

	RefuteReadString := func(s string) {
		var x interface{}
		defer func() {
			if r := recover(); r == nil {
				t.Fatalf("%v.ReadString() should fail but is %v", s, x)
			}
		}()
		x = NewScanner(strings.NewReader(s), dummyReader).ReadString()
	}

	ConfirmReadString("\"\"", "")
	ConfirmReadString("\"A\"", "A")
	ConfirmReadString("\"1\"", "1")
	RefuteReadString("")
	RefuteReadString("\"")
	RefuteReadString("\"A")
	RefuteReadString("\"1")
}

func TestReadNumber(t *testing.T) {
	ConfirmReadNumber := func(s string, r interface{}) {
		if x := NewScanner(strings.NewReader(s), dummyReader).ReadNumber(); x != r {
			t.Fatalf("%v.ReadNumber() should be %v but is %v", s, r, x)
		}
	}

	RefuteReadNumber := func(s string) {
		var x interface{}
		defer func() {
			if r := recover(); r == nil {
				t.Fatalf("%v.ReadNumber() should fail but is %v", s, x)
			}
		}()
		x = NewScanner(strings.NewReader(s), dummyReader).ReadNumber()
	}

	RefuteReadNumber("]")
	ConfirmReadNumber("0", int64(0))
	ConfirmReadNumber("1", int64(1))
	ConfirmReadNumber("2", int64(2))
	ConfirmReadNumber("3", int64(3))
	ConfirmReadNumber("4", int64(4))
	ConfirmReadNumber("5", int64(5))
	ConfirmReadNumber("6", int64(6))
	ConfirmReadNumber("7", int64(7))
	ConfirmReadNumber("8", int64(8))
	ConfirmReadNumber("9", int64(9))
	ConfirmReadNumber("10", int64(10))
	RefuteReadNumber("A")

	RefuteReadNumber("#]r0")
	ConfirmReadNumber("#2r0", int64(0))
	ConfirmReadNumber("#2r1", int64(1))
	ConfirmReadNumber("#2r10", int64(2))
	RefuteReadNumber("#2r]")
	RefuteReadNumber("#2r2")

	ConfirmReadNumber("#3r0", int64(0))
	ConfirmReadNumber("#3r1", int64(1))
	ConfirmReadNumber("#3r2", int64(2))
	ConfirmReadNumber("#3r10", int64(3))
	RefuteReadNumber("#3r3")

	ConfirmReadNumber("#03r0", int64(0))
	ConfirmReadNumber("#03r1", int64(1))
	ConfirmReadNumber("#03r2", int64(2))
	ConfirmReadNumber("#03r10", int64(3))
	RefuteReadNumber("#03r3")

	ConfirmReadNumber("#8r0", int64(0))
	ConfirmReadNumber("#8r1", int64(1))
	ConfirmReadNumber("#8r2", int64(2))
	ConfirmReadNumber("#8r3", int64(3))
	ConfirmReadNumber("#8r4", int64(4))
	ConfirmReadNumber("#8r5", int64(5))
	ConfirmReadNumber("#8r6", int64(6))
	ConfirmReadNumber("#8r7", int64(7))
	ConfirmReadNumber("#8r10", int64(8))
	RefuteReadNumber("#8r8")

	ConfirmReadNumber("#16r0", int64(0))
	ConfirmReadNumber("#16r1", int64(1))
	ConfirmReadNumber("#16r2", int64(2))
	ConfirmReadNumber("#16r3", int64(3))
	ConfirmReadNumber("#16r4", int64(4))
	ConfirmReadNumber("#16r5", int64(5))
	ConfirmReadNumber("#16r6", int64(6))
	ConfirmReadNumber("#16r7", int64(7))
	ConfirmReadNumber("#16r8", int64(8))
	ConfirmReadNumber("#16r9", int64(9))
	ConfirmReadNumber("#16rA", int64(10))
	ConfirmReadNumber("#16rB", int64(11))
	ConfirmReadNumber("#16rC", int64(12))
	ConfirmReadNumber("#16rD", int64(13))
	ConfirmReadNumber("#16rE", int64(14))
	ConfirmReadNumber("#16rF", int64(15))
	ConfirmReadNumber("#16r10", int64(16))
	RefuteReadNumber("#16rG")

	ConfirmReadNumber("08", int64(0))

	RefuteReadNumber("0b")
	ConfirmReadNumber("0b0", int64(0))
	ConfirmReadNumber("0b1", int64(1))
	ConfirmReadNumber("0b10", int64(2))
	RefuteReadNumber("0b2")

	ConfirmReadNumber("0", int64(0))
	ConfirmReadNumber("00", int64(0))
	ConfirmReadNumber("01", int64(1))
	ConfirmReadNumber("02", int64(2))
	ConfirmReadNumber("03", int64(3))
	ConfirmReadNumber("04", int64(4))
	ConfirmReadNumber("05", int64(5))
	ConfirmReadNumber("06", int64(6))
	ConfirmReadNumber("07", int64(7))
	ConfirmReadNumber("010", int64(8))

	RefuteReadNumber("0x")
	ConfirmReadNumber("0x0", int64(0))
	ConfirmReadNumber("0x1", int64(1))
	ConfirmReadNumber("0x2", int64(2))
	ConfirmReadNumber("0x3", int64(3))
	ConfirmReadNumber("0x4", int64(4))
	ConfirmReadNumber("0x5", int64(5))
	ConfirmReadNumber("0x6", int64(6))
	ConfirmReadNumber("0x7", int64(7))
	ConfirmReadNumber("0x8", int64(8))
	ConfirmReadNumber("0x9", int64(9))
	ConfirmReadNumber("0xA", int64(10))
	ConfirmReadNumber("0xB", int64(11))
	ConfirmReadNumber("0xC", int64(12))
	ConfirmReadNumber("0xD", int64(13))
	ConfirmReadNumber("0xE", int64(14))
	ConfirmReadNumber("0xF", int64(15))
	ConfirmReadNumber("0x10", int64(16))
	RefuteReadNumber("0xG")

	ConfirmReadNumber("0.", float64(0))
	ConfirmReadNumber("0.0", float64(0))
	ConfirmReadNumber("0.1", float64(0.1))
	ConfirmReadNumber("0.19", float64(0.19))
	ConfirmReadNumber("0.19", float64(0.19))

	ConfirmReadNumber("0.19+1", float64(0.19))
	ConfirmReadNumber("0.19e1", float64(0.19e+1))
	ConfirmReadNumber("0.19e+1", float64(0.19e+1))
	ConfirmReadNumber("0.19e+10", float64(0.19e+10))

	ConfirmReadNumber("0.19-1", float64(0.19))
	ConfirmReadNumber("0.19e-1", float64(0.19e-1))
	ConfirmReadNumber("0.19e-10", float64(0.19e-10))

	ConfirmReadNumber("1.", float64(1))
	ConfirmReadNumber("1.0", float64(1))
	ConfirmReadNumber("1.1", float64(1.1))
	ConfirmReadNumber("1.19", float64(1.19))

	ConfirmReadNumber("1.19+1", float64(1.19))
	ConfirmReadNumber("1.19e1", float64(1.19e+1))
	ConfirmReadNumber("1.19e+1", float64(1.19e+1))
	ConfirmReadNumber("1.19e+10", float64(1.19e+10))

	ConfirmReadNumber("1.19-1", float64(1.19))
	ConfirmReadNumber("1.19e-1", float64(1.19e-1))
	ConfirmReadNumber("1.19e-10", float64(1.19e-10))

	ConfirmReadNumber("9991.", float64(9991))
	ConfirmReadNumber("9991.0", float64(9991))
	ConfirmReadNumber("9991.1", float64(9991.1))
	ConfirmReadNumber("9991.19", float64(9991.19))

	ConfirmReadNumber("9991.19+1", float64(9991.19))
	ConfirmReadNumber("9991.19e1", float64(9991.19e+1))
	ConfirmReadNumber("9991.19e+1", float64(9991.19e+1))
	ConfirmReadNumber("9991.19e+10", float64(9991.19e+10))

	ConfirmReadNumber("9991.19-1", float64(9991.19))
	ConfirmReadNumber("9991.19e-1", float64(9991.19e-1))
	ConfirmReadNumber("9991.19e-10", float64(9991.19e-10))
}