package tests

import (
	"image"
	"nes-emu/src/emulator/ppu/internal/domain"
)

type screenFixture struct {
	image          *image.RGBA
	ShowImageCalls int
	LastBuffer     *[domain.MAX_FRAME_SCANLINE + 1][domain.MAX_PIXEL_PER_SCANLINE + 1]uint32
}

func NewScreenFixture() *screenFixture {
	return &screenFixture{
		image: image.NewRGBA(image.Rect(0, 0, 256, 240)),
	}
}

func (s *screenFixture) ShowImage(buffer *[domain.MAX_FRAME_SCANLINE + 1][domain.MAX_PIXEL_PER_SCANLINE + 1]uint32) {
	s.ShowImageCalls++
	s.LastBuffer = buffer
}
