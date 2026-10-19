package application

import "nes-emu/src/emulator/ppu/internal/domain"

type Screen interface {
	ShowImage(buffer *domain.PPUFrameBuffer)
}
