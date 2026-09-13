package application

func (c *CPU) GetValueByAbsoluteMode(address uint16) uint8 {
	return c.readFromMemory(address)
}

func (c *CPU) GetValueByZeroPageMode(address uint8) uint8 {
	return c.readFromMemory(uint16(address))
}

func (c *CPU) GetValueByIndexedAbsoluteMode(address uint16, index uint8) uint8 {
	return c.readFromMemory(c.GetAddressByIndexedAbsoluteMode(address, index))
}

func (c *CPU) GetValueByZeroPageIndexedMode(address uint8, index uint8) uint8 {
	return c.readFromMemory(uint16(address) + uint16(index))
}

func (c *CPU) GetValueByZeroPageIndexedModeWithDummyRead(address uint8, index uint8) uint8 {
	c.doDummyMemoryRead(uint16(address))
	return c.readFromMemory(uint16(address) + uint16(index))
}

func (c *CPU) GetValueByIndirectAbsoluteMode(initialAddress uint16) uint8 {
	lastByte := c.readFromMemory(initialAddress)
	firstByte := c.readFromMemory(initialAddress + 1)

	return c.readFromMemory(uint16(firstByte)<<8 | uint16(lastByte))
}

func (c *CPU) GetValueByIndexedIndirectXMode(initialAddress uint8) uint8 {
	c.readFromMemory(uint16(initialAddress))
	return c.readFromMemory(c.GetAddressByIndexedIndirectXMode(uint16(initialAddress)))
}

func (c *CPU) GetValueByIndirectIndexedYMode(initialAddress uint8) uint8 {
	return c.readFromMemory(c.GetAddressByIndirectIndexedYMode(initialAddress))
}

func (c *CPU) GetAddressByIndexedAbsoluteMode(address uint16, index uint8) uint16 {
	prevHighAddress := uint8(address >> 8)
	newAddress := address + uint16(index)
	currentHighAddress := uint8(newAddress >> 8)

	if prevHighAddress != currentHighAddress {
		c.doDummyMemoryRead(address)
	}

	return newAddress
}

func (c *CPU) GetAddressByIndexedAbsoluteModeWithDummyRead(address uint16, index uint8) uint16 {
	c.doDummyMemoryRead(address)
	return c.GetAddressByIndexedAbsoluteMode(address, index)
}

func (c *CPU) GetAddressByIndexedIndirectXModeWithDummyRead(initialAddress uint16) uint16 {
	c.doDummyMemoryRead(initialAddress)
	return c.GetAddressByIndexedIndirectXMode(initialAddress)
}

func (c *CPU) GetAddressByIndexedIndirectXMode(initialAddress uint16) uint16 {
	pivotAddress := uint16(initialAddress) + uint16(c.x)
	lastByte := c.readFromMemory(pivotAddress)
	firstByte := c.readFromMemory(pivotAddress + 1)

	return uint16(firstByte)<<8 | uint16(lastByte)
}

func (c *CPU) GetAddressByIndirectIndexedYMode(initialAddress uint8) uint16 {
	lowByte := c.readFromMemory(uint16(initialAddress))
	highByte := c.readFromMemory(uint16(initialAddress) + 1)
	baseAddress := uint16(highByte)<<8 | uint16(lowByte)
	newAddresss := baseAddress + uint16(c.y)

	if uint8(newAddresss>>8) != highByte {
		c.doDummyMemoryRead(baseAddress)
	}

	return newAddresss
}

func (c *CPU) GetAddressByIndirectIndexedYModeWithDummyRead(initialAddress uint8) uint16 {
	lastByte := c.readFromMemory(uint16(initialAddress))
	firstByte := c.readFromMemory(uint16(initialAddress) + 1)
	incompleteAddress := uint16(firstByte)<<8 | uint16(lastByte)

	c.doDummyMemoryRead(incompleteAddress)

	return incompleteAddress + uint16(c.y)
}
