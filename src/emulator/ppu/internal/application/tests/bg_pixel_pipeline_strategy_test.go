package tests

import (
	"nes-emu/src/emulator/memory"
	"nes-emu/src/emulator/ppu/internal/application"
	"nes-emu/src/emulator/ppu/internal/domain"
	"testing"
)

func TestFetchNameTableAndAttributeTable(t *testing.T) {
	vram, _, sut, state := createBGPixelPipelineStrategySut()

	writeChar1ToPatternTable(vram)

	vram.Write(vram_nametable_address, 0x1)
	vram.Write(initial_palette_address, 0x01)
	vram.Write(initial_palette_address+1, 0x01)
	vram.Write(initial_palette_address+2, 0x09)
	vram.Write(initial_palette_address+3, 0x1A)

	tileIndex := sut.GetTileIndex(state)

	if tileIndex != 0x1 {
		t.Errorf("Expected tile index to be 0x1 (at $2000), that is, to be at the start of nametable given v=0, got %d", tileIndex)
	}

	attribute := sut.GetAttribute(state)

	if attribute != 0x00 {
		t.Errorf("Expected attribute to be 0x00, got %d", attribute)
	}
}

func TestGetLowAndHighPatternAddress(t *testing.T) {
	const expectedLowPatternAddress = uint16(0b1_0000_0001_0101)
	const expectedHighPatternAddress = uint16(0b1_0000_0001_0101) + 8
	vram, ram, sut, state := createBGPixelPipelineStrategySut()

	writeChar1ToPatternTable(vram)

	vram.Write(vram_nametable_address, 0x1)
	ram.Write(0x2000, 0b1111_0000)           //VPHB SINN
	state.SetVDirectly(0b101_00_00000_00000) //yyy NN YYYYY XXXXX

	for range 3 {
		sut.FetchPixel(state)
		state.MoveToNextDot()
	}

	lowPatternAddress := sut.GetLowPatternAddress(state)

	if lowPatternAddress != expectedLowPatternAddress {
		t.Errorf("Expected low pattern address to be %X; got %X", expectedLowPatternAddress, lowPatternAddress)
	}

	highPatternAddress := sut.GetHighPatternAddress(state)

	if highPatternAddress != expectedHighPatternAddress {
		t.Errorf("Expected high pattern address to be %X; got %X", expectedHighPatternAddress, highPatternAddress)
	}
}

func createBGPixelPipelineStrategySut() (*memory.Memory, *memory.Memory, *application.BGPixelPipelineStrategy, *domain.PPU) {
	vram := memory.NewMemory(0xFFFF)
	ram := memory.NewMemory(0xFFFF)
	sut := application.NewBGPixelPipelineStrategy(vram, ram)
	state := domain.NewPPU()
	return vram, ram, sut, state
}
