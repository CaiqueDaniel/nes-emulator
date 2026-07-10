package cpu

const signalBitMask = 0b10000000

func (c *cpu) AddWithCarry(value uint8) {
	accSign := c.acc&signalBitMask != 0
	valSign := value&signalBitMask != 0

	result := uint16(c.acc) + uint16(value)
	if c.carry {
		result++
	}

	c.acc = uint8(result)
	c.carry = result > 0xFF
	c.zero = c.acc == 0
	c.negative = c.acc&signalBitMask != 0

	resultSign := c.acc&signalBitMask != 0
	c.overflow = (accSign == valSign) && (resultSign != accSign)
}

func (c *cpu) SubtractWithCarry(value uint8) {
	c.AddWithCarry(^value)
}

func (c *cpu) CompareWithRegister(value uint8, register string) {
	var registerValue uint8

	switch register {
	case ACCUMULATOR:
		registerValue = c.acc
	case REGISTER_X:
		registerValue = c.x
	case REGISTER_Y:
		registerValue = c.y
	}

	c.zero = registerValue == value
	c.carry = registerValue >= value
	c.negative = c.isValueNegative(registerValue & value)
}

func (c *cpu) IncrementMemory(address uint16) {
	oldValue := c.memory.Read(address)
	newValue := oldValue + 1

	c.memory.Write(address, oldValue)
	c.memory.Write(address, newValue)

	c.zero = c.isValueZero(newValue)
	c.negative = c.isValueNegative(newValue)
}

func (c *cpu) IncrementRegister(register string) {
	var registerValue uint8

	switch register {
	case REGISTER_X:
		c.x++
		registerValue = c.x
	case REGISTER_Y:
		c.y++
		registerValue = c.y
	}

	c.zero = c.isValueZero(registerValue)
	c.negative = c.isValueNegative(registerValue)
}

func (c *cpu) DecrementMemory(address uint16) {
	oldValue := c.memory.Read(address)
	newValue := oldValue - 1

	c.memory.Write(address, oldValue)
	c.memory.Write(address, newValue)

	c.zero = c.isValueZero(newValue)
	c.negative = c.isValueNegative(newValue)
}

func (c *cpu) DecrementRegister(register string) {
	var registerValue uint8

	switch register {
	case REGISTER_X:
		c.x--
		registerValue = c.x
	case REGISTER_Y:
		c.y--
		registerValue = c.y
	}

	c.zero = c.isValueZero(registerValue)
	c.negative = c.isValueNegative(registerValue)
}

func (c *cpu) isValueNegative(value uint8) bool {
	return value&signalBitMask != 0
}

func (c *cpu) isValueZero(value uint8) bool {
	return value == 0
}
