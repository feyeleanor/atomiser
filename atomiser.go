package atomiser

import(
//	"fmt"
	. "github.com/feyeleanor/chain"
	"github.com/feyeleanor/slices"
)

func ReadPair(s *Scanner) (r interface{}) {
	r = s.Read()
	if _, ok := r.(Delimiter); ok || r == nil {
		panic("missing item after .")
	}
	obj := s.Read()
	if _, ok := obj.(Delimiter); !ok && obj != nil {
		panic("extra item after .")
	}
	return
}

func ReadList(s *Scanner) (head *Cell) {
	if !s.IsListStart() {
		panic("Not a list")
	}
	s.NextNonWhitespace()
	if !s.IsListEnd() {
		v := s.Read()
		if s.IsEOF() {
			panic("EOF in list literal")
		}
		head = &Cell{ v, nil }
		tail := head
		for !s.IsListEnd() {
			s.NextNonWhitespace()
			switch {
			case s.IsEOF():					panic("EOF in list literal")
			case s.IsListEnd():				return
			case s.IsListStart():			tail.Append(ReadList(s))
			case s.IsDelimiter('.'):		tail.Append(ReadPair(s))
			default:						tail.Append(v)
			}
			tail = tail.Tail
		}
	}
	return
}

func ReadArray(s *Scanner) (array slices.Slice) {
	for obj := s.Read(); obj != Delimiter(']') ; obj = s.Read() {
		if obj == Delimiter('.') {
			obj = ReadPair(s)
			array = append(array, obj)
			break
		} else {
			array = append(array, obj)
		}
	}
	return
}