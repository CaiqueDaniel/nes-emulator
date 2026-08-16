package tests

import "nes-emu/src/emulator/application"

type PixelPipelineFixture struct {
}

func (p *PixelPipelineFixture) StepUpPipeline(currentDot uint, vValue uint16, fineY uint16) *application.PipelineResult {
	if currentDot%8 == 0 {
		return &application.PipelineResult{
			HighPalette:     0b00000001,
			LowPalette:      0b00000001,
			LowPatternByte:  0b00000010,
			HighPatternByte: 0b00000011,
		}
	}

	return nil
}
