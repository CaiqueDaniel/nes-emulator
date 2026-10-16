package application

import (
	"nes-emu/src/emulator/ppu/internal/domain"
	"nes-emu/src/emulator/shared"
)

type RenderGraphics struct {
	repository domain.PPURepository
	vram       shared.Memory
	screen     Screen
}

func NewRenderGraphics(repository domain.PPURepository, vram shared.Memory, screen Screen) *RenderGraphics {
	return &RenderGraphics{
		repository: repository,
		vram:       vram,
		screen:     screen,
	}
}

func (rg *RenderGraphics) Execute() {
	//TODO: Implement the logic to render graphics based on the PPU state and VRAM content.
}
