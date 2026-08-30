package bus

import (
	"nes-emu/src/emulator/application"
)

type bus struct {
	workMemory  application.Memory
	videoMemory application.Memory
	ppu         application.PPU
	tickCount   uint
	nmiMethod   func()
}

func NewBus() *bus {
	return &bus{
		tickCount: 0,
		nmiMethod: func() {},
	}
}

func NewBusWithWorkMemory(memory application.Memory) *bus {
	return &bus{
		tickCount:  0,
		nmiMethod:  func() {},
		workMemory: memory,
	}
}

func (b *bus) AtatchWorkMemory(memory application.Memory) {
	b.workMemory = memory
}

func (b *bus) AtatchVideoMemory(memory application.Memory) {
	b.videoMemory = memory
}

func (b *bus) AttachPictureProcessingUnit(ppu application.PPU) {
	b.ppu = ppu
}

func (b *bus) Tick() {
	if b.ppu != nil {
		for range 3 {
			b.ppu.Render()
		}
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
	if b.workMemory == nil {
		panic("read operation on unattached work memory!")
	}

	return b.workMemory.Read(address)
}

func (b *bus) WriteToMemory(address uint16, value uint8) {
	if b.workMemory == nil {
		panic("write operation on unattached work memory!")
	}

	b.workMemory.Write(address, value)
}

func (b *bus) ReadFromVideoMemory(address uint16) uint8 {
	if b.videoMemory == nil {
		panic("read operation on unattached video memory!")
	}

	return b.videoMemory.Read(address)
}

func (b *bus) WriteToVideoMemory(address uint16, value uint8) {
	if b.videoMemory == nil {
		panic("write operation on unattached video memory!")
	}

	b.videoMemory.Write(address, value)
}

func (b *bus) GetTickCount() uint {
	return b.tickCount
}

func (b *bus) ResetTickCount() {
	b.tickCount = 0
}
