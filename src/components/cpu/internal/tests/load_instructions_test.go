package tests

import (
	"nes-emu/src/components/cpu/internal"
	"nes-emu/src/components/memory"
	"testing"
)

func TestShouldLoadAcumulator(t *testing.T) {
	memory := memory.NewMemory()
	cpu := internal.NewCpu(memory)
	cpu.LoadValueIntoRegister(0x42, internal.ACCUMULATOR)

	if cpu.GetDebugData()["acc"] != 0x42 {
		t.Errorf("Expected acc to be 0x42, got %d", cpu.GetDebugData()["acc"])
	}

	if cpu.GetNegativeFlag() {
		t.Errorf("Expected negative flag to be false")
	}

	if cpu.GetZeroFlag() {
		t.Errorf("Expected zero flag to be false")
	}
}

func TestShouldLoadAcumulatorWithZero(t *testing.T) {
	memory := memory.NewMemory()
	cpu := internal.NewCpu(memory)
	cpu.LoadValueIntoRegister(0x00, internal.ACCUMULATOR)

	if cpu.GetDebugData()["acc"] != 0x00 {
		t.Errorf("Expected acc to be 0x00, got %d", cpu.GetDebugData()["acc"])
	}

	if !cpu.GetZeroFlag() {
		t.Errorf("Expected zero flag to be true")
	}

	if cpu.GetNegativeFlag() {
		t.Errorf("Expected negative flag to be false")
	}
}

func TestShouldLoadAcumulatorWithNegativeValue(t *testing.T) {
	memory := memory.NewMemory()
	cpu := internal.NewCpu(memory)
	cpu.LoadValueIntoRegister(128, internal.ACCUMULATOR)

	if cpu.GetDebugData()["acc"] != 128 {
		t.Errorf("Expected acc to be 128, got %d", cpu.GetDebugData()["acc"])
	}

	if !cpu.GetNegativeFlag() {
		t.Errorf("Expected negative flag to be true")
	}

	if cpu.GetZeroFlag() {
		t.Errorf("Expected zero flag to be false")
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
