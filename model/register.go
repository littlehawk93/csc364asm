package model

// Register represents a single register address in the CSC364 microprocessor architecture.
type Register string

const (
	//R0 register at address 0
	R0 Register = "r0"
	//R1 register at address 1
	R1 Register = "r1"
	//R2 register at address 2
	R2 Register = "r2"
	//R3 register at address 3
	R3 Register = "r3"
	//R4 register at address 4
	R4 Register = "r4"
	//R5 register at address 5
	R5 Register = "r5"
	//R6 register at address 6
	R6 Register = "r6"
	//R7 register at address 7
	R7 Register = "r7"
	//R8 register at address 8
	R8 Register = "r8"
	//R9 register at address 9
	R9 Register = "r9"
	//RA register at address 10
	RA Register = "ra"
	//RB register at address 11
	RB Register = "rb"
	//RC register at address 12
	RC Register = "rc"
	//RD register at address 13
	RD Register = "rd"
	//RE register at address 14
	RE Register = "re"
	//RF register at address 15
	RF Register = "rf"
	//R10 register at address 10
	R10 Register = "r10"
	//R11 register at address 11
	R11 Register = "r11"
	//R12 register at address 12
	R12 Register = "r12"
	//R13 register at address 13
	R13 Register = "r13"
	//R14 register at address 14
	R14 Register = "r14"
	//R15 register at address 15
	R15 Register = "r15"
	//RIN input register address (same as R6)
	RIN Register = "in"
	//ROUT0 first output address register (same as RD)
	ROUT0 Register = "out0"
	//ROUT1 second output address register (same as RE)
	ROUT1 Register = "out1"
	//RPC program counter register (same as RF)
	RPC Register = "pc"
)

// registerMap maps registers to their corresponding binary address value
var registerMap = map[Register]byte{
	R0:    0x00,
	R1:    0x01,
	R2:    0x02,
	R3:    0x03,
	R4:    0x04,
	R5:    0x05,
	R6:    0x06,
	R7:    0x07,
	R8:    0x08,
	R9:    0x09,
	RA:    0x0A,
	RB:    0x0B,
	RC:    0x0C,
	RD:    0x0D,
	RE:    0x0E,
	RF:    0x0F,
	RIN:   0x06,
	ROUT0: 0x0D,
	ROUT1: 0x0E,
	RPC:   0x0F,
	R10:   0x0A,
	R11:   0x0B,
	R12:   0x0C,
	R13:   0x0D,
	R14:   0x0E,
	R15:   0x0F,
}

// GetValue returns the binary value of the address of this register
// Returns the address of the register and true, or false if the register is invalid
func (me Register) GetValue() (byte, bool) {
	v, ok := registerMap[me]
	return v, ok
}
