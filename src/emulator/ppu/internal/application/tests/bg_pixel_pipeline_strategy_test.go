package tests

import (
	"nes-emu/src/emulator/memory"
	"nes-emu/src/emulator/ppu/internal/application"
	"nes-emu/src/emulator/ppu/internal/domain"
	"testing"
)

func Test(t *testing.T) {
	vram := memory.NewMemory(0xFFFF)
	sut := application.NewBGPixelPipelineStrategy(vram)
	state := domain.NewPPU()

	writeChar1ToPatternTable(vram)

	vram.Write(vram_nametable_address, 0x1)
	vram.Write(initial_palette_address, 0x01)
	vram.Write(initial_palette_address+1, 0x01)
	vram.Write(initial_palette_address+2, 0x09)
	vram.Write(initial_palette_address+3, 0x1A)

	tileIndex := sut.GetTileIndex(*state)

	if tileIndex != 0x1 {
		t.Errorf("Expected tile index to be 0x1 (at $2000), that is, to be at the start of nametable given v=0, got %d", tileIndex)
	}

	attribute := sut.GetAttribute(*state)

	if attribute != 0x00 {
		t.Errorf("Expected attribute to be 0x00, got %d", attribute)
	}
}
