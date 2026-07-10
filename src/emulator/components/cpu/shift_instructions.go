package cpu

func (c *cpu) ArithmeticShiftLeft() {
	prevValue := c.acc
	c.acc = c.acc << 1
	c.updateFlagsOnShift(c.acc)
	c.carry = c.acc < prevValue
}

func (c *cpu) ArithmeticShiftLeftAbsolute(address uint16, xIndexedMode bool) {
	prevValue, address := c.getValueByAddress(address, xIndexedMode)
	value := prevValue << 1

	c.doubleWriteToMemory(address, prevValue, value)
	c.updateFlagsOnShift(value)

	c.carry = value < prevValue
}

func (c *cpu) ArithmeticShiftLeftZeroPage(address uint8, xIndexedMode bool) {
	c.ArithmeticShiftLeftAbsolute(uint16(address), xIndexedMode)
}

func (c *cpu) LogicalShiftRight() {
	prevValue := c.acc
	c.acc = c.acc >> 1
	c.updateFlagsOnShift(c.acc)
	c.carry = prevValue&0x1 == 1
}

func (c *cpu) LogicalShiftRightAbsolute(address uint16, xIndexedMode bool) {
	prevValue, address := c.getValueByAddress(address, xIndexedMode)
	value := prevValue >> 1

	c.doubleWriteToMemory(address, prevValue, value)
	c.updateFlagsOnShift(value)

	c.carry = prevValue&0x1 == 1
}

func (c *cpu) LogicalShiftRightZeroPage(address uint8, xIndexedMode bool) {
	c.LogicalShiftRightAbsolute(uint16(address), xIndexedMode)
}

func (c *cpu) RotateLeft() {
	prevCarry := c.carry

	c.carry = c.acc&0b10000000 != 0
	c.acc = bitShiftLeftWithCarry(c.acc, prevCarry)

	c.updateFlagsOnShift(c.acc)
}

func (c *cpu) RotateLeftAbsolute(address uint16, xIndexedMode bool) {
	prevValue, address := c.getValueByAddress(address, xIndexedMode)
	prevCarry := c.carry

	c.carry = prevValue&0b10000000 != 0

	value := bitShiftLeftWithCarry(prevValue, prevCarry)

	c.doubleWriteToMemory(address, prevValue, value)
	c.updateFlagsOnShift(value)
}

func (c *cpu) RotateLeftZeroPage(address uint8, xIndexedMode bool) {
	c.RotateLeftAbsolute(uint16(address), xIndexedMode)
}

func (c *cpu) RotateRight() {
	prevCarry := c.carry

	c.carry = c.acc&0b00000001 != 0
	c.acc = bitShiftRightWithCarry(c.acc, prevCarry)

	c.updateFlagsOnShift(c.acc)
}

func (c *cpu) RotateRightAbsolute(address uint16, xIndexedMode bool) {
	prevValue, address := c.getValueByAddress(address, xIndexedMode)
	prevCarry := c.carry

	c.carry = prevValue&0b00000001 != 0

	value := bitShiftRightWithCarry(prevValue, prevCarry)

	c.doubleWriteToMemory(address, prevValue, value)
	c.updateFlagsOnShift(value)
}

func (c *cpu) RotateRightZeroPage(address uint8, xIndexedMode bool) {
	c.RotateRightAbsolute(uint16(address), xIndexedMode)
}

func (c *cpu) doubleWriteToMemory(address uint16, prevValue, currentValue uint8) {
	c.memory.Write(address, prevValue)
	c.memory.Write(address, currentValue)
}

func (c *cpu) updateFlagsOnShift(value uint8) {
	c.negative = c.isValueNegative(value)
	c.zero = c.isValueZero(value)
}

func (c *cpu) getValueByAddress(address uint16, xIndexedMode bool) (uint8, uint16) {
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
