package tests

import (
	"nes-emu/src/emulator/components/bus"
	"nes-emu/src/emulator/components/memory"
	"testing"
)

func TestReadEmptyValueFromMemory(t *testing.T) {
	b := bus.NewBusWithInternalType()
	sut := memory.NewMemory(b)
	result := sut.Read(0x1234)
	tickCount := b.GetTickCount()

	if result != 0x00 {
		t.Errorf("Expected to read 0x00, got %d", result)
	}

	if tickCount != 1 {
		t.Errorf("Expected tick count to be 1, got %d", tickCount)
	}
}

func TestWriteToMemory(t *testing.T) {
	b := bus.NewBusWithInternalType()
	sut := memory.NewMemory(b)
	sut.Write(0x1234, 0x42)
	result := sut.Read(0x1234)
	tickCount := b.GetTickCount()

	if result != 0x42 {
		t.Errorf("Expected to read 0x42, got %d", result)
	}

	if tickCount != 2 {
		t.Errorf("Expected tick count to be 1, got %d", tickCount)
	}
}
