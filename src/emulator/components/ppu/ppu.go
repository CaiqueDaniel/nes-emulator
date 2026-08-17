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
	v_blank_scanline_end   = 260
	v_blank_pixel_start    = 0
	max_dots_per_line      = 336
)

const frequency_in_mhz = 5.37
const base_nametable_address = 0x2000

const (
	max_value_for_15_bits = 32767
	max_value_for_3_bits  = 7
)

var colorPallet = [64]uint32{
	0x7C7C7C, 0x0000FC, 0x0000BC, 0x4428BC, 0x940084, 0xA80020, 0xA81000, 0x881400,
	0x503000, 0x007800, 0x006800, 0x005800, 0x004058, 0x000000, 0x000000, 0x000000,
	0xBCBCBC, 0x0070EC, 0x3C40FC, 0x7C00FA, 0xA800B4, 0xC43800, 0xE04000, 0xE86814,
	0x9C9400, 0x44B400, 0x34C400, 0x00D840, 0x00BCBC, 0x000000, 0x000000, 0x000000,
	0x3CBCFC, 0x0078F8, 0x0058F8, 0x6844FC, 0xD800CC, 0xE40058, 0xF83800, 0xE45C10,
	0xAC7C00, 0x00B800, 0x00A800, 0x00A844, 0x008888, 0x000000, 0x000000, 0x000000,
	0xF8F8F8, 0xA4E4FC, 0xB8B8F8, 0xD8B8F8, 0xF8B8F8, 0xF8A4C0, 0xF0D0B0, 0xFCE0A8,
	0xE8D878, 0xD8F878, 0xB8F8B8, 0xB8F8D8, 0x00FCFC, 0xD8D8D8, 0x000000, 0x000000,
}

type ppu struct {
	vMemory                    application.Memory
	scanline                   uint16
	pixel                      uint8
	enableRender               bool
	bus                        application.MNIBus
	v, t                       uint16
	x                          uint8
	w                          bool
	dots                       uint16
	lowPatternShiftRegister    uint16
	highPatternShiftRegister   uint16
	lowAttributeShiftRegister  uint16
	highAttributeShiftRegister uint16
	pipeline                   application.PixelPipeline
}

func NewPPU(bus application.MNIBus, vMemory application.Memory, pipeline application.PixelPipeline) *ppu {
	p := &ppu{
		enableRender: false,
		pipeline:     pipeline,
		bus:          bus,
		vMemory:      vMemory,
	}

	return p
}

func (p *ppu) Render() {
	if p.dots <= 255 && p.scanline < v_blank_scanline_start {
		p.renderPixel()
		p.shiftRegisters()
		p.advanceToNextScanlinePixel()
	}

	p.advanceToNextScanline()
	p.fetchGraphics()
	p.updateStatusRegister()

	if p.checkIfNMIShouldBeCalled() {
		p.bus.CallNMIHandler()
	}

	p.dots = (p.dots + 1) % max_dots_per_line
}

func (p *ppu) GetCurrentScanline() uint16 {
	return p.scanline
}

func (p *ppu) GetCurrentScanlinePixel() uint8 {
	return p.pixel
}

func (p *ppu) GetShiftRegisters() (uint16, uint16, uint16, uint16) {
	return p.lowPatternShiftRegister, p.highPatternShiftRegister, p.lowAttributeShiftRegister, p.highAttributeShiftRegister
}

func (p *ppu) renderPixel() {
	//add render logic

}

func (p *ppu) shiftRegisters() {
	p.lowAttributeShiftRegister = p.lowAttributeShiftRegister << 1
	p.highAttributeShiftRegister = p.highAttributeShiftRegister << 1
	p.lowPatternShiftRegister = p.lowPatternShiftRegister << 1
	p.highPatternShiftRegister = p.highPatternShiftRegister << 1
}

func (p *ppu) fetchGraphics() {
	result := p.pipeline.StepUpPipeline(uint(p.dots), p.v, p.getFineY())

	if result == nil {
		return
	}

	p.fillPatternShiftRegister(result.HighPatternByte, result.LowPatternByte)
	p.fillAttrShiftRegister(result.HighPalette, result.LowPalette)
}

func (p *ppu) advanceToNextScanlinePixel() {
	p.pixel++
}

func (p *ppu) advanceToNextScanline() {
	if p.dots == max_dots_per_line-1 {
		p.scanline++
		p.scanline %= (max_frame_scanline + 1)
	}
}

func (p *ppu) updateStatusRegister() {
	if p.isOnPreRender() {
		value := p.bus.ReadFromMemory(ppu_status) & 0b00011111
		p.bus.WriteToMemory(ppu_status, value)
	}

	if p.isVBlankStarted() {
		value := p.bus.ReadFromMemory(ppu_status) ^ 0b10000000
		p.bus.WriteToMemory(ppu_status, value)
	}
}

func (p *ppu) checkIfNMIShouldBeCalled() bool {
	return p.isNMIFlagEnabled() && p.isVBlankStarted()
}

func (p *ppu) isNMIFlagEnabled() bool {
	const nmiFlagMask = 0b10000000
	return p.bus.ReadFromMemory(ppu_control)&nmiFlagMask != 0
}

func (p *ppu) isOnPreRender() bool {
	return p.scanline == max_frame_scanline && p.pixel == 0
}

func (p *ppu) isVBlankStarted() bool {
	return p.scanline == v_blank_scanline_start && p.pixel == v_blank_pixel_start
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

func (p *ppu) getFineY() uint16 {
	const fine_y_mask = 0x7000
	return (p.v & fine_y_mask) >> 12
}

func (p *ppu) readVMemory(address uint16) uint8 {
	return p.vMemory.Read(address)
}

func (p *ppu) fillPatternShiftRegister(highByte byte, lowByte byte) {
	p.fillShiftRegister(highByte, &p.highPatternShiftRegister)
	p.fillShiftRegister(lowByte, &p.lowPatternShiftRegister)
}

func (p *ppu) fillAttrShiftRegister(highByte byte, lowByte byte) {
	p.fillShiftRegister(highByte, &p.highAttributeShiftRegister)
	p.fillShiftRegister(lowByte, &p.lowAttributeShiftRegister)
}

func (p *ppu) fillShiftRegister(value byte, register *uint16) {
	*register = *register | uint16(value)
}
