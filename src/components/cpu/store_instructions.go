package cpu

func (c *cpu) StoreRegisterIntoAbsoluteMemory(value uint16, register string) {
	switch register {
	case REGISTER_X:
		c.memory.Write(value, c.x)
	case REGISTER_Y:
		c.memory.Write(value, c.y)
	case ACCUMULATOR:
		c.memory.Write(value, c.acc)
	}
}
