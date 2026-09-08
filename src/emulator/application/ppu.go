package application

type PPU interface {
	Render()
	TriggerLatchWithWriteSignal(address uint16)
	TriggerLatchWithReadSignal(address uint16)
}

type PixelPipeline interface {
	StepUpPipeline(currentDot uint, vValue uint16, fineY uint16) bool
	RenderPixel(fineX byte) uint32
}
