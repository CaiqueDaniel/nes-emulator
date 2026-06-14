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

	cpu.JumpProgramCounterByIndirectValue(0x03FF)

	if cpu.GetProgramCounter() != 0x0300 {
		t.Errorf("Expected instruction to corretcly imitate a hardware bug in the 6502 processor where the higher byte is read twice, to have the program counter to be 255, got %d", cpu.GetProgramCounter())
	}
}
