package tests

import (
	"nes-emu/src/components/cpu/internal"
	"nes-emu/src/components/memory"
	"testing"
)

func TestBranchIfCarryIsClearAndValuePositive(t *testing.T) {
	memory := memory.NewMemory()
	cpu := internal.NewCpu(memory)

	cpu.BranchIfCarryIsClear(0x10)

	if cpu.GetProgramCounter() != 0x12 {
		t.Errorf("Expected pc to be 0x%04X, got 0x%04X", 0x12, cpu.GetProgramCounter())
	}
}

func TestBranchIfCarryIsClearAndValueNegative(t *testing.T) {
	memory := memory.NewMemory()
	cpu := internal.NewCpuWithProgramCounter(memory, 0xA)

	cpu.BranchIfCarryIsClear(128)

	if cpu.GetProgramCounter() != 0xB {
		t.Errorf("Expected pc to be 0x%04X, got 0x%04X", 0xB, cpu.GetProgramCounter())
	}
}

func TestBranchIfCarryIsClear_CarryIsSet(t *testing.T) {
	memory := memory.NewMemory()
	cpu := internal.NewCpuWithProgramCounter(memory, 0xA)

	cpu.AddWithCarry(0xFF)
	cpu.AddWithCarry(1)
	cpu.BranchIfCarryIsClear(10)

	if cpu.GetProgramCounter() != 0xA {
		t.Errorf("Expected pc to be 0x%04X, got 0x%04X", 0xA, cpu.GetProgramCounter())
	}
}

func TestBranchIfCarryIsSetAndValuePositive(t *testing.T) {
	memory := memory.NewMemory()
	cpu := internal.NewCpu(memory)

	cpu.AddWithCarry(255)
	cpu.AddWithCarry(1)
	cpu.BranchIfCarryIsSet(0x10)

	if cpu.GetProgramCounter() != 0x12 {
		t.Errorf("Expected pc to be 0x%04X, got 0x%04X", 0x12, cpu.GetProgramCounter())
	}
}

func TestBranchIfCarryIsSetAndValueNegative(t *testing.T) {
	memory := memory.NewMemory()
	cpu := internal.NewCpuWithProgramCounter(memory, 0xA)

	cpu.AddWithCarry(255)
	cpu.AddWithCarry(1)
	cpu.BranchIfCarryIsSet(128)

	if cpu.GetProgramCounter() != 0xB {
		t.Errorf("Expected pc to be 0x%04X, got 0x%04X", 0xB, cpu.GetProgramCounter())
	}
}

func TestBranchIfCarryIsSet_CarryIsClear(t *testing.T) {
	memory := memory.NewMemory()
	cpu := internal.NewCpuWithProgramCounter(memory, 0xA)
	cpu.BranchIfCarryIsSet(10)

	if cpu.GetProgramCounter() != 0xA {
		t.Errorf("Expected pc to be 0x%04X, got 0x%04X", 0xA, cpu.GetProgramCounter())
	}
}

func TestBranchIfNegative(t *testing.T) {

}

func TestBranchIfZero(t *testing.T) {

}

func TestBranchIfNotZero(t *testing.T) {

}

func TestBranchIfOverflow(t *testing.T) {

}
