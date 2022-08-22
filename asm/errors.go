package asm

import "fmt"

// SyntaxError provides an error message for any illegally formed syntax during parsing
type SyntaxError struct {
	Message    string
	LineNumber int
}

func (me *SyntaxError) Error() string {
	return fmt.Sprintf("Syntax error on line %d: %s", me.LineNumber, me.Message)
}
