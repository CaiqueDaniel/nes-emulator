package application

import "nes-emu/src/emulator/ppu/internal/domain"

type DefaultPixelPipelineContext struct {
	bgPipeline PixelPipelineStrategy
}

func NewDefaultPixelPipelineContext(bgPipeline PixelPipelineStrategy) *DefaultPixelPipelineContext {
	return &DefaultPixelPipelineContext{
		bgPipeline: bgPipeline,
	}
}

func (ctx *DefaultPixelPipelineContext) FullfillFrameBuffer(ppu *domain.PPU, pipeline uint) {
	switch pipeline {
	case BG:
		ppu.PushPixelColorToFrameBuffer(ctx.bgPipeline.FetchPixel(ppu))
		return
	case OAM:
		break
	}
}
