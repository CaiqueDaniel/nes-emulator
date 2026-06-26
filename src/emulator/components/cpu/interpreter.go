package cpu

import (
	"slices"
)

type instructionSet [256]*instruction

type instruction struct {
	Method    func([]uint8)
	ArgsBytes int
}

//STA	STX	STY
//Transfer	TAX	TXA	TAY	TYA
//Arithmetic	ADC	SBC	INC	DEC	INX	DEX	INY	DEY
//Shift	ASL	LSR	ROL	ROR
//Bitwise	AND	ORA	EOR	BIT
//Compare	CMP	CPX	CPY
//Branch	BCC	BCS	BEQ	BNE	BPL	BMI	BVC	BVS
//Jump	JMP	JSR	RTS	BRK	RTI
//Stack	PHA	PLA	PHP	PLP	TXS	TSX
//Flags	CLC	SEC	CLI	SEI	CLD	SED	CLV

func (c *cpu) initInstructions() instructionSet {
	var instructionSet instructionSet

	c.appendLDAInstructions(&instructionSet)
	c.appendLDXInstructions(&instructionSet)
	c.appendLDYInstructions(&instructionSet)
	c.appendSTAInstructions(&instructionSet)

	// c.instructionSet[0xEA] = &instruction{
	// 	Method:    func(u []uint8) { c.NoOp() },
	// 	ArgsBytes: 1,
	// }
	// 0x00: &instruction{
	// 	Method:    func(u []uint8) { c.Break() },
	// 	ArgsBytes: 0,
	// },

	return instructionSet
}

func (c *cpu) appendLDAInstructions(instructionSet *instructionSet) {
	instructionSet[0xA9] = &instruction{
		Method:    func(u []uint8) { c.LoadValueIntoRegister(u[0], ACCUMULATOR) },
		ArgsBytes: 1,
	}
	instructionSet[0xA5] = &instruction{
		Method:    func(u []uint8) { c.LoadValueIntoRegister(c.GetValueByZeroPageMode(u[0]), ACCUMULATOR) },
		ArgsBytes: 1,
	}
	instructionSet[0xB5] = &instruction{
		Method:    func(u []uint8) { c.LoadValueIntoRegister(c.GetValueByZeroPageIndexedMode(u[0], c.x), ACCUMULATOR) },
		ArgsBytes: 1,
	}
	instructionSet[0xAD] = &instruction{
		Method:    func(u []uint8) { c.LoadValueIntoRegister(c.GetValueByAbsoluteMode(parseAddress(u)), ACCUMULATOR) },
		ArgsBytes: 2,
	}
	instructionSet[0xBD] = &instruction{
		Method: func(u []uint8) {
			c.LoadValueIntoRegister(c.GetValueByIndexedAbsoluteMode(parseAddress(u), c.x), ACCUMULATOR)
		},
		ArgsBytes: 2,
	}
	instructionSet[0xB9] = &instruction{
		Method: func(u []uint8) {
			c.LoadValueIntoRegister(c.GetValueByIndexedAbsoluteMode(parseAddress(u), c.y), ACCUMULATOR)
		},
		ArgsBytes: 2,
	}
	instructionSet[0xA1] = &instruction{
		Method:    func(u []uint8) { c.LoadValueIntoRegister(c.GetValueByIndexedIndirectXMode(u[0]), ACCUMULATOR) },
		ArgsBytes: 1,
	}
	instructionSet[0xB1] = &instruction{
		Method:    func(u []uint8) { c.LoadValueIntoRegister(c.GetValueByIndirectIndexedYMode(u[0]), ACCUMULATOR) },
		ArgsBytes: 1,
	}
}

func (c *cpu) appendLDXInstructions(instructionSet *instructionSet) {
	instructionSet[0xA2] = &instruction{
		Method:    func(u []uint8) { c.LoadValueIntoRegister(u[0], REGISTER_X) },
		ArgsBytes: 1,
	}
	instructionSet[0xA6] = &instruction{
		Method:    func(u []uint8) { c.LoadValueIntoRegister(c.GetValueByZeroPageMode(u[0]), REGISTER_X) },
		ArgsBytes: 1,
	}
	instructionSet[0xB6] = &instruction{
		Method:    func(u []uint8) { c.LoadValueIntoRegister(c.GetValueByZeroPageIndexedMode(u[0], c.x), REGISTER_X) },
		ArgsBytes: 1,
	}
	instructionSet[0xAE] = &instruction{
		Method:    func(u []uint8) { c.LoadValueIntoRegister(c.GetValueByAbsoluteMode(parseAddress(u)), REGISTER_X) },
		ArgsBytes: 2,
	}
	instructionSet[0xB3] = &instruction{
		Method: func(u []uint8) {
			c.LoadValueIntoRegister(c.GetValueByIndexedAbsoluteMode(parseAddress(u), c.y), REGISTER_X)
		},
		ArgsBytes: 2,
	}
}

