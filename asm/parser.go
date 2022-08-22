package asm

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/littlehawk93/csc364asm/model"
)

const (
	commentLeadingChar = '#'
)

// Parser parses data from a reader and validates it according to the ASM364 language syntax
type Parser struct {
	reader      *bufio.Scanner
	currentLine int
}

// Next retrieves the next line of assembly code from the reader and parses it.
// Returns the parsed data in binary form or any errors that occurred during reading or parsing.
// Also returns true if there is another line parsed, false otherwise.
func (me *Parser) Next() ([]byte, bool, error) {

	if ok := me.reader.Scan(); !ok {
		return nil, ok, nil
	}

	me.currentLine++
	line := strings.ToLower(strings.TrimSpace(me.reader.Text()))

	// skip empty lines or comments
	if line == "" || line[0] == commentLeadingChar {
		return me.Next()
	}

	tokens := strings.Fields(line)
	var opCode byte
	var params []model.InstructionParamType
	var ok bool

	instruction := model.Instruction(tokens[0])

	if opCode, ok = instruction.GetValue(); !ok {
		return nil, true, &SyntaxError{
			LineNumber: me.currentLine,
			Message:    fmt.Sprintf("Invalid instruction provided '%s'", strings.ToUpper(tokens[0])),
		}
	}

	if params = instruction.GetParams(); params == nil {
		return nil, true, &SyntaxError{
			LineNumber: me.currentLine,
			Message:    fmt.Sprintf("Invalid instruction provided '%s'", strings.ToUpper(tokens[0])),
		}
	}

	if len(tokens)-1 != len(params) {
		return nil, true, &SyntaxError{
			LineNumber: me.currentLine,
			Message:    fmt.Sprintf("Instruction '%s' expects %d parameters, %d provided", strings.ToUpper(tokens[0]), len(params), len(tokens)-1),
		}
	}

	data := make([]byte, 4)
	data[0] = opCode

	var tmp byte
	var err error

	for i, p := range params {
		if p == model.ParamRegister {
			if tmp, err = parseRegister(tokens[i+1]); err != nil {
				return nil, true, &SyntaxError{LineNumber: me.currentLine, Message: err.Error()}
			}
		} else {
			if tmp, err = parseNumeric(tokens[i+1]); err != nil {
				return nil, true, &SyntaxError{LineNumber: me.currentLine, Message: err.Error()}
			}
		}
		data[i+1] = tmp
	}

	return compressData(data), true, nil
}

// NewParser creates an initializes a new ASM364 parser using the provided reader.
func NewParser(r io.Reader) *Parser {
	return &Parser{
		reader:      bufio.NewScanner(r),
		currentLine: 0,
	}
}

// parseRegister parses a given string token into a register address
func parseRegister(token string) (byte, error) {

	reg := model.Register(token)

	if val, ok := reg.GetValue(); ok {
		return val, nil
	}

	return 0, fmt.Errorf("Invalid register value provided '%s'", token)
}

// parseNumeric parses a given string token into a numeric byte value
func parseNumeric(token string) (byte, error) {

	var base int

	if len(token) > 2 && token[:2] == "0x" {
		base = 16
	} else {
		base = 10
	}

	val, err := strconv.ParseInt(token[2:], base, 8)

	if err != nil {
		return 0, fmt.Errorf("Invalid numeric value provided '%s': %s", token, err.Error())
	}

	if val < 0 || val > 15 {
		return 0, fmt.Errorf("Invalid numeric value provided '%s': Value must be between 0 and 15 inclusive", token)
	}

	return byte(val), nil
}

// compressData compress the least significant 4 bits from each of the 4 bytes provided into 2 bytes
func compressData(d []byte) []byte {

	data := make([]byte, 2)
	data[0] = (d[0] & 0xF0) | (d[1] & 0x0F)
	data[1] = ((d[2] & 0x0F) << 4) | (d[3] & 0x0F)
	return data
}
