package atomiser

import(
//	"fmt"
	. "github.com/feyeleanor/chain"
	"github.com/feyeleanor/slices"
)

type Delimiter	rune

type Reader		func(*Scanner) interface{}

func ReadPair(s *Scanner, read Reader) (r interface{}) {
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

func ReadList(s *Scanner, read Reader) (head *Cell) {
	if !s.IsDelimiter('(') {
		panic("Not a list")
	}
	s.Next()
	if v := read(s); v != Delimiter(')') {
		if s.IsEOF() {
			panic("EOF in list literal")
		}
		head = &Cell{ v, nil }
		tail := head
		for !s.IsDelimiter(')') {
			s.Next()
			if s.IsEOF() {
				panic("EOF in list literal")
			}

			switch v := read(s); {
			case v == Delimiter('.'):		tail.Append(ReadPair(s, read))
//			case v == Delimiter('('):		tail.Append(ReadList(s, read))
			default:						tail.Append(v)
			}
			tail = tail.Tail
		}
	}
	return
}

func ReadArray(s *Scanner, read Reader) (array slices.Slice) {
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