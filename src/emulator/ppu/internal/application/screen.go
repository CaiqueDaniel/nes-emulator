package application

import "nes-emu/src/emulator/ppu/internal/domain"

type Screen interface {
	ShowImage(buffer *[domain.MAX_VISIBLE_SCANLINES][domain.MAX_VISIBLE_DOTS]uint32)
}
