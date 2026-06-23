package memory

import (
	"nes-emu/src/emulator/application"
)

type memory struct {
	bus     application.Bus
	storage map[uint16]uint8
}

func NewMemory(bus application.Bus) application.Memory {
	return &memory{
		bus:     bus,
		storage: make(map[uint16]uint8),
	}
}

func (mc *memory) Read(address uint16) uint8 {
	mc.bus.Tick(1)
	return mc.storage[address]
}

func (mc *memory) Write(address uint16, value uint8) {
	mc.bus.Tick(1)
	mc.storage[address] = value
}
