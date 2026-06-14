package internal

func (c *cpu) ArithmeticShiftLeft() {
	prevValue := c.acc
	c.acc = c.acc << 1
	c.negative = c.isValueNegative(c.acc)
	c.zero = c.isValueZero(c.acc)
	c.carry = c.acc < prevValue
}

func (c *cpu) ArithmeticShiftLeftAbsolute(address uint16, xIndexedMode bool) {
	if xIndexedMode {
		address = c.getAddressByIndexedAbsoluteMode(address, c.x)
	}

	prevValue := c.GetValueByAbsoluteMode(address)
	value := prevValue << 1

	c.memory.Write(address, prevValue)
	c.memory.Write(address, value)

	c.negative = c.isValueNegative(value)
	c.zero = c.isValueZero(value)
	c.carry = value < prevValue
}

func (c *cpu) ArithmeticShiftLeftZeroPage(address uint8, xIndexedMode bool) {
	c.ArithmeticShiftLeftAbsolute(uint16(address), xIndexedMode)
}

func (c *cpu) LogicalShiftRight() {
	prevValue := c.acc
	c.acc = c.acc >> 1
	c.negative = false
	c.zero = c.isValueZero(c.acc)
	c.carry = prevValue&0x1 == 1
}

func (c *cpu) LogicalShiftRightAbsolute(address uint16, xIndexedMode bool) {
	if xIndexedMode {
		address = c.getAddressByIndexedAbsoluteMode(address, c.x)
	}

	prevValue := c.GetValueByAbsoluteMode(address)
	value := prevValue >> 1

	c.memory.Write(address, prevValue)
	c.memory.Write(address, value)

	c.negative = false
	c.zero = c.isValueZero(value)
	c.carry = prevValue&0x1 == 1
}

func (c *cpu) LogicalShiftRightZeroPage(address uint8, xIndexedMode bool) {
	c.LogicalShiftRightAbsolute(uint16(address), xIndexedMode)
}

func (c *cpu) RotateLeft() {
	swapCarry := c.transformFlagIntoUint8(c.carry)
	c.carry = c.acc&0b10000000 != 0
	c.acc <<= 1
	c.acc |= swapCarry

	c.negative = c.isValueNegative(c.acc)
	c.zero = c.isValueZero(c.acc)
}

func (c *cpu) RotateLeftAbsolute(address uint16, xIndexedMode bool) {
	if xIndexedMode {
		address = c.getAddressByIndexedAbsoluteMode(address, c.x)
	}

	prevValue := c.GetValueByAbsoluteMode(address)
	swapCarry := c.transformFlagIntoUint8(c.carry)
	c.carry = prevValue&0b10000000 != 0
	value := prevValue << 1
	value |= swapCarry

	c.memory.Write(address, prevValue)
	c.memory.Write(address, value)

	c.negative = c.isValueNegative(value)
	c.zero = c.isValueZero(value)
}

func (c *cpu) RotateLeftZeroPage(address uint8, xIndexedMode bool) {
	c.RotateLeftAbsolute(uint16(address), xIndexedMode)
}

func (c *cpu) RotateRight() {
	swapCarry := c.transformFlagIntoUint8(c.carry) << 7
	c.carry = c.acc&0b00000001 != 0
	c.acc = c.acc >> 1
	c.acc |= swapCarry

	c.negative = c.isValueNegative(c.acc)
	c.zero = c.isValueZero(c.acc)
}

func (c *cpu) RotateRightAbsolute(address uint16, xIndexedMode bool) {
	if xIndexedMode {
		address = c.getAddressByIndexedAbsoluteMode(address, c.x)
	}

	prevValue := c.GetValueByAbsoluteMode(address)
	swapCarry := c.transformFlagIntoUint8(c.carry) << 7
	c.carry = prevValue&0b00000001 != 0
	value := prevValue >> 1
	value |= swapCarry

	c.memory.Write(address, prevValue)
	c.memory.Write(address, value)

	c.negative = c.isValueNegative(value)
	c.zero = c.isValueZero(value)
}

func (c *cpu) RotateRightZeroPage(address uint8, xIndexedMode bool) {
	c.RotateRightAbsolute(uint16(address), xIndexedMode)
}
