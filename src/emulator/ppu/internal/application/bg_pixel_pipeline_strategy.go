package application

import "nes-emu/src/emulator/ppu/internal/domain"

type BGPixelPipelineStrategy struct {
}

func (s *BGPixelPipelineStrategy) FetchPixel(ppu domain.PPU) uint32 {
	return 0
}
