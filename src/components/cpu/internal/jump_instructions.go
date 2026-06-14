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
