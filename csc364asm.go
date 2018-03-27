package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/littlehawk93/ihex"
)

// RegisterMap - Lookup map for string register names to their binary code
var RegisterMap = map[string]byte{
	"r0":   0x00,
	"r1":   0x01,
	"r2":   0x02,
	"r3":   0x03,
	"r4":   0x04,
	"r5":   0x05,
	"r6":   0x06,
	"r7":   0x07,
	"r8":   0x08,
	"r9":   0x09,
	"ra":   0x0A,
	"rb":   0x0B,
	"rc":   0x0C,
	"rd":   0x0D,
	"re":   0x0E,
	"rf":   0x0F,
	"in":   0x06,
	"out1": 0x0D,
	"out2": 0x0E,
	"pc":   0x0F,
	"r10":  0x0A,
	"r11":  0x0B,
	"r12":  0x0C,
	"r13":  0x0D,
	"r14":  0x0E,
	"r15":  0x0F,
}

// InstructionMap - Lookup map for string instruction names to their binary code
var InstructionMap = map[string]byte{
	"mov":   0x00,
	"not":   0x10,
	"and":   0x20,
	"or":    0x30,
	"add":   0x40,
	"sub":   0x50,
	"addi":  0x60,
	"subi":  0x70,
	"set":   0x80,
	"seth":  0x90,
	"incz":  0xA0,
	"decn":  0xB0,
	"movz":  0xC0,
	"movx":  0xD0,
	"movp":  0xE0,
	"movn":  0xF0,
	"move":  0x00,
	"inciz": 0xA0,
	"decin": 0xB0,
	"movez": 0xC0,
	"movex": 0xD0,
	"movep": 0xE0,
	"moven": 0xF0,
}

func main() {

	sourceFile := flag.String("f", "", "The source code file")
	outputFile := flag.String("o", "", "Optional. The output ROM file")

	flag.Parse()

	if *sourceFile == "" {
		log.Fatalln("No source file provided")
	}

	if _, err := os.Stat(*sourceFile); os.IsNotExist(err) {
		log.Fatalln("Source file does not exist")
	}

	if *outputFile == "" {
		ext := filepath.Ext(*sourceFile)

		var tmp string

		if ext == "" {
			tmp = fmt.Sprintf("%s.hex", *sourceFile)
		} else {
			tmp = strings.Replace(*sourceFile, ext, ".hex", -1)
		}

		outputFile = &tmp
	}

	file, err := os.Open(*sourceFile)

	if err != nil {
		log.Fatalf("Error opening file '%s': %s", *sourceFile, err.Error())
	}

	defer file.Close()

	var romFile ihex.HexFile
	romFile.Type = ihex.HexFileTypeI8HEX

	err = parseFile(file, &romFile)

	if err != nil {
		log.Fatal(err)
	}

	outFile, err := os.Create(*outputFile)

	if err != nil {
		log.Fatalf("Error creating file '%s': %s", *outputFile, err.Error())
	}

	_, err = romFile.WriteTo(outFile)

	if err != nil {
		log.Fatalf("Error writing file data: %s", err.Error())
	}
}

func parseFile(src *os.File, dest *ihex.HexFile) error {

	scanner := bufio.NewScanner(src)

	lineNumber := 1

	var addressOffset uint16

	addressOffset = 0

	for scanner.Scan() {

		record, err := parseLine(scanner.Text(), addressOffset)

		if err != nil {
			return fmt.Errorf("Error on line %d: %s", lineNumber, err.Error())
		}

		if record != nil {
			dest.AddRecord(record)
			addressOffset++
		}

		lineNumber++
	}

	var endFileRecord ihex.Record

	endFileRecord.AddressOffset = 0
	endFileRecord.Data = make([]byte, 0)
	endFileRecord.RecordType = ihex.RecordTypeEndOfFile

	dest.AddRecord(&endFileRecord)

	return nil
}

