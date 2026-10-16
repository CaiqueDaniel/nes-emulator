package tests

import (
	"nes-emu/src/emulator/memory"
	"testing"
)

func TestMemoryWriteAndRead(t *testing.T) {
	sut := createSut()

	sut.Write(0x0000, 0xAB)
	value := sut.Read(0x0000)

	if value != 0xAB {
		t.Errorf("Expected value 0xAB, but got 0x%X", value)
	}
}

func TestReadUninitializedMemory(t *testing.T) {
	sut := createSut()

	value := sut.Read(0x0001)

	if value != 0x00 {
		t.Errorf("Expected value 0x00 for uninitialized memory, but got 0x%X", value)
	}
}

func createSut() *memory.Memory {
	return memory.NewMemory(1024)
}
