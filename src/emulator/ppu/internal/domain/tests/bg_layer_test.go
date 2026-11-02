package tests

import (
	"nes-emu/src/emulator/ppu/internal/domain"
	"testing"
)

func TestBGLayer_LoadShiftRegistersFromLatches(t *testing.T) {
	bglayer := domain.NewBGLayer()
	bglayer.SetLowPatternLatch(0x12)
	bglayer.SetHighPatternLatch(0x34)
	bglayer.SetTileAttributeIndex(0b11)

	bglayer.LoadShiftRegistersFromLatches()

	if bglayer.GetLowPatternShiftRegister() != 0x12 {
		t.Errorf("Failed to load shift registers from latches")
	}

	if bglayer.GetHighPatternShiftRegister() != 0x34 {
		t.Errorf("Failed to load shift registers from latches")
	}

	if bglayer.GetLowAttributeShiftRegister() != 255 {
		t.Errorf("Failed to load shift registers from latches")
	}

	if bglayer.GetHighAttributeShiftRegister() != 255 {
		t.Errorf("Failed to load shift registers from latches")
	}
}

func TestBGLayer_ShiftRegisters(t *testing.T) {
	bglayer := domain.NewBGLayer()
	bglayer.SetLowPatternLatch(0x12)
	bglayer.SetHighPatternLatch(0x34)
	bglayer.SetTileAttributeIndex(0b11)

	bglayer.LoadShiftRegistersFromLatches()
	bglayer.ShiftRegisters()

	if bglayer.GetLowPatternShiftRegister() != (0x12 << 1) {
		t.Errorf("Failed to load shift registers from latches")
	}

	if bglayer.GetHighPatternShiftRegister() != (0x34 << 1) {
		t.Errorf("Failed to load shift registers from latches")
	}

	if bglayer.GetLowAttributeShiftRegister() != (0xFF << 1) {
		t.Errorf("Failed to load shift registers from latches")
	}

	if bglayer.GetHighAttributeShiftRegister() != (0xFF << 1) {
		t.Errorf("Failed to load shift registers from latches")
	}
}

func TestBGLayer_LoadShiftRegistersFromLatches_ShouldUpdateFirst8BitsFromRegister(t *testing.T) {
	bglayer := domain.NewBGLayer()
	bglayer.SetLowPatternLatch(0x12)
	bglayer.SetHighPatternLatch(0x34)
	bglayer.SetTileAttributeIndex(0b11)

	bglayer.LoadShiftRegistersFromLatches()

	for range 8 {
		bglayer.ShiftRegisters()
	}

	bglayer.LoadShiftRegistersFromLatches()

	if bglayer.GetLowPatternShiftRegister() != (0x12<<8)|0x12 {
		t.Errorf("Failed to load shift registers from latches")
	}

	if bglayer.GetHighPatternShiftRegister() != (0x34<<8)|0x34 {
		t.Errorf("Failed to load shift registers from latches")
	}

	if bglayer.GetLowAttributeShiftRegister() != (0xFF<<8)|0xFF {
		t.Errorf("Failed to load shift registers from latches")
	}

	if bglayer.GetHighAttributeShiftRegister() != (0xFF<<8)|0xFF {
		t.Errorf("Failed to load shift registers from latches")
	}
}
