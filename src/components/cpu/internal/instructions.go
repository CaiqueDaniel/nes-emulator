package internal

const REGISTER_X = "X"
const REGISTER_Y = "Y"
const ACCUMULATOR = "ACC"

func (c *cpu) LoadValueIntoRegister(value uint8, register string) {
	switch register {
	case REGISTER_X:
		c.x = value
	case REGISTER_Y:
		c.y = value
	case ACCUMULATOR:
		c.acc = value
	}
}

func (c *cpu) LoadValueFromMemoryIntoRegister(address uint16, register string) {
	value := c.memory.Read(address)
	c.LoadValueIntoRegister(value, register)
}

func (c *cpu) LoadValueFromMemoryIntoRegisterZeroPage(address uint8, register string) {
	c.LoadValueFromMemoryIntoRegister(uint16(address), register)
}

func (c *cpu) LoadValueFromMemoryIntoAccumulatorWithIndexAbsolute(address uint16, register string) {
	var registerValue uint8

	switch register {
	case REGISTER_X:
		registerValue = c.x
	case REGISTER_Y:
		registerValue = c.y
	}

	finalAddress := address + uint16(registerValue)
	value := c.memory.Read(finalAddress)
	c.acc = value
}

func (c *cpu) LoadValueFromMemoryIntoAccumulatorWithIndexZeroPage(address uint8, register string) {
	var registerValue uint8

	switch register {
	case REGISTER_X:
		registerValue = c.x
	case REGISTER_Y:
		registerValue = c.y
	}

	finalAddress := address + registerValue
	value := c.memory.Read(uint16(finalAddress))
	c.acc = value
}
