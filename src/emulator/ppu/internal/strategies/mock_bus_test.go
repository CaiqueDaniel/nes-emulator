package strategies

import (
	shared "nes-emu/src/emulator/shared/application"
)

type mockBus struct {
	workMemory  map[uint16]byte
	videoMemory map[uint16]byte
	nmiCalled   int
}

func newMockBus() *mockBus {
	return &mockBus{
		workMemory:  make(map[uint16]byte),
		videoMemory: make(map[uint16]byte),
	}
}

func (b *mockBus) Tick() {}

func (b *mockBus) ReadFromMemory(address uint16) uint8 {
	if v, ok := b.workMemory[address]; ok {
		return v
	}
	return 0
}

func (b *mockBus) WriteToMemory(address uint16, value uint8) {
	b.workMemory[address] = value
}

func (b *mockBus) CallNMIHandler() {
	b.nmiCalled++
}

func (b *mockBus) ReadFromVideoMemory(address uint16) uint8 {
	if v, ok := b.videoMemory[address]; ok {
		return v
	}
	return 0
}

func (b *mockBus) WriteToVideoMemory(address uint16, value uint8) {
	b.videoMemory[address] = value
}

func (b *mockBus) AtatchWorkMemory(memory shared.Memory)  {}
func (b *mockBus) AtatchVideoMemory(memory shared.Memory) {}
