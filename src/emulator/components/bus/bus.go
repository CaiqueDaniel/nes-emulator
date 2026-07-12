package bus

import (
	"nes-emu/src/emulator/application"
)

type bus struct {
	memory    application.Memory
	tickables []application.Tickable
	tickCount uint
	nmiMethod func()
}

func NewBus(memory application.Memory) *bus {
	return &bus{
		tickables: make([]application.Tickable, 0),
		tickCount: 0,
		nmiMethod: func() {},
		memory:    memory,
	}
}

func (b *bus) AttachTickable(tickable application.Tickable) {
	b.tickables = append(b.tickables, tickable)
}

func (b *bus) Tick(cycles int) {
	for _, tickable := range b.tickables {
		tickable.Tick()
	}
	b.tickCount++
}

func (b *bus) AttachNMI(cpu application.CPU) {
	b.nmiMethod = func() { cpu.SetNMI() }
}

func (b *bus) CallNMIHandler() {
	b.nmiMethod()
}

func (b *bus) ReadFromMemory(address uint16) uint8 {
	b.Tick(1)
	return b.memory.Read(address)
}

func (b *bus) WriteToMemory(address uint16, value uint8) {
	b.Tick(1)
	b.memory.Write(address, value)
}

func (b *bus) GetTickCount() uint {
	return b.tickCount
}

func (b *bus) ResetTickCount() {
	b.tickCount = 0
}
