package tests

import (
	"nes-emu/src/emulator/components/ppu"
	"testing"
)

func TestPPURender_ShouldDrawAPixel(t *testing.T) {
	ppu := ppu.NewPPUWithInternal()

	ppu.Render()

	if ppu.GetCurrentScanline() != 0 {
		t.Fatal("scanline expected to be 0")
	}

	if ppu.GetCurrentScanlinePixel() != 1 {
		t.Fatal("expected to draw a pixel")
	}
}

func TestPPURender_ShouldWrapScanlineToStart(t *testing.T) {
	ppu := ppu.NewPPUWithInternal()
	i := 0

	for {
		ppu.Render()
		i++

		if i == 256*262 {
			break
		}
	}

	if ppu.GetCurrentScanline() != 0 {
		t.Errorf("scanline expected to be 0, got %d", ppu.GetCurrentScanline())
	}

	if ppu.GetCurrentScanlinePixel() != 0 {
		t.Errorf("expected to be at pixel 0, got %d", ppu.GetCurrentScanlinePixel())
	}
}
