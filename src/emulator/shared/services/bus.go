package services

import (
	"nes-emu/src/emulator/shared/application"
	shared "nes-emu/src/emulator/shared/application"
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
	workMemory           shared.Memory
	videoMemory          shared.Memory
	ppu                  application.PPU
	tickCount            uint
	nmiMethod            func()
	lastOperationAddress uint16
	lastOperationIsWrite bool
}

func NewBus() *bus {
	return &bus{
		tickCount: 0,
		nmiMethod: func() {},
	}
}

func NewBusWithWorkMemory(memory shared.Memory) *bus {
	return &bus{
		tickCount:  0,
		nmiMethod:  func() {},
		workMemory: memory,
	}
}

func (b *bus) AtatchWorkMemory(memory shared.Memory) {
	b.workMemory = memory
}

func (b *bus) AtatchVideoMemory(memory shared.Memory) {
	b.videoMemory = memory
}

func (b *bus) AttachPictureProcessingUnit(ppu application.PPU) {
	b.ppu = ppu
}

func (b *bus) Tick() {
	if b.ppu != nil {
		for i := range 3 {
			if i == 0 {
				b.ppu.Render(&application.PPUIOEvent{
					Address: b.lastOperationAddress,
					IsWrite: b.lastOperationIsWrite,
				})
			} else {
				b.ppu.Render(nil)
			}
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

	translatedAddress := translateMemoryAddress(address)

	b.lastOperationAddress = translatedAddress
	b.lastOperationIsWrite = false

	return b.workMemory.Read(translatedAddress)
}

func (b *bus) WriteToMemory(address uint16, value uint8) {
	if b.workMemory == nil {
		panic("write operation on unattached work memory!")
	}

	translatedAddress := translateMemoryAddress(address)

	b.lastOperationAddress = translatedAddress
	b.lastOperationIsWrite = true

	b.workMemory.Write(translatedAddress, value)
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
	if address >= initial_ppu_memory_address && address < initial_apu_memory_address {
		index := (address - initial_ppu_memory_address) % max_ppu_latches_size
		address = initial_ppu_memory_address + index
	}

	if address < initial_ppu_memory_address {
		index := (address - initial_work_memory_address) % max_work_memory_size
		address = initial_work_memory_address + index
	}

	return address
}
