package sim6502

import "fmt"

type instructionImpl interface {
	Mnemonic() string
	// Exec executes the specified instruction
	// The uint8 will contain, as appropriate, the byte being operated
	// the uint16 will contain the memory address (where memory is accessed)
	// or jump address depending on addressing mode
	Exec(*Processor, AddressingMode, uint8, uint16) error
}

// instruction consists of it's mnemonic, addressing mode and opcode
type instruction struct {
	Impl           instructionImpl
	AddressingMode AddressingMode
	OpCode         uint8
	BaseCycles     uint8
}

// variableCycleInstruction is implemented by instructions who may have a variable
// cycle count depending on circumstances
// It returns the number of additional cycles, over and above the default for the instruction/mode
// that should be used
// instructionPC is the PC where the instruction was located
// address is the target address
// xPageBoundary is true if an ABS_X, ABS_Y or IND_Y increment to absolute address crossed a page boundary
type variableCycleInstruction interface {
	TweakCycles(m AddressingMode, instructionPC uint16, address uint16, xPageBoundary bool) uint8
}

// instructions is a complete instruction set indexed by opcode
var instructions = []*instruction{
	{&brk{}, IMMED, 0x00, 7},
	{&ora{}, X_IND, 0x01, 6},
	nil,
	nil,
	nil,
	{&ora{}, ZPG, 0x05, 3},
	{&asl{}, ZPG, 0x06, 5},
	nil,
	{&php{}, IMPL, 0x08, 3},
	{&ora{}, IMMED, 0x09, 2},
	{&asl{}, A, 0x0A, 2},
	nil,
	nil,
	{&ora{}, ABS, 0x0D, 4},
	{&asl{}, ABS, 0x0E, 4},
	nil,

	{&bpl{}, REL, 0x10, 2},
	{&ora{}, IND_Y, 0x11, 6},
	nil,
	nil,
	nil,
	{&ora{}, ZPG_X, 0x15, 4},
	{&asl{}, ZPG_X, 0x16, 6},
	nil,
	{&clc{}, IMPL, 0x18, 2},
	{&ora{}, ABS_Y, 0x19, 4},
	nil,
	nil,
	nil,
	{&ora{}, ABS_X, 0x1D, 4},
	{&asl{}, ABS_X, 0x1E, 7},
	nil,

	{&jsr{}, ABS, 0x20, 6},
	{&and{}, X_IND, 0x21, 6},
	nil,
	nil,
	{&bit{}, ZPG, 0x24, 3},
	{&and{}, ZPG, 0x25, 3},
	{&rol{}, ZPG, 0x26, 5},
	nil,
	{&plp{}, IMPL, 0x28, 4},
	{&and{}, IMMED, 0x29, 2},
	{&rol{}, A, 0x2A, 2},
	nil,
	{&bit{}, ABS, 0x2C, 4},
	{&and{}, ABS, 0x2D, 4},
	{&rol{}, ABS, 0x2E, 6},
	nil,

	{&bmi{}, REL, 0x30, 2},
	{&and{}, IND_Y, 0x31, 5},
	nil,
	nil,
	nil,
	{&and{}, ZPG_X, 0x35, 4},
	{&rol{}, ZPG_X, 0x36, 6},
	nil,
	{&SEC{}, IMPL, 0x38, 2},
	{&and{}, ABS_Y, 0x39, 4},
	nil,
	nil,
	nil,
	{&and{}, ABS_X, 0x3D, 4},
	{&rol{}, ABS_X, 0x3E, 7},
	nil,

	{&rti{}, IMPL, 0x40, 6},
	{&eor{}, X_IND, 0x41, 6},
	nil,
	nil,
	nil,
	{&eor{}, ZPG, 0x45, 3},
	{&lsr{}, ZPG, 0x46, 5},
	nil,
	{&pha{}, IMPL, 0x48, 3},
	{&eor{}, IMMED, 0x49, 2},
	{&lsr{}, A, 0x4A, 2},
	nil,
	{&jmp{}, ABS, 0x4C, 3},
	{&eor{}, ABS, 0x4D, 4},
	{&lsr{}, ABS, 0x4E, 6},
	nil,
	{&bvc{}, REL, 0x50, 2},
	{&eor{}, IND_Y, 0x51, 5},
	nil,
	nil,
	nil,
	{&eor{}, ZPG_X, 0x55, 4},
	{&lsr{}, ZPG_X, 0x56, 6},
	nil,
	{&cli{}, IMPL, 0x58, 2},
	{&eor{}, ABS_Y, 0x59, 4},
	nil,
	nil,
	nil,
	{&eor{}, ABS_X, 0x5D, 4},
	{&lsr{}, ABS_X, 0x5E, 7},
	nil,

	{&rts{}, IMPL, 0x60, 6},
	{&adc{}, X_IND, 0x61, 6},
	nil,
	nil,
	nil,
	{&adc{}, ZPG, 0x65, 3},
	{&ror{}, ZPG, 0x66, 5},
	nil,
	{&pla{}, IMPL, 0x68, 4},
	{&adc{}, IMMED, 0x69, 2},
	{&ror{}, A, 0x6A, 2},
	nil,
	{&jmp{}, IND, 0x6C, 5},
	{&adc{}, ABS, 0x6D, 4},
	{&ror{}, ABS, 0x6E, 6},
	nil,

	{&bvs{}, REL, 0x70, 2},
	{&adc{}, IND_Y, 0x71, 5},
	nil,
	nil,
	nil,
	{&adc{}, ZPG_X, 0x75, 4},
	{&ror{}, ZPG_X, 0x76, 6},
	nil,
	{&SEI{}, IMPL, 0x78, 2},
	{&adc{}, ABS_Y, 0x79, 4},
	nil,
	nil,
	nil,
	{&adc{}, ABS_X, 0x7d, 4},
	{&ror{}, ABS_X, 0x7e, 7},
	nil,

	nil,
	{&sta{}, X_IND, 0x81, 6},
	nil,
	nil,
	{&STY{}, ZPG, 0x84, 3},
	{&sta{}, ZPG, 0x85, 3},
	{&STX{}, ZPG, 0x86, 3},
	nil,
	{&dey{}, IMPL, 0x88, 2},
	nil,
	{&TXA{}, IMPL, 0x8A, 2},
	nil,
	{&STY{}, ABS, 0x8C, 4},
	{&sta{}, ABS, 0x8D, 4},
	{&STX{}, ABS, 0x8E, 4},
	nil,

	{&bcc{}, REL, 0x90, 2},
	{&sta{}, IND_Y, 0x91, 6},
	nil,
	nil,
	{&STY{}, ZPG_X, 0x94, 4},
	{&sta{}, ZPG_X, 0x95, 4},
	{&STX{}, ZPG_Y, 0x96, 4},
	nil,
	{&TYA{}, IMPL, 0x98, 2},
	{&sta{}, ABS_Y, 0x99, 5},
	{&TXS{}, IMPL, 0x9A, 2},
	nil,
	nil,
	{&sta{}, ABS_X, 0x9D, 5},
	nil,
	nil,

	{&ldy{}, IMMED, 0xA0, 2},
	{&lda{}, X_IND, 0xA1, 6},
	{&ldx{}, IMMED, 0xA2, 2},
	nil,
	{&ldy{}, ZPG, 0xA4, 3},
	{&lda{}, ZPG, 0xA5, 3},
	{&ldx{}, ZPG, 0xA6, 3},
	nil,
	{&TAY{}, IMPL, 0xA8, 2},
	{&lda{}, IMMED, 0xA9, 2},
	{&TAX{}, IMPL, 0xAA, 2},
	nil,
	{&ldy{}, ABS, 0xAC, 4},
	{&lda{}, ABS, 0xAD, 4},
	{&ldx{}, ABS, 0xAE, 4},
	nil,

	{&bcs{}, REL, 0xB0, 2},
	{&lda{}, IND_Y, 0xB1, 5},
	nil,
	nil,
	{&ldy{}, ZPG_X, 0xB4, 4},
	{&lda{}, ZPG_X, 0xB5, 4},
	{&ldx{}, ZPG_Y, 0xB6, 4},
	nil,
	{&clv{}, IMPL, 0xB8, 2},
	{&lda{}, ABS_Y, 0xB9, 4},
	{&TSX{}, IMPL, 0xBA, 2},
	nil,
	{&ldy{}, ABS_X, 0xBC, 4},
	{&lda{}, ABS_X, 0xBD, 4},
	{&ldx{}, ABS_Y, 0xBE, 4},
	nil,

	{&cpy{}, IMMED, 0xC0, 2},
	{&cmp{}, X_IND, 0xC1, 6},
	nil,
	nil,
	{&cpy{}, ZPG, 0xC4, 3},
	{&cmp{}, ZPG, 0xC5, 3},
	{&dec{}, ZPG, 0xC6, 5},
	nil,
	{&iny{}, IMPL, 0xC8, 2},
	{&cmp{}, IMMED, 0xC9, 2},
	{&dex{}, IMPL, 0xCA, 2},
	nil,
	{&cpy{}, ABS, 0xCC, 4},
	{&cmp{}, ABS, 0xCD, 4},
	{&dec{}, ABS, 0xCE, 6},
	nil,

	{&bne{}, REL, 0xD0, 2},
	{&cmp{}, IND_Y, 0xD1, 5},
	nil,
	nil,
	nil,
	{&cmp{}, ZPG_X, 0xD5, 4},
	{&dec{}, ZPG_X, 0xD6, 6},
	nil,
	{&cld{}, IMPL, 0xD8, 2},
	{&cmp{}, ABS_Y, 0xD9, 4},
	nil,
	nil,
	nil,
	{&cmp{}, ABS_X, 0xDD, 4},
	{&dec{}, ABS_X, 0xDE, 7},
	nil,

	{&cpx{}, IMMED, 0xE0, 2},
	{&sbc{}, X_IND, 0xE1, 6},
	nil,
	nil,
	{&cpx{}, ZPG, 0xE4, 3},
	{&sbc{}, ZPG, 0xE5, 3},
	{&inc{}, ZPG, 0xE6, 5},
	nil,
	{&inx{}, IMPL, 0xE8, 2},
	{&sbc{}, IMMED, 0xE9, 2},
	{&nop{}, IMPL, 0xEA, 2},
	nil,
	{&cpx{}, ABS, 0xEC, 4},
	{&sbc{}, ABS, 0xED, 4},
	{&inc{}, ABS, 0xEE, 6},
	nil,

	{&beq{}, REL, 0xF0, 2},
	{&sbc{}, IND_Y, 0xF1, 5},
	nil,
	nil,
	nil,
	{&sbc{}, ZPG_X, 0xF5, 4},
	{&inc{}, ZPG_X, 0xF6, 6},
	nil,
	{&SED{}, IMPL, 0xF8, 2},
	{&sbc{}, ABS_Y, 0xF9, 4},
	nil,
	nil,
	nil,
	{&sbc{}, ABS_X, 0xFD, 4},
	{&inc{}, ABS_X, 0xFE, 7},
	nil,
}

func loadAdditionalInstruction(instructions []*instruction, opcode uint8, impl instructionImpl, am AddressingMode, cycles uint8) {
	if instructions[opcode] != nil {
		panic(fmt.Sprintf("Cannot load new instruction over existing instruction at opcode 0x%02x", opcode))
	}
	instructions[opcode] = &instruction{Impl: impl, AddressingMode: am, OpCode: opcode, BaseCycles: cycles}
}
