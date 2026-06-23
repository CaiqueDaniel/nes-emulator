package bus

import (
	"nes-emu/src/emulator/application"
)

type Bus interface{}

type bus struct {
	cpu    application.CPU
	memory application.Memory
}

func NewBus(cpu application.CPU, memory application.Memory) Bus {
	return &bus{cpu, memory}
}
