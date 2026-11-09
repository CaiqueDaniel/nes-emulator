package application

import (
	"nes-emu/src/emulator/memory"
	"nes-emu/src/emulator/ppu/internal/domain"
)

const (
	FETCH_TILE_INDEX_CYCLE  = 2
	FETCH_ATTRIBUTE_CYCLE   = 4
	FETCH_LOW_PATTERN_BYTE  = 6
	FETCH_HIGH_PATTERN_BYTE = 0
	TOTAL_PIPELINE_CYCLES   = 8
)

const BASE_NAMETABLE_ADDRESS = 0x2000

type BGPixelPipelineStrategy struct {
	vram *memory.Memory
	ram  *memory.Memory
}

func NewBGPixelPipelineStrategy(vram *memory.Memory, ram *memory.Memory) *BGPixelPipelineStrategy {
	return &BGPixelPipelineStrategy{vram: vram, ram: ram}
}

func (s *BGPixelPipelineStrategy) FetchPixel(ppu *domain.PPU) uint32 {
	switch {
	case ppu.GetDot()%TOTAL_PIPELINE_CYCLES == FETCH_TILE_INDEX_CYCLE:
		ppu.SetCurrentBgTileIndex(s.GetTileIndex(ppu))
	case ppu.GetDot()%TOTAL_PIPELINE_CYCLES == FETCH_ATTRIBUTE_CYCLE:
		ppu.SetCurrentBgAttribute(s.GetAttribute(ppu))
	case ppu.GetDot()%TOTAL_PIPELINE_CYCLES == FETCH_LOW_PATTERN_BYTE:
		ppu.SetLowBgPatternLatch(s.GetLowPatternByte(ppu))
	case ppu.GetDot()%TOTAL_PIPELINE_CYCLES == FETCH_HIGH_PATTERN_BYTE:
		ppu.SetHighBgPatternLatch(s.GetHighPatternByte(ppu))
		ppu.LoadBgShiftRegistersFromLatches()
		return s.FetchColor(ppu)
	}

	return 0x0
}

func (s *BGPixelPipelineStrategy) GetTileIndex(ppu *domain.PPU) byte {
	return s.vram.Read(ppu.GetNameTableAddress())
}

func (s *BGPixelPipelineStrategy) GetAttribute(ppu *domain.PPU) byte {
	return s.vram.Read(ppu.GetAttributeTableAddress())
}

func (s *BGPixelPipelineStrategy) GetLowPatternAddress(ppu *domain.PPU) uint16 {
	const ppu_ctrl_bg_pattern_table_bit = 0b1_0000

	bgPatternTableBit := (s.ram.Read(domain.PPU_CONTROL) & ppu_ctrl_bg_pattern_table_bit) >> 4

	return ((uint16(bgPatternTableBit)) << 12) | (uint16(ppu.GetCurrentBgTileIndex()) << 4) | uint16(ppu.GetFineY())
}

func (s *BGPixelPipelineStrategy) GetHighPatternAddress(ppu *domain.PPU) uint16 {
	return s.GetLowPatternAddress(ppu) + 8
}

func (s *BGPixelPipelineStrategy) GetLowPatternByte(ppu *domain.PPU) byte {
	return s.vram.Read(s.GetLowPatternAddress(ppu))
}

func (s *BGPixelPipelineStrategy) GetHighPatternByte(ppu *domain.PPU) byte {
	return s.vram.Read(s.GetHighPatternAddress(ppu))
}

func (s *BGPixelPipelineStrategy) FetchColor(ppu *domain.PPU) uint32 {
	const palette_offset uint16 = 0x3F00
	var colorIndex byte = s.vram.Read(uint16(palette_offset))

	index := ppu.GetPaletteIndex()

	if index&0b11 != 0 {
		colorIndex = s.vram.Read(palette_offset | uint16(colorIndex))
	}

	return domain.ColorPallet[colorIndex]
}
