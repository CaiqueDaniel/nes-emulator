package tests

import (
	"nes-emu/src/emulator/components/bus"
	internal "nes-emu/src/emulator/components/cpu"
	"nes-emu/src/emulator/components/memory"
	"testing"
)

func TestCarryFlagInstructions(t *testing.T) {
	mem := memory.NewMemory()
	b := bus.NewBus(mem)
	cpu := internal.NewCpuWithInternal(b)

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
	b := bus.NewBus(mem)
	cpu := internal.NewCpuWithInternal(b)

	// Set Interrupt Flag
	cpu.SetInterruptFlag()
	if cpu.GetIRQFlag() {
		t.Error("Expected IRQ flag to be false after SetInterruptFlag")
	}

	// Clear Interrupt Flag
	cpu.ClearInterruptFlag()
	if !cpu.GetIRQFlag() {
		t.Error("Expected IRQ flag to be true after ClearInterruptFlag")
	}
}

func TestOverflowFlagInstructions(t *testing.T) {
	mem := memory.NewMemory()
	b := bus.NewBus(mem)
	cpu := internal.NewCpuWithInternal(b)

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
	b := bus.NewBus(mem)
	cpu := internal.NewCpuWithInternal(b)

	// Save initial CPU state
	initialDebug := cpu.GetDebugData()
	initialCarry := cpu.GetCarryFlag()
	initialOverflow := cpu.GetOverflowFlag()
	initialIRQ := cpu.GetIRQFlag()

	// Call NoOp
	cpu.NoOp()

	// Verify state remains unchanged
	currentDebug := cpu.GetDebugData()
	if currentDebug["acc"] != initialDebug["acc"] || currentDebug["x"] != initialDebug["x"] || currentDebug["y"] != initialDebug["y"] {
		t.Error("Expected CPU registers to remain unchanged after NoOp")
	}

	if cpu.GetCarryFlag() != initialCarry ||
		cpu.GetOverflowFlag() != initialOverflow || cpu.GetIRQFlag() != initialIRQ {
		t.Error("Expected CPU flags to remain unchanged after NoOp")
	}
}

func TestGetZeroFlag(t *testing.T) {
	mem := memory.NewMemory()
	b := bus.NewBus(mem)
	cpu := internal.NewCpuWithInternal(b)

	// Default should be false
	if cpu.GetZeroFlag() {
		t.Error("Expected zero flag to be false initially")
	}

	// Load zero into accumulator and perform AND to trigger zero flag
	cpu.LoadValueIntoRegister(0x00, internal.ACCUMULATOR)
	cpu.And(0x00)
	if !cpu.GetZeroFlag() {
		t.Error("Expected zero flag to be true after AND resulting in zero")
	}

	// Load non-zero value and trigger clear of zero flag
	cpu.LoadValueIntoRegister(0x01, internal.ACCUMULATOR)
	cpu.And(0xFF)
	if cpu.GetZeroFlag() {
		t.Error("Expected zero flag to be false after AND resulting in non-zero")
	}
}

func TestGetNegativeFlag(t *testing.T) {
	mem := memory.NewMemory()
	b := bus.NewBus(mem)
	cpu := internal.NewCpuWithInternal(b)

	// Default should be false
	if cpu.GetNegativeFlag() {
		t.Error("Expected negative flag to be false initially")
	}

	// Load a value with bit 7 set to trigger negative flag
	cpu.LoadValueIntoRegister(0xFF, internal.ACCUMULATOR)
	cpu.And(0b10000000)
	if !cpu.GetNegativeFlag() {
		t.Error("Expected negative flag to be true when bit 7 of result is set")
	}

	// Load a value with bit 7 clear to clear negative flag
	cpu.LoadValueIntoRegister(0xFF, internal.ACCUMULATOR)
	cpu.And(0b01111111)
	if cpu.GetNegativeFlag() {
		t.Error("Expected negative flag to be false when bit 7 of result is clear")
	}
}

func TestDecimalFlagInstructions(t *testing.T) {
	mem := memory.NewMemory()
	b := bus.NewBus(mem)
	cpu := internal.NewCpuWithInternal(b)

	// Default should be false
	if cpu.GetDecimalFlag() {
		t.Error("Expected decimal flag to be false initially")
	}

	// Set Decimal Flag
	cpu.SetDecimalFlag()
	if !cpu.GetDecimalFlag() {
		t.Error("Expected decimal flag to be true after SetDecimalFlag")
	}

	// Clear Decimal Flag
	cpu.ClearDecimalFlag()
	if cpu.GetDecimalFlag() {
		t.Error("Expected decimal flag to be false after ClearDecimalFlag")
	}
}
