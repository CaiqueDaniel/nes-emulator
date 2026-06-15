package internal

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
