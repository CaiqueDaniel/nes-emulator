package ppu

import "nes-emu/src/emulator/application"

type pipeline struct {
	tileIndex       byte
	highPallette    byte
	lowPallette     byte
	lowPatternByte  byte
	highPatternByte byte
	vMemory         application.Memory
	bus             application.MNIBus
}

func NewPipeline(vMemory application.Memory, bus application.MNIBus) application.PixelPipeline {
	return &pipeline{
		vMemory: vMemory,
		bus:     bus,
	}
}

func (p *pipeline) StepUpPipeline(currentDot uint, vValue uint16, fineY uint16) *application.PipelineResult {
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

		return &application.PipelineResult{
			HighPalette:     p.highPallette,
			LowPalette:      p.lowPallette,
			LowPatternByte:  p.lowPatternByte,
			HighPatternByte: p.highPatternByte,
		}
	}

	return nil
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
	return p.vMemory.Read(address)
}

func (p *pipeline) getPatternTableTileIndexFromControl() uint8 {
	const index_mask = 0b1_0000
	return (p.bus.ReadFromMemory(ppu_control) & index_mask) >> 4
}

func (p *pipeline) readVMemory(address uint16) uint8 {
	return p.vMemory.Read(address)
}
