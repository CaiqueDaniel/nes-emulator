package application

import (
	"nes-emu/src/emulator/ppu/internal/domain"
	"nes-emu/src/emulator/shared"
)

type RenderGraphics struct {
	repository           domain.PPURepository
	vram                 shared.Memory
	screen               Screen
	pixelPipelineContext PixelPipelineContext
}

func NewRenderGraphics(
	repository domain.PPURepository,
	vram shared.Memory,
	screen Screen,
	pixelPipelineContext PixelPipelineContext,
) *RenderGraphics {
	return &RenderGraphics{
		repository:           repository,
		vram:                 vram,
		screen:               screen,
		pixelPipelineContext: pixelPipelineContext,
	}
}

func (rg *RenderGraphics) Execute() {
	state := rg.repository.Get()

	state.MoveToNextDot()
	rg.pixelPipelineContext.FullfillFrameBuffer(&state, BG_PIPELINE)

	if state.IsVBlankStarted() {
		rg.screen.ShowImage(state.GetFrameBuffer())
	}

	rg.repository.Save(state)
}
