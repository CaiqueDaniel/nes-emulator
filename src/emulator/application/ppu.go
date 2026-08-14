package application

type PPU interface {
	Render()
}

type PixelPipeline interface {
	StepUpPipeline(currentDot uint, vValue uint16, fineY uint16) *PipelineResult
}

type PipelineResult struct {
	HighPalette     byte
	LowPalette      byte
	LowPatternByte  byte
	HighPatternByte byte
}
