package application

import "nes-emu/src/emulator/ppu/internal/domain"

const (
	BG_PIPELINE  = 0
	OAM_PIPELINE = 1
)

type PixelPipelineContext interface {
	FullfillFrameBuffer(ppu *domain.PPU, pipeline uint)
}

type PixelPipelineStrategy interface {
	FetchPixel(ppu *domain.PPU) uint32
}
