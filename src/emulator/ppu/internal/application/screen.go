package application

import "nes-emu/src/emulator/ppu/internal/domain"

type Screen interface {
	ShowImage(buffer *[domain.MAX_FRAME_SCANLINE + 1][domain.MAX_PIXEL_PER_SCANLINE + 1]uint32)
}