func parseLine(line string, recordOffset uint16) (*ihex.Record, error) {

	line = strings.TrimSpace(line)

	// Ignore empty lines and comments
	if len(line) == 0 || line[0] == '#' {
		return nil, nil
	}

	var record ihex.Record

	record.RecordType = ihex.RecordTypeData
	record.AddressOffset = recordOffset
	record.Data = make([]byte, 2)

	tokens := strings.Fields(strings.ToLower(line))

	if len(tokens) < 3 || len(tokens) > 4 {
		return nil, errors.New("Invalid number of tokens")
	}

	// First token is always the instruction
	instruction, err := parseInstruction(tokens[0])

	if err != nil {
		return nil, err
	}

	// MOVE and NOT instructions only need 2 arguments (3 tokens) all others require 3 arguments (4 tokens)
	if instruction != 0x00 && instruction != 0x10 && len(tokens) < 4 {
		return nil, errors.New("Not enough arguments provided for instruction")
	}

	// Second token is always the destination register
	regDest, err := parseRegister(tokens[1])

	if err != nil {
		return nil, err
	}

	record.Data[0] = byte(instruction | byte(regDest&0x0F))

	regA, err := parseRegA(instruction, tokens[2])

	if err != nil {
		return nil, err
	}

	var regB byte

	if len(tokens) > 3 {
		regB, err = parseRegB(instruction, tokens[3])

		if err != nil {
			return nil, err
		}
	} else {
		regB = 0
	}

	record.Data[1] = byte(byte(regA<<4) | byte(regB&0x0F))

	return &record, nil
}

func parseRegA(instruction byte, token string) (byte, error) {

	// SET, SETH, INCIZ, DECIN instructions are the only instructions that read the 2nd argument as a numeric literal
	if instruction == 0x80 || instruction == 0x90 || instruction == 0xA0 || instruction == 0xB0 {
		return parseValue(token)
	}

	// All other instructions read 2nd argument as a register address
	return parseRegister(token)
}

func parseRegB(instruction byte, token string) (byte, error) {

	// ADDI, SUBI, SET, SETH instructions are the only instructions that read the 3rd argument as a numeric literal
	if instruction == 0x60 || instruction == 0x70 || instruction == 0x80 || instruction == 0x90 {
		return parseValue(token)
	}

	// All other instructions read 3rd argument as a register address
	return parseRegister(token)
}

func parseRegister(token string) (byte, error) {

	val, ok := RegisterMap[token]

	if !ok {
		return val, fmt.Errorf("Unable to parse token: '%s' as register address", token)
	}

	return val, nil
}

func parseInstruction(token string) (byte, error) {

	val, ok := InstructionMap[token]

	if !ok {
		return val, fmt.Errorf("Unable to parse token '%s' as instruction", token)
	}

	return val, nil
}

func parseValue(token string) (byte, error) {

	var val byte

	var err error

	if len(token) > 2 && token[0:2] == "0x" {

		val, err = parseHex(token)

		if err != nil {
			return val, fmt.Errorf("Unable to parse token '%s': %s", token, err.Error())
		}
	} else {

		val, err = parseDecimal(token)

		if err != nil {
			return val, fmt.Errorf("Unable to parse token '%s': %s", token, err.Error())
		}
	}

	return val, nil
}

func parseHex(token string) (byte, error) {

	val, err := strconv.ParseInt(token, 16, 8)

	if err != nil {
		return 0, err
	}

	if val < 0 || val > 15 {
		return 0, fmt.Errorf("Invalid numeric value provided '%s'. Must be between 0 and 15 inclusive", token)
	}

	return byte(val), nil
}

func parseDecimal(token string) (byte, error) {

	val, err := strconv.ParseInt(token, 10, 8)

	if err != nil {
		return 0, err
	}

	if val < 0 || val > 15 {
		return 0, fmt.Errorf("Invalid numeric value provided '%s'. Must be between 0 and 15 inclusive", token)
	}

	return byte(val), nil
}
