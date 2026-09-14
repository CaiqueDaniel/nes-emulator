package application

type PixelPipeline interface {
	StepUpPipeline(currentDot uint, vValue uint16, fineY uint16) bool
	RenderPixel(fineX byte) uint32
}
