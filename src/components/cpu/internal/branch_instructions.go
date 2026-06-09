package internal

func (c *cpu) BranchIfCarryIsClear(value uint8) {
	if c.carry {
		return
	}

	c.branchByValue(value)
}

func (c *cpu) BranchIfCarryIsSet(value uint8) {
	if !c.carry {
		return
	}

	c.branchByValue(value)
}

func (c *cpu) branchByValue(value uint8) {
	if c.isValueNegative(value) {
		c.programCounter = c.programCounter - (uint16(value^0b10000000) + 1) + 2
		return
	}

	c.programCounter += uint16(value) + 2
}
