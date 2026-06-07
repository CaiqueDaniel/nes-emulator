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

func And() {}

func Or() {}

func Xor() {}

func Bit() {}

func CompareWithRegister() {}
