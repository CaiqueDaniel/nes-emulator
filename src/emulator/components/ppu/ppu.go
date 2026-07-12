package ppu

import (
	"nes-emu/src/emulator/application"
)

const (
	ppu_control = 0x2000
	ppu_mask    = 0x2001
	ppu_status  = 0x2002
	oam_address = 0x2003
	oam_data    = 0x2004
	ppu_scroll  = 0x2005
	ppu_address = 0x2006
	ppu_data    = 0x2007
	oam_dma     = 0x4014
)

const (
	max_frame_scanline     = 261
	max_pixel_per_scanline = 255
	v_blank_scanline_start = 240
	v_blank_pixel_start    = 0
)

const frequency_in_mhz = 5.37

const (
	max_value_for_15_bits = 32767
	max_value_for_3_bits  = 7
)

type ppu struct {
	scanline     uint16
	pixel        uint8
	enableRender bool
	nmiBus       application.MNIBus
	memory       application.Memory
	v, t         uint16
	x            uint8
	w            bool
}

func NewPPU(memory application.Memory, nmiBus application.MNIBus) *ppu {
	return NewPPUWithInternal(memory, nmiBus)
}

func NewPPUWithInternal(memory application.Memory, nmiBus application.MNIBus) *ppu {
	return &ppu{
		enableRender: false,
		nmiBus:       nmiBus,
		memory:       memory,
	}
}

func (p *ppu) Render() {
	//render logic
	//have windows for fetching data
	// nametable is 1kb of memory to layout backgrond
	//pattern table is the shape of graphs for boyh sprites and background

	p.fetchDots()
	p.advanceToNextPixel()

	if p.checkIfNMIShouldBeCalled() {
		p.nmiBus.CallNMIHandler()
	}
}

func (p *ppu) GetCurrentScanline() uint16 {
	return p.scanline
}

func (p *ppu) GetCurrentScanlinePixel() uint8 {
	return p.pixel
}

func (p *ppu) fetchDots() {

}

func (p *ppu) advanceToNextPixel() {
	if p.pixel == max_pixel_per_scanline {
		p.scanline++
		p.scanline %= (max_frame_scanline + 1)
	}

	p.pixel++
}

func (p *ppu) checkIfNMIShouldBeCalled() bool {
	return p.isNMIFlagEnabled() && p.scanline == v_blank_scanline_start && p.pixel == v_blank_pixel_start
}

func (p *ppu) isNMIFlagEnabled() bool {
	const nmiFlagMask = 0b10000000
	return p.memory.Read(ppu_control)&nmiFlagMask != 0
}

func (p *ppu) setV(value uint16) {
	p.v = value % max_value_for_15_bits
}

func (p *ppu) setT(value uint16) {
	p.t = value % max_value_for_15_bits
}

func (p *ppu) setX(value uint8) {
	p.x = value % max_value_for_3_bits
}

func (p *ppu) setW(value bool) {
	p.w = value
}
