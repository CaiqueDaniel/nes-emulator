package tests

import (
	"testing"

	"nes-emu/src/emulator/ppu/internal/domain"
)

func TestNewPPUInitialState(t *testing.T) {
	ppu := domain.NewPPU()

	if ppu == nil {
		t.Fatal("NewPPU() returned nil")
	}

	if got := ppu.GetScanline(); got != 0 {
		t.Errorf("GetScanline() = %d, want 0", got)
	}

	if got := ppu.GetDot(); got != 0 {
		t.Errorf("GetDot() = %d, want 0", got)
	}

	if got := ppu.IsBackgroundRenderingEnabled(); got {
		t.Error("IsBackgroundRenderingEnabled() = true, want false")
	}

	if got := ppu.IsSpriteRenderingEnabled(); got {
		t.Error("IsSpriteRenderingEnabled() = true, want false")
	}

	if got := ppu.IsRenderingEnabled(); got {
		t.Error("IsRenderingEnabled() = true, want false")
	}
}

func TestPPUMoveToNextDot(t *testing.T) {
	ppu := domain.NewPPU()
	ppu.MoveToNextDot()

	if got := ppu.GetDot(); got != 1 {
		t.Errorf("GetDot() after MoveToNextDot() = %d, want 1", got)
	}
}

func TestPPUMoveToNextScanline(t *testing.T) {
	ppu := domain.NewPPU()

	for i := 0; i < domain.MAX_DOTS_PER_SCANLINE; i++ {
		ppu.MoveToNextDot()
	}

	if got := ppu.GetDot(); got != 0 {
		t.Errorf("GetDot() after MoveToNextDot() = %d, want 0", got)
	}

	if got := ppu.GetScanline(); got != 1 {
		t.Errorf("GetScanline() after MoveToNextDot() = %d, want 1", got)
	}
}

func TestPPUShouldWrapAroundScanlines(t *testing.T) {
	ppu := domain.NewPPU()

	for scanline := 0; scanline < domain.MAX_SCANLINES; scanline++ {
		for i := 0; i < domain.MAX_DOTS_PER_SCANLINE; i++ {
			ppu.MoveToNextDot()
		}
	}

	if got := ppu.GetDot(); got != 0 {
		t.Errorf("GetDot() after MoveToNextDot() = %d, want 0", got)
	}

	if got := ppu.GetScanline(); got != 0 {
		t.Errorf("GetScanline() after MoveToNextDot() = %d, want 0", got)
	}
}
