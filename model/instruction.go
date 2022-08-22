package model

// Instruction represents a type of instruction supported by the CSC364 assembly language
type Instruction string

// InstructionParamType defines what kind of parameter an instruction is expecting
type InstructionParamType byte

const (

	// ParamRegister parameter for a register address
	ParamRegister InstructionParamType = 1

	// ParamNumeric for a numeric constant
	ParamNumeric InstructionParamType = 2

	// MOV move value from one register to another
	MOV Instruction = "mov"
	// NOT perform binary not operation on a register
	NOT Instruction = "not"
	// AND perform binary and operation on two registers
	AND Instruction = "and"
	// OR perform binary or operation on two registers
	OR Instruction = "or"
	// ADD add two registers together
	ADD Instruction = "add"
	// SUB subtract one register from another
	SUB Instruction = "sub"
	// ADDI increment a register by a certain value
	ADDI Instruction = "addi"
	// SUBI decrement a register by a certain value
	SUBI Instruction = "subi"
	// SET set the lower 8 bits of a register
	SET Instruction = "set"
	// SETH set the higher 8 bits of a register
	SETH Instruction = "seth"
	// INCZ increment a register if another register equals zero
	INCZ Instruction = "incz"
	// DECN decrement a register if another register is negative
	DECN Instruction = "decn"
	// MOVZ move a value to another register if a register equals zero
	MOVZ Instruction = "movz"
	// MOVX move a value to another reigster if a register is not equal to zero
	MOVX Instruction = "movx"
	// MOVP move a value to another register if a register is a positive number
	MOVP Instruction = "movp"
	// MOVN move a value to another register if a register is negative
	MOVN Instruction = "movn"
	// MOVE alias for MOV
	MOVE Instruction = "move"
	// INCIZ alias for INCZ
	INCIZ Instruction = "inciz"
	// DECIN alias for DECN
	DECIN Instruction = "decin"
	// MOVEZ alias for MOVZ
	MOVEZ Instruction = "movez"
	// MOVEX alias for MOVX
	MOVEX Instruction = "movex"
	// MOVEP alias for MOVP
	MOVEP Instruction = "movep"
	// MOVEN alias for MOVN
	MOVEN Instruction = "moven"
)

// instructionMap maps an instruction to its binary OP code
var instructionMap = map[Instruction]byte{
	MOV:   0x00,
	NOT:   0x10,
	AND:   0x20,
	OR:    0x30,
	ADD:   0x40,
	SUB:   0x50,
	ADDI:  0x60,
	SUBI:  0x70,
	SET:   0x80,
	SETH:  0x90,
	INCZ:  0xA0,
	DECN:  0xB0,
	MOVZ:  0xC0,
	MOVX:  0xD0,
	MOVP:  0xE0,
	MOVN:  0xF0,
	MOVE:  0x00,
	INCIZ: 0xA0,
	DECIN: 0xB0,
	MOVEZ: 0xC0,
	MOVEX: 0xD0,
	MOVEP: 0xE0,
	MOVEN: 0xF0,
}

// instructionParamMap maps the expected parameter types for each instruction
var instructionParamMap = map[Instruction][]InstructionParamType{
	MOV:   {ParamRegister, ParamRegister},
	NOT:   {ParamRegister, ParamRegister},
	AND:   {ParamRegister, ParamRegister, ParamRegister},
	OR:    {ParamRegister, ParamRegister, ParamRegister},
	ADD:   {ParamRegister, ParamRegister, ParamRegister},
	SUB:   {ParamRegister, ParamRegister, ParamRegister},
	ADDI:  {ParamRegister, ParamRegister, ParamNumeric},
	SUBI:  {ParamRegister, ParamRegister, ParamNumeric},
	SET:   {ParamRegister, ParamNumeric, ParamNumeric},
	SETH:  {ParamRegister, ParamNumeric, ParamNumeric},
	INCZ:  {ParamRegister, ParamNumeric, ParamRegister},
	DECN:  {ParamRegister, ParamNumeric, ParamRegister},
	MOVZ:  {ParamRegister, ParamRegister, ParamRegister},
	MOVX:  {ParamRegister, ParamRegister, ParamRegister},
	MOVP:  {ParamRegister, ParamRegister, ParamRegister},
	MOVN:  {ParamRegister, ParamRegister, ParamRegister},
	MOVE:  {ParamRegister, ParamRegister},
	INCIZ: {ParamRegister, ParamNumeric, ParamRegister},
	DECIN: {ParamRegister, ParamNumeric, ParamRegister},
	MOVEZ: {ParamRegister, ParamRegister, ParamRegister},
	MOVEX: {ParamRegister, ParamRegister, ParamRegister},
	MOVEP: {ParamRegister, ParamRegister, ParamRegister},
	MOVEN: {ParamRegister, ParamRegister, ParamRegister},
}

// GetValue get the binary OP code for this instruction.
// Returns the OP code for this instruction and true, or false is the instruction isn't valid.
func (me Instruction) GetValue() (byte, bool) {
	v, ok := instructionMap[me]
	return v, ok
}

// GetParams get the set of parameter types expected for this instruction.
// Returns nil if the instruction isn't valid.
func (me Instruction) GetParams() []InstructionParamType {

	p, ok := instructionParamMap[me]

	if !ok {
		return nil
	}
	return p
}
