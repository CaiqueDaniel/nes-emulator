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
