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

type ppu struct {
	scanline     uint16
	pixel        uint8
	enableRender bool
}

func NewPPU() application.PPU {
	return NewPPUWithInternal()
}

func NewPPUWithInternal() *ppu {
	return &ppu{
		enableRender: false,
	}
}

func (p *ppu) Render() {
	//render logic

	if p.pixel == max_pixel_per_scanline {
		p.scanline++
		p.scanline %= (max_frame_scanline + 1)
	}

	p.pixel++
}

func (p *ppu) GetCurrentScanline() uint16 {
	return p.scanline
}

func (p *ppu) GetCurrentScanlinePixel() uint8 {
	return p.pixel
}
