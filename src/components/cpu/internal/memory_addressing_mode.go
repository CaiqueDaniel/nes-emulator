package internal

func (c *cpu) GetValueByAbsoluteMode(address uint16) uint8 {
	return c.memory.Read(address)
}

func (c *cpu) GetAddressByIndexedAbsoluteMode(address uint16, register string) uint16 {
	var registerValue uint8

	switch register {
	case REGISTER_X:
		registerValue = c.x
	case REGISTER_Y:
		registerValue = c.y
	}

	return address + uint16(registerValue)
}

func (c *cpu) GetAddressByIndexedZeroPageMode(address uint8, register string) uint16 {
	var registerValue uint8

	switch register {
	case REGISTER_X:
		registerValue = c.x
	case REGISTER_Y:
		registerValue = c.y
	}

	finalAddress := address + registerValue
	return uint16(finalAddress)
}

func (c *cpu) GetAddressByIndirectAbsoluteMode(initialAddress uint16) uint16 {
	lastByte := c.memory.Read(initialAddress)
	firstByte := c.memory.Read(initialAddress + 1)
	return uint16(firstByte)<<8 | uint16(lastByte)
}

func (c *cpu) GetAddressByIndexedIndirectXMode(initialAddress uint8) uint16 {
	pivotAddress := uint16(initialAddress) + uint16(c.x)
	lastByte := c.memory.Read(pivotAddress)
	firstByte := c.memory.Read(pivotAddress + 1)
	return uint16(firstByte)<<8 | uint16(lastByte)
}

func (c *cpu) GetAddressByIndirectIndexedYMode(initialAddress uint8) uint16 {
	lastByte := c.memory.Read(uint16(initialAddress))
	firstByte := c.memory.Read(uint16(initialAddress) + 1)
	return uint16(firstByte)<<8 | uint16(lastByte) + uint16(c.y)
}
