package bus

import (
	"nes-emu/src/emulator/application"
)

const (
	initial_work_memory_address = 0x0
	initial_ppu_memory_address  = 0x2000
	initial_apu_memory_address  = 0x4000
	initial_io_memory_address   = 0x4018
	initial_rom_memory_address  = 0x8000
)

const (
	max_work_memory_size = 2048
	max_ppu_latches_size = 8
	max_rom_memory_size  = 16 * 1024 // 16KB
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

	return b.workMemory.Read(translateMemoryAddress(address))
}

func (b *bus) WriteToMemory(address uint16, value uint8) {
	if b.workMemory == nil {
		panic("write operation on unattached work memory!")
	}

	b.workMemory.Write(translateMemoryAddress(address), value)
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

func translateMemoryAddress(address uint16) uint16 {
	if address < initial_apu_memory_address {
		index := (address - initial_ppu_memory_address) % max_ppu_latches_size
		address = initial_ppu_memory_address + index
	}

	if address < initial_ppu_memory_address {
		index := (address - initial_work_memory_address) % max_work_memory_size
		address = initial_work_memory_address + index
	}

	return address
}