func (c *cpu) appendLDYInstructions(instructionSet *instructionSet) {
	instructionSet[0xA0] = &instruction{
		Method:    func(u []uint8) { c.LoadValueIntoRegister(u[0], REGISTER_X) },
		ArgsBytes: 1,
	}
	instructionSet[0xA4] = &instruction{
		Method:    func(u []uint8) { c.LoadValueIntoRegister(c.GetValueByZeroPageMode(u[0]), REGISTER_X) },
		ArgsBytes: 1,
	}
	instructionSet[0xB4] = &instruction{
		Method:    func(u []uint8) { c.LoadValueIntoRegister(c.GetValueByZeroPageIndexedMode(u[0], c.x), REGISTER_X) },
		ArgsBytes: 1,
	}
	instructionSet[0xAC] = &instruction{
		Method:    func(u []uint8) { c.LoadValueIntoRegister(c.GetValueByAbsoluteMode(parseAddress(u)), REGISTER_X) },
		ArgsBytes: 2,
	}
	instructionSet[0xBC] = &instruction{
		Method: func(u []uint8) {
			c.LoadValueIntoRegister(c.GetValueByIndexedAbsoluteMode(parseAddress(u), c.y), REGISTER_X)
		},
		ArgsBytes: 2,
	}
}

func (c *cpu) appendSTAInstructions(instructionSet *instructionSet) {
	c.instructionSet[0x85] = &instruction{
		Method:    func(u []uint8) { c.StoreRegisterIntoAbsoluteMemory(uint16(u[0]), ACCUMULATOR) },
		ArgsBytes: 1,
	}

	c.instructionSet[0x95] = &instruction{
		Method: func(u []uint8) {
			c.StoreRegisterIntoAbsoluteMemory(c.GetAddressByIndexedAbsoluteMode(uint16(u[0]), c.x), ACCUMULATOR)
		},
		ArgsBytes: 1,
	}

	c.instructionSet[0x8D] = &instruction{
		Method:    func(u []uint8) { c.StoreRegisterIntoAbsoluteMemory(parseAddress(u), ACCUMULATOR) },
		ArgsBytes: 2,
	}

	c.instructionSet[0x9D] = &instruction{
		Method: func(u []uint8) {
			c.StoreRegisterIntoAbsoluteMemory(c.GetAddressByIndexedAbsoluteMode(parseAddress(u), c.x), ACCUMULATOR)
		},
		ArgsBytes: 2,
	}

	c.instructionSet[0x99] = &instruction{
		Method: func(u []uint8) {
			c.StoreRegisterIntoAbsoluteMemory(c.GetAddressByIndexedAbsoluteMode(parseAddress(u), c.y), ACCUMULATOR)
		},
		ArgsBytes: 2,
	}

	c.instructionSet[0x81] = &instruction{
		Method: func(u []uint8) {
			c.StoreRegisterIntoAbsoluteMemory(c.GetAddressByIndexedIndirectXMode(uint16(u[0])), ACCUMULATOR)
		},
		ArgsBytes: 1,
	}

	c.instructionSet[0x91] = &instruction{
		Method: func(u []uint8) {
			c.StoreRegisterIntoAbsoluteMemory(c.GetAddressByIndirectIndexedYMode(u[0]), ACCUMULATOR)
		},
		ArgsBytes: 1,
	}

}

func (c *cpu) interpretInstruction(opCode uint8) {
	instruction := c.instructionSet[opCode]

	if instruction == nil {
		return
	}

	bytes := make([]uint8, 0)

	for i := 0; i < instruction.ArgsBytes; i++ {
		bytes = append(bytes, c.memory.Read(c.programCounter))
		c.programCounter++
	}

	slices.Reverse(bytes)
	instruction.Method(bytes)
}

func parseAddress(bytes []uint8) uint16 {
	if len(bytes) != 2 {
		return 0
	}

	return (uint16(bytes[0]) << 8) + uint16(bytes[1])
}
