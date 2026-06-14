package tests

import (
	"nes-emu/src/components/cpu/internal"
	"nes-emu/src/components/memory"
	"testing"
)

func TestJump(t *testing.T) {
	mem := memory.NewMemory()
	cpu := internal.NewCpu(mem)

	cpu.JumpProgramCounterToValue(0x1234)

	if cpu.GetProgramCounter() != 0x1234 {
		t.Errorf("Expected program counter to be 0x1234, got %d", cpu.GetProgramCounter())
	}
}

func TestJumpProgramCounterByIndirectValue(t *testing.T) {
	mem := memory.NewMemory()
	cpu := internal.NewCpu(mem)

	mem.Write(0x03FF, 1)
	mem.Write(0x4000, 123)

	cpu.JumpProgramCounterByIndirectValue(0x03FF)

	if cpu.GetProgramCounter() != 0x0300 {
		t.Errorf("Expected instruction to corretcly imitate a hardware bug in the 6502 processor where the higher byte is read twice, to have the program counter to be 255, got %d", cpu.GetProgramCounter())
	}
}

func TestJumpProgramCounterToSubRoutine(t *testing.T) {
	mem := memory.NewMemory()
	cpu := internal.NewCpuWithProgramCounter(mem, 0x1234)

	cpu.JumpProgramCounterToSubRoutine(0x5678)

	if cpu.GetProgramCounter() != 0x5678 {
		t.Errorf("Expected program counter to be 0x5678, got %d", cpu.GetProgramCounter())
	}

	if cpu.PullValueFromStack() != 0x34 {
		t.Errorf("Expected high address to be 0x12, got %d", cpu.PullValueFromStack())
	}

	if cpu.PullValueFromStack() != 0x12 {
		t.Errorf("Expected low address to be 0x34, got %d", cpu.PullValueFromStack())
	}
}
