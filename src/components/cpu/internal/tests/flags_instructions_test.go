package tests

import (
	"nes-emu/src/components/cpu/internal"
	"nes-emu/src/components/memory"
	"testing"
)

func TestCarryFlagInstructions(t *testing.T) {
	mem := memory.NewMemory()
	cpu := internal.NewCpu(mem)

	// Default should be false
	if cpu.GetCarryFlag() {
		t.Error("Expected carry flag to be false initially")
	}

	// Set Carry Flag
	cpu.SetCarryFlag()
	if !cpu.GetCarryFlag() {
		t.Error("Expected carry flag to be true after SetCarryFlag")
	}

	// Clear Carry Flag
	cpu.ClearCarryFlag()
	if cpu.GetCarryFlag() {
		t.Error("Expected carry flag to be false after ClearCarryFlag")
	}
}

func TestInterruptFlagInstructions(t *testing.T) {
	mem := memory.NewMemory()
	cpu := internal.NewCpu(mem)

	// Default should be false
	if cpu.GetInterruptFlag() {
		t.Error("Expected interrupt flag to be false initially")
	}

	// Set Interrupt Flag
	cpu.SetInterruptFlag()
	if !cpu.GetInterruptFlag() {
		t.Error("Expected interrupt flag to be true after SetInterruptFlag")
	}
	if cpu.GetIRQFlag() {
		t.Error("Expected IRQ flag to be false after SetInterruptFlag")
	}

	// Clear Interrupt Flag
	cpu.ClearInterruptFlag()
	if cpu.GetInterruptFlag() {
		t.Error("Expected interrupt flag to be false after ClearInterruptFlag")
	}
	if !cpu.GetIRQFlag() {
		t.Error("Expected IRQ flag to be true after ClearInterruptFlag")
	}
}

func TestOverflowFlagInstructions(t *testing.T) {
	mem := memory.NewMemory()
	cpu := internal.NewCpu(mem)

	// Default should be false
	if cpu.GetOverflowFlag() {
		t.Error("Expected overflow flag to be false initially")
	}

	// Set Overflow Flag
	cpu.SetOverflowFlag()
	if !cpu.GetOverflowFlag() {
		t.Error("Expected overflow flag to be true after SetOverflowFlag")
	}

	// Clear Overflow Flag
	cpu.ClearOverflowFlag()
	if cpu.GetOverflowFlag() {
		t.Error("Expected overflow flag to be false after ClearOverflowFlag")
	}
}

func TestNoOpInstruction(t *testing.T) {
	mem := memory.NewMemory()
	cpu := internal.NewCpu(mem)

	// Save initial CPU state
	initialDebug := cpu.GetDebugData()
	initialCarry := cpu.GetCarryFlag()
	initialInterrupt := cpu.GetInterruptFlag()
	initialOverflow := cpu.GetOverflowFlag()
	initialIRQ := cpu.GetIRQFlag()

	// Call NoOp
	cpu.NoOp()

	// Verify state remains unchanged
	currentDebug := cpu.GetDebugData()
	if currentDebug["acc"] != initialDebug["acc"] || currentDebug["x"] != initialDebug["x"] || currentDebug["y"] != initialDebug["y"] {
		t.Error("Expected CPU registers to remain unchanged after NoOp")
	}
	if cpu.GetCarryFlag() != initialCarry || cpu.GetInterruptFlag() != initialInterrupt ||
		cpu.GetOverflowFlag() != initialOverflow || cpu.GetIRQFlag() != initialIRQ {
		t.Error("Expected CPU flags to remain unchanged after NoOp")
	}
}

func TestBreakInstruction(t *testing.T) {
	mem := memory.NewMemory()
	cpu := internal.NewCpu(mem)

	// Default irq should be false
	if cpu.GetIRQFlag() {
		t.Error("Expected IRQ flag to be false initially")
	}

	// Call Break
	cpu.Break()
	if !cpu.GetIRQFlag() {
		t.Error("Expected IRQ flag to be true after Break")
	}
}
