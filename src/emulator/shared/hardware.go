package shared

type Memory interface {
	Write(address uint16, value byte)
	Read(address uint16) byte
}

type CPU interface {
	Run()
}

type PPU interface {
	Tick()
}

type Bus interface {
	Tick()
	AttachRAM(ram Memory)
	AttachVRAM(vram Memory)
	AttachPPU(ppu PPU)
}
