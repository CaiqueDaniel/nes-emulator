package internal

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

func (c *cpu) And(value uint8) {
	c.acc = c.acc & value
	c.zero = c.isValueZero(c.acc)
	c.negative = c.isValueNegative(c.acc)
}

func (c *cpu) Or(value uint8) {
	c.acc = c.acc | value
	c.zero = c.isValueZero(c.acc)
	c.negative = c.isValueNegative(c.acc)
}

func (c *cpu) Xor(value uint8) {
	c.acc = c.acc ^ value
	c.zero = c.isValueZero(c.acc)
	c.negative = c.isValueNegative(c.acc)
}

func (c *cpu) Bit(value uint8) {
	result := c.acc & value

	c.zero = c.isValueZero(result)
	c.negative = c.isValueNegative(result)
	c.overflow = result&0b01000000 != 0
}


func (c *cpu) CompareWithRegister(value uint8) {

}

func (c *cpu) isValueNegative(value uint8) bool {
	return value&signalBitMask != 0
}

func (c *cpu) isValueZero(value uint8) bool {
	return value == 0
}
