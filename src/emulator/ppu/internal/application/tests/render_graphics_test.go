package tests

import (
	"nes-emu/src/emulator/memory"
	"nes-emu/src/emulator/ppu/internal/application"
	"nes-emu/src/emulator/ppu/internal/domain"
	"nes-emu/src/emulator/ppu/internal/persistence"
	"nes-emu/src/emulator/ppu/internal/services"
	"slices"
	"testing"
)

const vram_nametable_address = 0x2000
const initial_palette_address = 0x3F00
const initial_pattern_address = 0x0000

func TestShouldIncrementDot(t *testing.T) {
	sut, repository, _, _ := createSut()

	sut.Execute()

	state := repository.Get()

	if state.GetDot() != 1 {
		t.Errorf("Expected dot to be 1, got %d", state.GetDot())
	}
}

func TestShouldRenderGraphics(t *testing.T) {
	sut, repository, memory, screen := createSut()

	writeChar1ToPatternTable(memory)

	memory.Write(vram_nametable_address, 0x1)
	memory.Write(initial_palette_address, 0x01)
	memory.Write(initial_palette_address+1, 0x01)
	memory.Write(initial_palette_address+2, 0x09)
	memory.Write(initial_palette_address+3, 0x1A)

	for i := 0; i < domain.MAX_VISIBLE_SCANLINES; i++ {
		for y := 0; y < domain.MAX_DOTS_PER_SCANLINE; y++ {
			sut.Execute()
		}
	}

	state := repository.Get()

	if state.GetDot() != 0 {
		t.Errorf("Expected dot 0")
	}

	if got := state.GetScanline(); got != domain.INITIAL_V_BLANK_SCANLINE-1 {
		t.Errorf("Expected v-blank, got %d", got)
	}

	if screen.GetImageBuffer() == nil {
		t.Errorf("Expected image buffer to be rendered, but it was nil")
	}

	var chunks []uint32 = screen.GetImageBuffer()[0][0:4]

	if slices.Compare(chunks, []uint32{0x0000FC, 0x0000FC, 0x0000FC, 0x007800}) != 4 {
		t.Errorf("Pixels rendered different from expected; got %X, %X, %X, %X", chunks[0], chunks[1], chunks[2], chunks[3])
	}
}

func createSut() (*application.RenderGraphics, *persistence.PPUMemoryRepository, *memory.Memory, *services.MemoryScreen) {
	vram := memory.NewMemory(0xFFFF)
	repo := persistence.NewPPUMemoryRepository()
	screen := services.NewMemoryScreen()
	return application.NewRenderGraphics(repo, vram, screen), repo, vram, screen
}
