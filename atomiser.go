package atomiser

import(
	"fmt"
	. "github.com/feyeleanor/chain"
	"github.com/feyeleanor/slices"
	"scanner"
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