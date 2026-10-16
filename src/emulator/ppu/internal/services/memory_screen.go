package services

import (
	"nes-emu/src/emulator/ppu/internal/domain"
)

type MemoryScreen struct {
	buffer *[domain.MAX_VISIBLE_SCANLINES][domain.MAX_VISIBLE_DOTS]uint32
}

func NewMemoryScreen() *MemoryScreen {
	return &MemoryScreen{}
}

func (s *MemoryScreen) ShowImage(buffer *[domain.MAX_VISIBLE_SCANLINES][domain.MAX_VISIBLE_DOTS]uint32) {
	s.buffer = buffer
}

func (s *MemoryScreen) GetImageBuffer() *[domain.MAX_VISIBLE_SCANLINES][domain.MAX_VISIBLE_DOTS]uint32 {
	return s.buffer
}
