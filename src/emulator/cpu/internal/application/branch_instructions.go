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

	previousPage := c.programCounter & 0xFF00
	c.programCounter += uint16(int16(int8(value)))

	if previousPage != c.programCounter&0xFF00 {
		c.doDummyMemoryRead(previousPage | (c.programCounter & 0x00FF))
	}
}
