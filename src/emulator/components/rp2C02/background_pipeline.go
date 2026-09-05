package rp2C02

import "nes-emu/src/emulator/application"

type pipeline struct {
	tileIndex                  byte
	highPallette               byte
	lowPallette                byte
	lowPatternByte             byte
	highPatternByte            byte
	lowPatternShiftRegister    uint16
	highPatternShiftRegister   uint16
	lowAttributeShiftRegister  uint16
	highAttributeShiftRegister uint16
	bus                        application.MNIBus
}

func NewPipeline(bus application.MNIBus) application.PixelPipeline {
	return &pipeline{
		bus: bus,
	}
}

func (p *pipeline) StepUpPipeline(currentDot uint, vValue uint16, fineY uint16) bool {
	const total_dots_per_tile = 8
	const (
		fetch_tile_index_dot           = 2
		fetch_attr_table_dot           = 4
		fetch_lower_pattern_table_dot  = 6
		fetch_higher_pattern_table_dot = 0
	)

	switch currentDot % total_dots_per_tile {
	case fetch_tile_index_dot:
		p.tileIndex = p.getTileIndex(vValue)
	case fetch_attr_table_dot:
		p.highPallette, p.lowPallette = p.extractPaletteBits(p.readVMemory(p.getAttrTableAddress(vValue)), vValue)
	case fetch_lower_pattern_table_dot:
		p.lowPatternByte = p.getLowByteFromPatternTable(p.tileIndex, fineY)
	case fetch_higher_pattern_table_dot:
		p.highPatternByte = p.getHighByteFromPatternTable(p.tileIndex, fineY)

		p.fillPatternShiftRegister(p.highPatternByte, p.lowPatternByte)
		p.fillAttrShiftRegister(p.highPallette, p.lowPallette)

		return true
	}

	return false
}

func (p *pipeline) RenderPixel(fineX byte) uint32 {
	const highestBit = 15
	bitPosition := highestBit - fineX

	lowPatternBit := p.getBitFromBitPosition(p.lowPatternShiftRegister, bitPosition)
	highPatternBit := p.getBitFromBitPosition(p.highPatternShiftRegister, bitPosition) << 1
	lowAttrBit := p.getBitFromBitPosition(p.lowAttributeShiftRegister, bitPosition) << 2
	highAttrBit := p.getBitFromBitPosition(p.highAttributeShiftRegister, bitPosition) << 3

	p.shiftRegisters()

	pixelData := lowPatternBit | highPatternBit | lowAttrBit | highAttrBit
	return colorPallet[pixelData]
}

func (p *pipeline) GetShiftRegisters() (uint16, uint16, uint16, uint16) {
	return p.lowPatternShiftRegister, p.highPatternShiftRegister, p.lowAttributeShiftRegister, p.highAttributeShiftRegister
}

func (p *pipeline) getTileIndex(vValue uint16) uint8 {
	return p.readVMemory(base_nametable_address | (vValue & 0x0FFF))
}

func (p *pipeline) getAttrTableAddress(vValue uint16) uint16 {
	const fixed_offset = 0x03C0
	const nametable_address_mask = 0x0C00
	const coarse_y_mask = 0x03E0
	const coarse_x_mask = 0xE0

	nametable := vValue & nametable_address_mask
	offsetX := (vValue & coarse_x_mask) >> 2
	offsetY := vValue >> 2 & coarse_y_mask

	return base_nametable_address | nametable | fixed_offset | offsetY | offsetX
}

// Extrai os 2 bits da paleta apropriados a partir do byte lido da Attribute Table
func (p *pipeline) extractPaletteBits(attrByte byte, vValue uint16) (byte, byte) {
	coarseY := (vValue >> 5) & 0x1F
	coarseX := vValue & 0x1F

	// Determina o shift (0, 2, 4 ou 6) com base no quadrante do tile 2x2
	shift := ((coarseY & 2) << 1) | (coarseX & 2)

	// Retorna os 2 bits que correspondem à paleta (valores de 0 a 3)
	palette := (attrByte >> shift) & 0x03
	highPalette := palette & 0b10
	lowPalette := palette & 0b1

	if highPalette != 0 {
		highPalette = 255
	}

	if lowPalette != 0 {
		lowPalette = 255
	}

	return highPalette, lowPalette
}

func (p *pipeline) getLowByteFromPatternTable(tileIndex byte, fineY uint16) byte {
	return p.getByteFromPatternTable(tileIndex, fineY, false)
}

func (p *pipeline) getHighByteFromPatternTable(tileIndex byte, fineY uint16) byte {
	return p.getByteFromPatternTable(tileIndex, fineY, true)
}

func (p *pipeline) getByteFromPatternTable(tileIndex byte, fineY uint16, fetchHighByte bool) byte {
	const high_byte_offset = 0b1000
	var address uint16

	patternIndex := uint16(p.getPatternTableTileIndexFromControl()) << 12

	if !fetchHighByte {
		address = patternIndex | (uint16(tileIndex) << 4) | fineY
	}

	address = patternIndex | (uint16(tileIndex) << 4) | high_byte_offset | fineY
	return p.readVMemory(address)
}

func (p *pipeline) getPatternTableTileIndexFromControl() uint8 {
	const index_mask = 0b1_0000
	return (p.bus.ReadFromMemory(ppu_control) & index_mask) >> 4
}

func (p *pipeline) getBitFromBitPosition(register uint16, bitPosition byte) uint16 {
	value := register & (1 << bitPosition)
	return value >> bitPosition
}

func (p *pipeline) fillPatternShiftRegister(highByte byte, lowByte byte) {
	p.fillShiftRegister(highByte, &p.highPatternShiftRegister)
	p.fillShiftRegister(lowByte, &p.lowPatternShiftRegister)
}

func (p *pipeline) fillAttrShiftRegister(highByte byte, lowByte byte) {
	p.fillShiftRegister(highByte, &p.highAttributeShiftRegister)
	p.fillShiftRegister(lowByte, &p.lowAttributeShiftRegister)
}

func (p *pipeline) fillShiftRegister(value byte, register *uint16) {
	*register = *register | uint16(value)
}

func (p *pipeline) shiftRegisters() {
	p.lowAttributeShiftRegister = p.lowAttributeShiftRegister << 1
	p.highAttributeShiftRegister = p.highAttributeShiftRegister << 1
	p.lowPatternShiftRegister = p.lowPatternShiftRegister << 1
	p.highPatternShiftRegister = p.highPatternShiftRegister << 1
}

func (p *pipeline) readVMemory(address uint16) uint8 {
	return p.bus.ReadFromVideoMemory(address)
}
