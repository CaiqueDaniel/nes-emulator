package cpu

func (c *cpu) PushAccToStack() {
	c.PushValueToStack(c.acc)
}

func (c *cpu) PullAccFromStack() {
	c.acc = c.PullValueFromStack()
	c.zero = c.isValueZero(c.acc)
	c.negative = c.isValueNegative(c.acc)
}

func (c *cpu) TransferXToStack() {
	c.stackPointer = c.x
}

func (c *cpu) TransferStackToX() {
	c.x = c.stackPointer
	c.zero = c.isValueZero(c.x)
	c.negative = c.isValueNegative(c.x)
}

func (c *cpu) PushStatusIntoStack() {
	c.bFlag = true
	c.PushFlagsIntoStack()
}

func (c *cpu) PullStatusFromStack() {
	value := c.PullValueFromStack()

	c.carry = value&0b00000001 != 0
	c.zero = value&0b00000010 != 0
	c.irq = value&0b00000100 != 0
	c.decimal = value&0b00001000 != 0
	c.overflow = value&0b01000000 != 0
	c.negative = value&0b10000000 != 0
}
