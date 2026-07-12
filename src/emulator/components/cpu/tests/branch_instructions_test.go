package tests

import (
	"nes-emu/src/emulator/components/bus"
	internal "nes-emu/src/emulator/components/cpu"
	"nes-emu/src/emulator/components/memory"
	"testing"
)

func TestBranchIfCarryIsClearAndValuePositive(t *testing.T) {
	memory := memory.NewMemory()
	b := bus.NewBus(memory)
	cpu := internal.NewCpuWithInternal(b)

	cpu.BranchIfCarryIsClear(0x10)

	if cpu.GetProgramCounter() != 0x12 {
		t.Errorf("Expected pc to be 0x%04X, got 0x%04X", 0x12, cpu.GetProgramCounter())
	}
}

func TestBranchIfCarryIsClearAndValueNegative(t *testing.T) {
	memory := memory.NewMemory()
	b := bus.NewBus(memory)
	cpu := internal.NewCpuWithProgramCounter(0xA, b)

	cpu.BranchIfCarryIsClear(128)

	if cpu.GetProgramCounter() != 0xB {
		t.Errorf("Expected pc to be 0x%04X, got 0x%04X", 0xB, cpu.GetProgramCounter())
	}
}

func TestBranchIfCarryIsClear_CarryIsSet(t *testing.T) {
	memory := memory.NewMemory()
	b := bus.NewBus(memory)
	cpu := internal.NewCpuWithProgramCounter(0xA, b)

	cpu.AddWithCarry(0xFF)
	cpu.AddWithCarry(1)
	cpu.BranchIfCarryIsClear(10)

	if cpu.GetProgramCounter() != 0xA {
		t.Errorf("Expected pc to be 0x%04X, got 0x%04X", 0xA, cpu.GetProgramCounter())
	}
}

func TestBranchIfCarryIsSetAndValuePositive(t *testing.T) {
	memory := memory.NewMemory()
	b := bus.NewBus(memory)
	cpu := internal.NewCpuWithInternal(b)

	cpu.AddWithCarry(255)
	cpu.AddWithCarry(1)
	cpu.BranchIfCarryIsSet(0x10)

	if cpu.GetProgramCounter() != 0x12 {
		t.Errorf("Expected pc to be 0x%04X, got 0x%04X", 0x12, cpu.GetProgramCounter())
	}
}

func TestBranchIfCarryIsSetAndValueNegative(t *testing.T) {
	memory := memory.NewMemory()
	b := bus.NewBus(memory)
	cpu := internal.NewCpuWithProgramCounter(0xA, b)

	cpu.AddWithCarry(255)
	cpu.AddWithCarry(1)
	cpu.BranchIfCarryIsSet(128)

	if cpu.GetProgramCounter() != 0xB {
		t.Errorf("Expected pc to be 0x%04X, got 0x%04X", 0xB, cpu.GetProgramCounter())
	}
}

func TestBranchIfCarryIsSet_CarryIsClear(t *testing.T) {
	memory := memory.NewMemory()
	b := bus.NewBus(memory)
	cpu := internal.NewCpuWithProgramCounter(0xA, b)

	cpu.BranchIfCarryIsSet(10)

	if cpu.GetProgramCounter() != 0xA {
		t.Errorf("Expected pc to be 0x%04X, got 0x%04X", 0xA, cpu.GetProgramCounter())
	}
}

func TestBranchIfEqualWithValuePositive(t *testing.T) {
	memory := memory.NewMemory()
	b := bus.NewBus(memory)
	cpu := internal.NewCpuWithInternal(b)

	cpu.AddWithCarry(0)
	cpu.BranchIfEqual(0x10)

	if cpu.GetProgramCounter() != 0x12 {
		t.Errorf("Expected pc to be 0x%04X, got 0x%04X", 0x12, cpu.GetProgramCounter())
	}
}

func TestBranchIfEqualWithValueNegative(t *testing.T) {
	memory := memory.NewMemory()
	b := bus.NewBus(memory)
	cpu := internal.NewCpuWithProgramCounter(0xA, b)

	cpu.AddWithCarry(0)
	cpu.BranchIfEqual(128)

	if cpu.GetProgramCounter() != 0xB {
		t.Errorf("Expected pc to be 0x%04X, got 0x%04X", 0xB, cpu.GetProgramCounter())
	}
}

func TestBranchIfEqual_ZeroIsNotClear(t *testing.T) {
	memory := memory.NewMemory()
	b := bus.NewBus(memory)
	cpu := internal.NewCpuWithProgramCounter(0xA, b)

	cpu.BranchIfEqual(10)

	if cpu.GetProgramCounter() != 0xA {
		t.Errorf("Expected pc to be 0x%04X, got 0x%04X", 0xA, cpu.GetProgramCounter())
	}
}

func TestBranchIfEqual_ZeroIsClear(t *testing.T) {
	memory := memory.NewMemory()
	b := bus.NewBus(memory)
	cpu := internal.NewCpuWithProgramCounter(0xA, b)

	cpu.AddWithCarry(1)
	cpu.BranchIfEqual(10)

	if cpu.GetProgramCounter() != 0xA {
		t.Errorf("Expected pc to be 0x%04X, got 0x%04X", 0xA, cpu.GetProgramCounter())
	}
}

