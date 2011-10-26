package atomiser

import (
	. "github.com/feyeleanor/chain"
	"strings"
	"testing"
)

func TestReadPair(t *testing.T) {}

func TestReadList(t *testing.T) {
	var read func(*Scanner) interface{}

	read = func(s *Scanner) interface{} {
		for s.SkipWhitespace(); !s.IsEOF(); s.NextNonWhitespace() {
//			switch {
//			case s.IsListStart():			return ReadList(s, read)
//			case s.IsListEnd():				return s.List.End
//			}

			switch c := s.Peek(); c {
			case '#':						fallthrough
			case '0':						fallthrough
			case '1':						fallthrough
			case '2':						fallthrough
			case '3':						fallthrough
			case '4':						fallthrough
			case '5':						fallthrough
			case '6':						fallthrough
			case '7':						fallthrough
			case '8':						fallthrough
			case '9':						return s.ReadNumber()

//			case '(':						return ReadList(s, read)
//			case ')':						return Delimiter(')')
			}
		}
		return nil
	}
	ConfirmReadList := func(s string, r *Cell) {
		if x := ReadList(NewScanner(strings.NewReader(s), read)); !r.Equal(x) {
			t.Fatalf("%v.ReadList() should be %v but is %v", s, r, x)
		}
	}

	RefuteReadList := func(s string) {
		var x interface{}
		defer func() {
			if r := recover(); r == nil {
				t.Fatalf("%v.ReadList() should fail but is %v", s, x)
			}
		}()
		x = ReadList(NewScanner(strings.NewReader(s), read))
	}

	RefuteReadList("")
	RefuteReadList("(")
	RefuteReadList(")")

	ConfirmReadList("()", Cons())
	ConfirmReadList("(0)", Cons(int64(0)))
	ConfirmReadList("(0 1 2 3)", Cons(int64(0), int64(1), int64(2), int64(3)))
	ConfirmReadList("((0 1 (2 3)))", Cons(Cons(int64(0), int64(1), Cons(int64(2), int64(3)))))
}