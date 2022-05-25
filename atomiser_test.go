package atomiser

import (
	"testing"

	"github.com/feyeleanor/chain"
	"github.com/feyeleanor/slices"
)

func dummyReader(a *Atomiser) any {
	return nil
}

func runScanner(s string, f func(*Atomiser)) {
	scanner := NewAtomiser(s)
	for ; !scanner.IsEOF(); scanner.Next() {
		f(scanner)
	}
}

func TestIsLineBreak(t *testing.T) {
	ConfirmLineBreak := func(s string) {
		runScanner(s, func(scanner *Atomiser) {
			if !scanner.IsLineBreak() {
				t.Fatalf("IsLineBreak() at %v should be true", scanner.Pos())
			}
		})
	}
	RefuteLineBreak := func(s string) {
		scanner := NewAtomiser(s)
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
	ConfirmSkipWhitespace := func(s string, r rune) {
		scanner := NewAtomiser(s)
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
	ConfirmDelimiter := func(s string, d rune) {
		runScanner(s, func(scanner *Atomiser) {
			if !scanner.IsDelimiter(d) {
				t.Fatalf("IsDelimiter() at %v should be true", scanner.Pos())
			}
		})
	}
	RefuteDelimiter := func(s string, d rune) {
		runScanner(s, func(scanner *Atomiser) {
			if scanner.IsDelimiter(d) {
				t.Fatalf("IsDelimiter() at %v should be false", scanner.Pos())
			}
		})
	}

	ConfirmDelimiter(")", ')')
	ConfirmDelimiter("))", ')')
	RefuteDelimiter("(", ')')
}

func TestReadSymbol(t *testing.T) {
	ConfirmReadSymbol := func(s, r string) {
		if x := NewAtomiser(s).ReadSymbol(); x != Symbol(r) {
			t.Fatalf("%v.ReadSymbol() should be %v but is %v", s, r, x)
		}
	}

	ConfirmReadSymbol("", "")
	ConfirmReadSymbol("(", "")
	ConfirmReadSymbol(")", "")
	ConfirmReadSymbol("[", "")
	ConfirmReadSymbol("]", "")
	ConfirmReadSymbol("\"", "")
	ConfirmReadSymbol("\t", "")
	ConfirmReadSymbol("\r", "")
	ConfirmReadSymbol("\n", "")

	ConfirmReadSymbol("{", "{")
	ConfirmReadSymbol("}", "}")
	ConfirmReadSymbol("A", "A")
	ConfirmReadSymbol("A ", "A")
	ConfirmReadSymbol("A)", "A")
	ConfirmReadSymbol("A]", "A")
	ConfirmReadSymbol("A\"", "A")
	ConfirmReadSymbol("A\tB", "A")
	ConfirmReadSymbol("A\rB", "A")
	ConfirmReadSymbol("A\nB", "A")
	ConfirmReadSymbol("A+", "A+")
	ConfirmReadSymbol("'A", "'A")
}

func TestReadString(t *testing.T) {
	ConfirmReadString := func(s string, r string) {
		if x := NewAtomiser(s).ReadString(); x != r {
			t.Fatalf("%v.ReadString() should be %v but is %v", s, r, x)
		}
	}

	RefuteReadString := func(s string) {
		var x any
		defer func() {
			if r := recover(); r == nil {
				t.Fatalf("%v.ReadString() should fail but is %v", s, x)
			}
		}()
		x = NewAtomiser(s).ReadString()
	}

	ConfirmReadString("\"\"", "")
	ConfirmReadString("\"A\"", "A")
	ConfirmReadString("\"1\"", "1")
	RefuteReadString("")
	RefuteReadString("\"")
	RefuteReadString("\"A")
	RefuteReadString("\"1")
}

func TestReadList(t *testing.T) {
	ConsSymbols := func(values ...any) (r *chain.Cell) {
		if len(values) > 0 {
			if n, ok := values[0].(string); ok {
				r = &chain.Cell{Head: Symbol(n)}
			} else {
				r = &chain.Cell{Head: values[0]}
			}
			c := r
			for _, v := range values[1:] {
				if n, ok := v.(string); ok {
					c.Append(Symbol(n))
				} else {
					c.Append(v)
				}
				c = c.Tail
			}
		}
		return
	}

	ConfirmReadList := func(s string, r *chain.Cell) {
		if x := NewAtomiser(s).ReadList(); !r.Equal(x) {
			if x == nil {
				t.Fatalf("%v.ReadList() should be %v but is nil", s, r)
			} else {
				t.Fatalf("%v.ReadList() should be %v [%T] but is %v [%T]", s, r, r.Head, x, x.Head)
			}
		}
	}

	ConfirmReadList("()", nil)
	ConfirmReadList("()", chain.Cons())
	ConfirmReadList("()", (*chain.Cell)(nil))
	ConfirmReadList("(0)", ConsSymbols("0"))
	ConfirmReadList("((0))", chain.Cons(ConsSymbols("0")))

	ConfirmReadList("(0 1 2 3)", ConsSymbols("0", "1", "2", "3"))
	ConfirmReadList("(0 (1))", chain.Cons(Symbol("0"), ConsSymbols("1")))
	ConfirmReadList("(0 (1 (2)))", chain.Cons(Symbol("0"), chain.Cons(Symbol("1"), ConsSymbols("2"))))
	ConfirmReadList("(0 (1) (2))", chain.Cons(Symbol("0"), ConsSymbols("1"), ConsSymbols("2")))
	ConfirmReadList("((0 1 (2 3)))", chain.Cons(chain.Cons(Symbol("0"), Symbol("1"), ConsSymbols("2", "3"))))
	ConfirmReadList("(0 (1 (2 (3))))", chain.Cons(chain.Cons(Symbol("0"), Symbol("1"), chain.Cons(Symbol("2"), ConsSymbols("3")))))

	ConsStrings := func(values ...any) (r *chain.Cell) {
		if len(values) > 0 {
			if n, ok := values[0].(string); ok {
				r = &chain.Cell{Head: String(n)}
			} else {
				r = &chain.Cell{Head: values[0]}
			}
			c := r
			for _, v := range values[1:] {
				if n, ok := v.(string); ok {
					c.Append(String(n))
				} else {
					c.Append(v)
				}
				c = c.Tail
			}
		}
		return
	}

	ConfirmReadList("(\"\")", ConsStrings(""))
	ConfirmReadList("(\"A\" \"B\" \"C\" \"D\")", ConsStrings("A", "B", "C", "D"))
	ConfirmReadList("((\"A\" \"B\" (\"C\" \"D\")))", ConsStrings(ConsStrings("A", "B", ConsStrings("C", "D"))))
	ConfirmReadList("(\"A\" (\"B\" (\"C\" (\"D\"))))", ConsStrings("A", ConsStrings("B", ConsStrings("C", ConsStrings("D")))))
}

func TestReadArray(t *testing.T) {
	ConsSymbols := func(values ...any) (r slices.Slice) {
		for _, v := range values {
			if n, ok := v.(string); ok {
				r = append(r, Symbol(n))
			} else {
				r = append(r, v)
			}
		}
		return
	}

	ConfirmReadArray := func(s string, r slices.Slice) {
		if x := NewAtomiser(s).ReadArray(); !r.Equal(x) {
			t.Fatalf("%v.ReadArray() should be %v but is %v", s, r, x)
		}
	}

	ConfirmReadArray("[]", nil)
	ConfirmReadArray("[]", make(slices.Slice, 0))
	ConfirmReadArray("[0]", ConsSymbols("0"))
	ConfirmReadArray("[[0]]", ConsSymbols(ConsSymbols("0")))

	ConfirmReadArray("[0 1 2 3]", ConsSymbols("0", "1", "2", "3"))
	ConfirmReadArray("[0 [1]]", ConsSymbols("0", ConsSymbols("1")))
	ConfirmReadArray("[0 [1 [2]]]", ConsSymbols("0", ConsSymbols("1", ConsSymbols("2"))))
	ConfirmReadArray("[0 [1] [2]]", ConsSymbols("0", ConsSymbols("1"), ConsSymbols("2")))
	ConfirmReadArray("[[0 1 [2 3]]]", ConsSymbols(ConsSymbols("0", Symbol("1"), ConsSymbols("2", "3"))))
	ConfirmReadArray("[0 [1 [2 [3]]]]", ConsSymbols("0", ConsSymbols("1", ConsSymbols("2", ConsSymbols("3")))))

	ConsStrings := func(values ...any) (r slices.Slice) {
		for _, v := range values {
			if n, ok := v.(string); ok {
				r = append(r, String(n))
			} else {
				r = append(r, v)
			}
		}
		return
	}

	ConfirmReadArray("[\"\"]", ConsStrings(""))
	ConfirmReadArray("[\"A\" \"B\" \"C\" \"D\"]", ConsStrings("A", "B", "C", "D"))
	ConfirmReadArray("[[\"A\" \"B\" [\"C\" \"D\"]]]", ConsStrings(ConsStrings("A", "B", ConsStrings("C", "D"))))
	ConfirmReadArray("[\"A\" [\"B\" [\"C\" [\"D\"]]]]", ConsStrings("A", ConsStrings("B", ConsStrings("C", ConsStrings("D")))))
}
