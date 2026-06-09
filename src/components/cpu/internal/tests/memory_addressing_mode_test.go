package tests

import (
	"testing"

	"nes-emu/src/components/cpu/internal"
	"nes-emu/src/components/memory"
)

func TestGetValueByAbsoluteMode(t *testing.T) {
	memory := memory.NewMemory()
	cpu := internal.NewCpu(memory)

	memory.Write(0x1000, 0xFF)

	value := cpu.GetValueByAbsoluteMode(0x1000)

	if value != 0xFF {
		t.Errorf("Expected %d, got %d", 0xFF, value)
	}
}
