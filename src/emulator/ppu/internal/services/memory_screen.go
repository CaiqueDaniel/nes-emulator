package services

import (
	"nes-emu/src/emulator/ppu/internal/domain"
)

type MemoryScreen struct {
	buffer *domain.PPUFrameBuffer
}

func NewMemoryScreen() *MemoryScreen {
	return &MemoryScreen{}
}

func (s *MemoryScreen) ShowImage(buffer *domain.PPUFrameBuffer) {
	s.buffer = buffer
}

func (s *MemoryScreen) GetImageBuffer() *domain.PPUFrameBuffer {
	return s.buffer
}
