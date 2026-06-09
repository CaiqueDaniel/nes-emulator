package tests

import (
	"nes-emu/src/components/cpu/internal"
	"nes-emu/src/components/memory"
	"testing"
)

// Imediate addressing mode
func TestShouldLoadAcumulator(t *testing.T) {
	memory := memory.NewMemory()
	cpu := internal.NewCpu(memory)
	cpu.LoadValueIntoRegister(0x42, internal.ACCUMULATOR)

	if cpu.GetDebugData()["acc"] != 0x42 {
		t.Errorf("Expected acc to be 0x42, got %d", cpu.GetDebugData()["acc"])
	}
}

func TestShouldLoadXRegister(t *testing.T) {
	memory := memory.NewMemory()
	cpu := internal.NewCpu(memory)
	cpu.LoadValueIntoRegister(0x42, internal.REGISTER_X)

	if cpu.GetDebugData()["x"] != 0x42 {
		t.Errorf("Expected x to be 0x42, got %d", cpu.GetDebugData()["x"])
	}
}

func TestShouldLoadYRegister(t *testing.T) {
	memory := memory.NewMemory()
	cpu := internal.NewCpu(memory)
	cpu.LoadValueIntoRegister(0x42, internal.REGISTER_Y)

	if cpu.GetDebugData()["y"] != 0x42 {
		t.Errorf("Expected y to be 0x42, got %d", cpu.GetDebugData()["y"])
	}
}
