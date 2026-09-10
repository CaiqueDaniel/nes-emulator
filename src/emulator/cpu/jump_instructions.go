package cpu

func (c *cpu) JumpProgramCounterToValue(value uint16) {
	c.programCounter = value
}

func (c *cpu) JumpProgramCounterByIndirectValue(address uint16) {
	ptrValue := c.readFromMemory(address)
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

func (c *cpu) ReturnFromSubRoutine() {
	lowAddress := c.PullValueFromStack()
	highAddress := c.PullValueFromStack()

	c.programCounter = uint16(lowAddress) + uint16(highAddress)<<8
	c.programCounter++
}

func (c *cpu) ReturnFromInterrupt() {
	flags := c.PullValueFromStack()
	lowAddress := c.PullValueFromStack()
	highAddress := c.PullValueFromStack()

	c.carry = flags&0b00000001 != 0
	c.zero = flags&0b00000010 != 0
	c.irq = flags&0b00000100 != 0
	c.decimal = flags&0b00001000 != 0
	c.overflow = flags&0b01000000 != 0
	c.negative = flags&0b10000000 != 0
	c.programCounter = uint16(lowAddress) + uint16(highAddress)<<8
}

func (c *cpu) Break() {
	pcHighAddress := uint8(c.programCounter >> 8)
	pcLowAddress := uint8(c.programCounter)

	c.irq = true
	c.bFlag = true

	c.PushValueToStack(pcHighAddress)
	c.PushValueToStack(pcLowAddress)
	c.PushStatusIntoStack()

	newPcLowAddress := uint16(c.readFromMemory(0xFFFE))
	newPcHighAddress := uint16(c.readFromMemory(0xFFFF))

	c.programCounter = newPcLowAddress | (newPcHighAddress << 8)
}
