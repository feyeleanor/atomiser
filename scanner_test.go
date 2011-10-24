package atomiser

import (
	"strings"
	"testing"
)

func TestIsLineBreak(t *testing.T) {
	ConfirmLineBreak := func(s string) {
		scanner := NewScanner(strings.NewReader(s))
		for ; !scanner.IsEOF(); scanner.Next() {
			if !scanner.IsLineBreak() {
				t.Fatalf("IsLineBreak() at %v should be true", scanner.Pos())
			}
		}
	}
	RefuteLineBreak := func(s string) {
		scanner := NewScanner(strings.NewReader(s))
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

func TestIsDelimiter(t *testing.T) {
	ConfirmDelimiter := func(s string, d Delimiter) {
		scanner := NewScanner(strings.NewReader(s))
		for ; !scanner.IsEOF(); scanner.Next() {
			if !scanner.IsDelimiter(d) {
				t.Fatalf("IsDelimiter() at %v should be true", scanner.Pos())
			}
		}
	}
	RefuteDelimiter := func(s string, d Delimiter) {
		scanner := NewScanner(strings.NewReader(s))
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