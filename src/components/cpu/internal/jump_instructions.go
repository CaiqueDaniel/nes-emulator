package internal

func (c *cpu) JumpProgramCounterToValue(value uint16) {
	c.programCounter = value
}

func (c *cpu) JumpProgramCounterByIndirectValue(address uint16) {
	ptrValue := c.memory.Read(address)
	lowAddress := uint8(address) + ptrValue
	highAddress := address & 0xFF00

	c.JumpProgramCounterToValue(highAddress + uint16(lowAddress))
}

func (c *cpu) JumpProgramCounterToSubRoutine(value uint16) {
	lowAddress := uint8(c.programCounter)
	highAddress := uint8(c.programCounter & 0xFF00 >> 8)

	c.PushValueToStack(highAddress)
	c.PushValueToStack(lowAddress)

	c.programCounter = value
}
