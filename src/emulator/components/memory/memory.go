package memory

import (
	"nes-emu/src/emulator/application"
)

type memory struct {
	storage map[uint16]uint8
}

func NewMemory() application.Memory {
	return &memory{
		storage: make(map[uint16]uint8),
	}
}

func (mc *memory) Read(address uint16) uint8 {
	return mc.storage[address]
}

func (mc *memory) Write(address uint16, value uint8) {
	mc.storage[address] = value
}
