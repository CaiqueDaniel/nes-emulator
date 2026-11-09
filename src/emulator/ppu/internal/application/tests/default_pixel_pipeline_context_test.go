package tests

import (
	"nes-emu/src/emulator/ppu/internal/application"
	"nes-emu/src/emulator/ppu/internal/domain"
	"testing"
)

func TestShouldPushBgPixelColor(t *testing.T) {
	state := domain.NewPPU()
	sut := createDefaultPixelPipelineContext()
	sut.FullfillFrameBuffer(state, application.BG)

	if state.GetFrameBuffer()[0][0] != 0xFC {
		t.Errorf("Expected background pixel color to be 0xFC, got %x", state.GetFrameBuffer()[0][0])
	}
}

// TODO: Not implemented
func TestShouldPushOamPixelColor(t *testing.T) {
	state := domain.NewPPU()
	sut := createDefaultPixelPipelineContext()
	sut.FullfillFrameBuffer(state, application.OAM)

	if state.GetFrameBuffer()[0][0] != 0x0 {
		t.Errorf("Expected object pixel color to be 0xFC, got %x", state.GetFrameBuffer()[0][0])
	}
}

func createDefaultPixelPipelineContext() *application.DefaultPixelPipelineContext {
	stub := &StubBGPipelineStrategy{}
	return application.NewDefaultPixelPipelineContext(stub)
}

type StubBGPipelineStrategy struct{}

func (stub *StubBGPipelineStrategy) FetchPixel(ppu *domain.PPU) uint32 {
	return 0xFC
}
