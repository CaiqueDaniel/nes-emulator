package bus

import (
	"nes-emu/src/components/cpu"
	"nes-emu/src/components/memory"
)

type Bus struct {
	cpu    cpu.CPU
	memory memory.Memory
}

func NewBus() *Bus {
	return &Bus{}
}
