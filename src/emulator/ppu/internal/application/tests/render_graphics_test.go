package tests

import (
	"nes-emu/src/emulator/memory"
	"nes-emu/src/emulator/ppu/internal/application"
	"nes-emu/src/emulator/ppu/internal/persistence"
	"nes-emu/src/emulator/ppu/internal/services"
	"slices"
	"testing"
)

const vram_nametable_address = 0x2000
const initial_palette_address = 0x3F00
const initial_pattern_address = 0x0000

func TestShouldRenderGraphics(t *testing.T) {
	sut, _, memory, screen := createSut()

	writeChar1ToPatternTable(memory)

	memory.Write(vram_nametable_address, 0x1)
	memory.Write(initial_palette_address, 0x01)
	memory.Write(initial_palette_address+1, 0x01)
	memory.Write(initial_palette_address+2, 0x09)
	memory.Write(initial_palette_address+3, 0x1A)

	sut.Execute()

	if screen.GetImageBuffer() == nil {
		t.Errorf("Expected image buffer to be rendered, but it was nil")
	}

	var chunks []uint32 = screen.GetImageBuffer()[0][0:4]

	if slices.Compare(chunks, []uint32{0x0000FC, 0x0000FC, 0x0000FC, 0x007800}) != 4 {
		t.Errorf("Pixels rendered different from expected")
	}
}

func createSut() (*application.RenderGraphics, *persistence.PPUMemoryRepository, *memory.Memory, *services.MemoryScreen) {
	vram := memory.NewMemory(0xFFFF)
	repo := persistence.NewPPUMemoryRepository()
	screen := services.NewMemoryScreen()
	return application.NewRenderGraphics(repo, vram, screen), repo, vram, screen
}

func writeChar1ToPatternTable(memory *memory.Memory) {
	patternTable := []byte{
		0x10, // Linha 0: ...#....
		0x30, // Linha 1: ..##....
		0x10, // Linha 2: ...#....
		0x10, // Linha 3: ...#....
		0x10, // Linha 4: ...#....
		0x10, // Linha 5: ...#....
		0x70, // Linha 6: .#####..
		0x00, // Linha 7: ........
		0x10, // (Plano Alto - High Byte) correspondente às mesmas linhas
		0x30,
		0x10,
		0x10,
		0x10,
		0x10,
		0x70,
		0x00,
	}

	for i, b := range patternTable {
		memory.Write(uint16(initial_pattern_address)+uint16(i), b)
	}
}
