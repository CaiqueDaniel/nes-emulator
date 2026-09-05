package application

type PPU interface {
	Render()
}

type PixelPipeline interface {
	StepUpPipeline(currentDot uint, vValue uint16, fineY uint16) bool
	RenderPixel(fineX byte) uint32
}

type PipelineResult struct {
	HighPalette     byte
	LowPalette      byte
	LowPatternByte  byte
	HighPatternByte byte
}
