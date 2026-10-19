package application

import "nes-emu/src/emulator/ppu/internal/domain"

const (
	BG  = 0
	OAM = 1
)

type PixelPipelineContext interface {
	FullfillFrameBuffer(ppu *domain.PPU, pipeline uint)
}

type PixelPipelineStrategy interface {
	FetchPixel(ppu domain.PPU) uint32
}
