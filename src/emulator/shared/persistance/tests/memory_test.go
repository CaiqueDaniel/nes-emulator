package tests

import (
	"nes-emu/src/emulator/components/memory"
	"testing"
)

func TestReadEmptyValueFromMemory(t *testing.T) {
	sut := memory.NewMemory()
	result := sut.Read(0x1234)

	if result != 0x00 {
		t.Errorf("Expected to read 0x00, got %d", result)
	}
}

func TestWriteToMemory(t *testing.T) {
	sut := memory.NewMemory()
	sut.Write(0x1234, 0x42)
	result := sut.Read(0x1234)

	if result != 0x42 {
		t.Errorf("Expected to read 0x42, got %d", result)
	}
}
