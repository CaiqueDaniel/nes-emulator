package tests

import (
	"nes-emu/src/emulator/components/bus"
	internal "nes-emu/src/emulator/components/cpu"
	"nes-emu/src/emulator/components/memory"
	"testing"
)

func TestCarryFlagInstructions(t *testing.T) {
	b := bus.NewBus()
	mem := memory.NewMemory(b)
	cpu := internal.NewCpu(mem, b)

	// Default should be false
	if cpu.GetCarryFlag() {
		t.Error("Expected carry flag to be false initially")
	}

	// Set Carry Flag
	cpu.SetCarryFlag()
	if !cpu.GetCarryFlag() {
		t.Error("Expected carry flag to be true after SetCarryFlag")
	}
	if cpu.GetTotalCycles() != 1 {
		t.Errorf("Expected total cycles to be 1, got %d", cpu.GetTotalCycles())
	}

	// Clear Carry Flag
	cpu.ClearCarryFlag()
	if cpu.GetCarryFlag() {
		t.Error("Expected carry flag to be false after ClearCarryFlag")
	}
	if cpu.GetTotalCycles() != 2 {
		t.Errorf("Expected total cycles to be 2, got %d", cpu.GetTotalCycles())
	}
}

func TestInterruptFlagInstructions(t *testing.T) {
	b := bus.NewBus()
	mem := memory.NewMemory(b)
	cpu := internal.NewCpu(mem, b)

	// Set Interrupt Flag
	cpu.SetInterruptFlag()
	if cpu.GetIRQFlag() {
		t.Error("Expected IRQ flag to be false after SetInterruptFlag")
	}
	if cpu.GetTotalCycles() != 1 {
		t.Errorf("Expected total cycles to be 1, got %d", cpu.GetTotalCycles())
	}

	// Clear Interrupt Flag
	cpu.ClearInterruptFlag()
	if !cpu.GetIRQFlag() {
		t.Error("Expected IRQ flag to be true after ClearInterruptFlag")
	}
	if cpu.GetTotalCycles() != 2 {
		t.Errorf("Expected total cycles to be 2, got %d", cpu.GetTotalCycles())
	}
}

func TestOverflowFlagInstructions(t *testing.T) {
	b := bus.NewBus()
	mem := memory.NewMemory(b)
	cpu := internal.NewCpu(mem, b)

	// Default should be false
	if cpu.GetOverflowFlag() {
		t.Error("Expected overflow flag to be false initially")
	}

	// Set Overflow Flag
	cpu.SetOverflowFlag()
	if !cpu.GetOverflowFlag() {
		t.Error("Expected overflow flag to be true after SetOverflowFlag")
	}
	if cpu.GetTotalCycles() != 1 {
		t.Errorf("Expected total cycles to be 1, got %d", cpu.GetTotalCycles())
	}

	// Clear Overflow Flag
	cpu.ClearOverflowFlag()
	if cpu.GetOverflowFlag() {
		t.Error("Expected overflow flag to be false after ClearOverflowFlag")
	}
	if cpu.GetTotalCycles() != 2 {
		t.Errorf("Expected total cycles to be 2, got %d", cpu.GetTotalCycles())
	}
}

func TestNoOpInstruction(t *testing.T) {
	b := bus.NewBus()
	mem := memory.NewMemory(b)
	cpu := internal.NewCpu(mem, b)

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

func TestBreakInstruction(t *testing.T) {
	b := bus.NewBusWithInternalType()
	mem := memory.NewMemory(b)
	mem.Write(0xFFFE, 0x12)
	cpu := internal.NewCpuWithProgramCounter(mem, 0x1234, b)

	// Default irq should be false
	if cpu.GetIRQFlag() {
		t.Error("Expected IRQ flag to be false initially")
	}

	b.ResetTickCount()
	cpu.Break()

	if b.GetTickCount() != 6 {
		t.Errorf("Expected total cycles to be 6 from instruction, got %d", b.GetTickCount())
	}

	if !cpu.GetIRQFlag() {
		t.Error("Expected IRQ flag to be true after Break")
	}

	// Verify PC is read from 0xFFFE
	if cpu.GetProgramCounter() != 0x0012 {
		t.Errorf("Expected Program Counter to be 0x0012, got 0x%04X", cpu.GetProgramCounter())
	}

	// Verify stack values (pulled in reverse order: Status, Low Address, High Address)
	status := cpu.PullValueFromStack()
	expectedStatusFlags := uint8(0b00110100) // Bit 5 (always 1), Bit 4 (B-flag), Bit 2 (IRQ)
	if status&expectedStatusFlags != expectedStatusFlags {
		t.Errorf("Expected Status with B-flag, unused and IRQ flag set, got %08b", status)
	}

	lowAddr := cpu.PullValueFromStack()
	if lowAddr != 0x34 {
		t.Errorf("Expected Low Address 0x34, got 0x%02X", lowAddr)
	}

	highAddr := cpu.PullValueFromStack()
	if highAddr != 0x12 {
		t.Errorf("Expected High Address 0x12, got 0x%02X", highAddr)
	}
}

func TestGetZeroFlag(t *testing.T) {
	b := bus.NewBus()
	mem := memory.NewMemory(b)
	cpu := internal.NewCpu(mem, b)

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
	b := bus.NewBus()
	mem := memory.NewMemory(b)
	cpu := internal.NewCpu(mem, b)

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
	b := bus.NewBus()
	mem := memory.NewMemory(b)
	cpu := internal.NewCpu(mem, b)

	// Default should be false
	if cpu.GetDecimalFlag() {
		t.Error("Expected decimal flag to be false initially")
	}

	// Set Decimal Flag
	cpu.SetDecimalFlag()
	if !cpu.GetDecimalFlag() {
		t.Error("Expected decimal flag to be true after SetDecimalFlag")
	}
	if cpu.GetTotalCycles() != 1 {
		t.Errorf("Expected total cycles to be 1, got %d", cpu.GetTotalCycles())
	}

	// Clear Decimal Flag
	cpu.ClearDecimalFlag()
	if cpu.GetDecimalFlag() {
		t.Error("Expected decimal flag to be false after ClearDecimalFlag")
	}
	if cpu.GetTotalCycles() != 2 {
		t.Errorf("Expected total cycles to be 2, got %d", cpu.GetTotalCycles())
	}
}
