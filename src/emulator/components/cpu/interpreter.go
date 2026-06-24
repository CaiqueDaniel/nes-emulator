package cpu

import "slices"

type instructionSet map[uint8]*instruction

type instruction struct {
	Method    func([]uint8)
	ArgsBytes int
}

func (c *cpu) initInstructions() map[uint8]*instruction {
	return map[uint8]*instruction{
		0xA9: &instruction{
			Method:    func(u []uint8) { c.LoadValueIntoRegister(u[0], ACCUMULATOR) },
			ArgsBytes: 1,
		},
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
