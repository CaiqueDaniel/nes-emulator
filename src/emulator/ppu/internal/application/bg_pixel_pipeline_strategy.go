package application

import (
	"nes-emu/src/emulator/memory"
	"nes-emu/src/emulator/ppu/internal/domain"
)

const BASE_NAMETABLE_ADDRESS = 0x2000

type BGPixelPipelineStrategy struct {
	vram *memory.Memory
}

func NewBGPixelPipelineStrategy(vram *memory.Memory) *BGPixelPipelineStrategy {
	return &BGPixelPipelineStrategy{vram: vram}
}

func (s *BGPixelPipelineStrategy) FetchPixel(ppu domain.PPU) uint32 {
	return 0
}

func (s *BGPixelPipelineStrategy) GetTileIndex(ppu domain.PPU) byte {
	return s.vram.Read(ppu.GetNameTableAddress())
}

func (s *BGPixelPipelineStrategy) GetAttribute(ppu domain.PPU) byte {
	return s.vram.Read(ppu.GetAttributeTableAddress())
}
