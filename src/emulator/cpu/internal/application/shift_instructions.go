package application

func (c *CPU) ArithmeticShiftLeft() {
	prevValue := c.acc
	c.acc = c.acc << 1
	c.updateFlagsOnShift(c.acc)
	c.carry = c.acc < prevValue
}

func (c *CPU) ArithmeticShiftLeftAbsolute(address uint16, xIndexedMode bool) {
	prevValue, address := c.getValueByAddress(address, xIndexedMode)
	value := prevValue << 1

	c.doubleWriteToMemory(address, prevValue, value)
	c.updateFlagsOnShift(value)

	c.carry = value < prevValue
}

func (c *CPU) ArithmeticShiftLeftZeroPage(address uint8, xIndexedMode bool) {
	c.shiftMemoryZeroPage(address, xIndexedMode,
		func(value uint8) uint8 { return value << 1 },
		func(value uint8) bool { return value&0x80 != 0 })
}

func (c *CPU) LogicalShiftRight() {
	prevValue := c.acc
	c.acc = c.acc >> 1
	c.updateFlagsOnShift(c.acc)
	c.carry = prevValue&0x1 == 1
}

func (c *CPU) LogicalShiftRightAbsolute(address uint16, xIndexedMode bool) {
	prevValue, address := c.getValueByAddress(address, xIndexedMode)
	value := prevValue >> 1

	c.doubleWriteToMemory(address, prevValue, value)
	c.updateFlagsOnShift(value)

	c.carry = prevValue&0x1 == 1
}

func (c *CPU) LogicalShiftRightZeroPage(address uint8, xIndexedMode bool) {
	c.shiftMemoryZeroPage(address, xIndexedMode,
		func(value uint8) uint8 { return value >> 1 },
		func(value uint8) bool { return value&0x01 != 0 })
}

func (c *CPU) RotateLeft() {
	prevCarry := c.carry

	c.carry = c.acc&0b10000000 != 0
	c.acc = bitShiftLeftWithCarry(c.acc, prevCarry)

	c.updateFlagsOnShift(c.acc)
}

func (c *CPU) RotateLeftAbsolute(address uint16, xIndexedMode bool) {
	prevValue, address := c.getValueByAddress(address, xIndexedMode)
	prevCarry := c.carry

	c.carry = prevValue&0b10000000 != 0

	value := bitShiftLeftWithCarry(prevValue, prevCarry)

	c.doubleWriteToMemory(address, prevValue, value)
	c.updateFlagsOnShift(value)
}

func (c *CPU) RotateLeftZeroPage(address uint8, xIndexedMode bool) {
	prevCarry := c.carry
	c.shiftMemoryZeroPage(address, xIndexedMode,
		func(value uint8) uint8 { return bitShiftLeftWithCarry(value, prevCarry) },
		func(value uint8) bool { return value&0x80 != 0 })
}

func (c *CPU) RotateRight() {
	prevCarry := c.carry

	c.carry = c.acc&0b00000001 != 0
	c.acc = bitShiftRightWithCarry(c.acc, prevCarry)

	c.updateFlagsOnShift(c.acc)
}

func (c *CPU) RotateRightAbsolute(address uint16, xIndexedMode bool) {
	prevValue, address := c.getValueByAddress(address, xIndexedMode)
	prevCarry := c.carry

	c.carry = prevValue&0b00000001 != 0

	value := bitShiftRightWithCarry(prevValue, prevCarry)

	c.doubleWriteToMemory(address, prevValue, value)
	c.updateFlagsOnShift(value)
}

func (c *CPU) RotateRightZeroPage(address uint8, xIndexedMode bool) {
	prevCarry := c.carry
	c.shiftMemoryZeroPage(address, xIndexedMode,
		func(value uint8) uint8 { return bitShiftRightWithCarry(value, prevCarry) },
		func(value uint8) bool { return value&0x01 != 0 })
}

func (c *CPU) shiftMemoryZeroPage(
	address uint8,
	xIndexedMode bool,
	transform func(uint8) uint8,
	carry func(uint8) bool,
) {
	if xIndexedMode {
		address = uint8(c.GetAddressByZeroPageIndexedModeWithDummyRead(address, c.x))
	}

	prevValue := c.GetValueByAbsoluteMode(uint16(address))
	value := transform(prevValue)
	c.doubleWriteToMemory(uint16(address), prevValue, value)
	c.updateFlagsOnShift(value)
	c.carry = carry(prevValue)
}

func (c *CPU) doubleWriteToMemory(address uint16, prevValue, currentValue uint8) {
	c.writeToMemory(address, prevValue)
	c.writeToMemory(address, currentValue)
}

func (c *CPU) updateFlagsOnShift(value uint8) {
	c.negative = c.isValueNegative(value)
	c.zero = c.isValueZero(value)
}

func (c *CPU) getValueByAddress(address uint16, xIndexedMode bool) (uint8, uint16) {
	if xIndexedMode {
		address = c.GetAddressByIndexedAbsoluteModeWithDummyRead(address, c.x)
	}

	return c.GetValueByAbsoluteMode(address), address
}

func bitShiftLeftWithCarry(value uint8, carry bool) uint8 {
	return bitShiftWithCarry(value, carry, false)
}

func bitShiftRightWithCarry(value uint8, carry bool) uint8 {
	return bitShiftWithCarry(value, carry, true)
}

func bitShiftWithCarry(value uint8, carry bool, isRight bool) uint8 {
	carryAsInt := transformFlagIntoUint8(carry)

	if isRight {
		carryAsInt <<= 7
		return value>>1 | carryAsInt
	}

	return value<<1 | carryAsInt
}