func TestBranchIfPositiveWithValuePositive(t *testing.T) {
	memory := memory.NewMemory()
	b := bus.NewBus(memory)
	cpu := internal.NewCpuWithProgramCounter(0xA, b)

	cpu.AddWithCarry(1)
	cpu.BranchIfPositive(10)

	if cpu.GetProgramCounter() != 0x16 {
		t.Errorf("Expected pc to be 0x%04X, got 0x%04X", 0x16, cpu.GetProgramCounter())
	}
}

func TestBranchIfPositiveWithValueNegative(t *testing.T) {
	memory := memory.NewMemory()
	b := bus.NewBus(memory)
	cpu := internal.NewCpuWithProgramCounter(0xA, b)

	cpu.AddWithCarry(128)
	cpu.BranchIfPositive(128)

	if cpu.GetProgramCounter() != 0xA {
		t.Errorf("Expected pc to be 0x%04X, got 0x%04X", 0xA, cpu.GetProgramCounter())
	}
}

func TestBranchIfNegativeWithValuePositive(t *testing.T) {
	memory := memory.NewMemory()
	b := bus.NewBus(memory)
	cpu := internal.NewCpuWithProgramCounter(0xA, b)

	cpu.AddWithCarry(1)
	cpu.BranchIfNegative(10)

	if cpu.GetProgramCounter() != 0xA {
		t.Errorf("Expected pc to be 0x%04X, got 0x%04X", 0xA, cpu.GetProgramCounter())
	}
}

func TestBranchIfNegativeWithValueNegative(t *testing.T) {
	memory := memory.NewMemory()
	b := bus.NewBus(memory)
	cpu := internal.NewCpuWithProgramCounter(0xA, b)

	cpu.AddWithCarry(128)
	cpu.BranchIfNegative(128)

	if cpu.GetProgramCounter() != 0xB {
		t.Errorf("Expected pc to be 0x%04X, got 0x%04X", 0xB, cpu.GetProgramCounter())
	}
}

func TestBranchIfNotEqual_ZeroIsNotClear(t *testing.T) {
	memory := memory.NewMemory()
	b := bus.NewBus(memory)
	cpu := internal.NewCpuWithProgramCounter(0xA, b)

	cpu.AddWithCarry(1)
	cpu.BranchIfNotEqual(10)

	if cpu.GetProgramCounter() != 0x16 {
		t.Errorf("Expected pc to be 0x%04X, got 0x%04X", 0x16, cpu.GetProgramCounter())
	}
}

func TestBranchIfNotEqual_ZeroIsClear(t *testing.T) {
	memory := memory.NewMemory()
	b := bus.NewBus(memory)
	cpu := internal.NewCpuWithProgramCounter(0xA, b)

	cpu.AddWithCarry(0)
	cpu.BranchIfNotEqual(10)

	if cpu.GetProgramCounter() != 0xA {
		t.Errorf("Expected pc to be 0x%04X, got 0x%04X", 0xA, cpu.GetProgramCounter())
	}
}

func TestBranchIfOverflowClear_OverflowIsSet(t *testing.T) {
	memory := memory.NewMemory()
	b := bus.NewBus(memory)
	cpu := internal.NewCpuWithProgramCounter(0xA, b)

	cpu.AddWithCarry(127)
	cpu.AddWithCarry(1)
	cpu.BranchIfOverflowClear(10)

	if cpu.GetProgramCounter() != 0xA {
		t.Errorf("Expected pc to be 0x%04X, got 0x%04X", 0xA, cpu.GetProgramCounter())
	}
}

func TestBranchIfOverflowClear_OverflowIsClear(t *testing.T) {
	memory := memory.NewMemory()
	b := bus.NewBus(memory)
	cpu := internal.NewCpuWithProgramCounter(0xA, b)

	cpu.AddWithCarry(1)
	cpu.BranchIfOverflowClear(10)

	if cpu.GetProgramCounter() != 0x16 {
		t.Errorf("Expected pc to be 0x%04X, got 0x%04X", 0x16, cpu.GetProgramCounter())
	}
}

func TestBranchIfOverflowSet_OverflowIsClear(t *testing.T) {
	memory := memory.NewMemory()
	b := bus.NewBus(memory)
	cpu := internal.NewCpuWithProgramCounter(0xA, b)

	cpu.AddWithCarry(1)
	cpu.BranchIfOverflowSet(10)

	if cpu.GetProgramCounter() != 0xA {
		t.Errorf("Expected pc to be 0x%04X, got 0x%04X", 0xA, cpu.GetProgramCounter())
	}
}

func TestBranchIfOverflowSet_OverflowIsSet(t *testing.T) {
	memory := memory.NewMemory()
	b := bus.NewBus(memory)
	cpu := internal.NewCpuWithProgramCounter(0xA, b)

	cpu.AddWithCarry(127)
	cpu.AddWithCarry(1)
	cpu.BranchIfOverflowSet(10)

	if cpu.GetProgramCounter() != 0x16 {
		t.Errorf("Expected pc to be 0x%04X, got 0x%04X", 0x16, cpu.GetProgramCounter())
	}
}
