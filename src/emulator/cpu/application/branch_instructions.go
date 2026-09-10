package application

func (c *CPU) BranchIfCarryIsClear(value uint8) {
	if c.carry {
		return
	}

	c.branchByValue(value)
}

func (c *CPU) BranchIfCarryIsSet(value uint8) {
	if !c.carry {
		return
	}

	c.branchByValue(value)
}

func (c *CPU) BranchIfEqual(value uint8) {
	if !c.zero {
		return
	}

	c.branchByValue(value)
}

func (c *CPU) BranchIfNotEqual(value uint8) {
	if c.zero {
		return
	}

	c.branchByValue(value)
}

func (c *CPU) BranchIfNegative(value uint8) {
	if !c.negative {
		return
	}

	c.branchByValue(value)
}

func (c *CPU) BranchIfPositive(value uint8) {
	if c.negative {
		return
	}

	c.branchByValue(value)
}

func (c *CPU) BranchIfOverflowClear(value uint8) {
	if c.overflow {
		return
	}

	c.branchByValue(value)
}

func (c *CPU) BranchIfOverflowSet(value uint8) {
	if !c.overflow {
		return
	}

	c.branchByValue(value)
}

func (c *CPU) branchByValue(value uint8) {
	c.doDummyMemoryRead(c.programCounter)

	prevPCHighByte := c.programCounter & 0xFF00
	prevPCLowByte := c.programCounter & 0xFF

	if c.isValueNegative(value) {
		c.programCounter = c.programCounter - (uint16(value^0b10000000) + 1) + 2
	}

	if !c.isValueNegative(value) {
		c.programCounter += uint16(value) + 2
	}

	if prevPCHighByte != c.programCounter&0xFF00 {
		c.doDummyMemoryRead(prevPCHighByte | prevPCLowByte + uint16(value))
	}
}
