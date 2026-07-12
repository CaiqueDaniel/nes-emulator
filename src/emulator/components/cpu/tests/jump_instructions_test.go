package tests

import (
	"nes-emu/src/emulator/components/bus"
	internal "nes-emu/src/emulator/components/cpu"
	"nes-emu/src/emulator/components/memory"
	"testing"
)

func TestJump(t *testing.T) {
	mem := memory.NewMemory()
	b := bus.NewBus(mem)
	cpu := internal.NewCpuWithInternal(mem, b)

	cpu.JumpProgramCounterToValue(0x1234)

	if cpu.GetProgramCounter() != 0x1234 {
		t.Errorf("Expected program counter to be 0x1234, got %d", cpu.GetProgramCounter())
	}
}

func TestJumpProgramCounterByIndirectValue(t *testing.T) {
	mem := memory.NewMemory()
	b := bus.NewBus(mem)
	cpu := internal.NewCpuWithInternal(mem, b)

	mem.Write(0x03FF, 1)
	mem.Write(0x4000, 123)

	cpu.JumpProgramCounterByIndirectValue(0x03FF)

	if cpu.GetProgramCounter() != 0x0300 {
		t.Errorf("Expected instruction to corretcly imitate a hardware bug in the 6502 processor where the higher byte is read twice, to have the program counter to be 255, got %d", cpu.GetProgramCounter())
	}
}

func TestJumpProgramCounterToSubRoutine(t *testing.T) {
	mem := memory.NewMemory()
	b := bus.NewBus(mem)
	cpu := internal.NewCpuWithProgramCounter(mem, 0x1234, b)

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

func TestReturnFromSubRoutine(t *testing.T) {
	mem := memory.NewMemory()
	b := bus.NewBus(mem)
	cpu := internal.NewCpuWithProgramCounter(mem, 0x1234, b)

	cpu.JumpProgramCounterToSubRoutine(0x5678)
	cpu.ReturnFromSubRoutine()

	if cpu.GetProgramCounter() != 0x1235 {
		t.Errorf("Expected program counter to be 0x1235, got %d", cpu.GetProgramCounter())
	}
}

// ReturnFromInterrupt
// Stack layout pushed by the interrupt handler: high address, low address, flags (top of stack).
// ReturnFromInterrupt pops them in reverse: flags first, then low, then high.

func pushInterruptState(cpu interface {
	PushValueToStack(uint8)
}, flags, lowAddr, highAddr uint8) {
	// Push in reverse pull order: high first, low second, flags last (top of stack)
	cpu.PushValueToStack(highAddr)
	cpu.PushValueToStack(lowAddr)
	cpu.PushValueToStack(flags)
}

func TestReturnFromInterrupt_RestoresProgramCounter(t *testing.T) {
	mem := memory.NewMemory()
	b := bus.NewBus(mem)
	cpu := internal.NewCpuWithInternal(mem, b)

	pushInterruptState(cpu, 0x00, 0x34, 0x12)
	cpu.ReturnFromInterrupt()

	if cpu.GetProgramCounter() != 0x1234 {
		t.Errorf("Expected program counter to be 0x1234, got 0x%04x", cpu.GetProgramCounter())
	}
}

func TestReturnFromInterrupt_RestoresProgramCounter_HighByteOnly(t *testing.T) {
	mem := memory.NewMemory()
	b := bus.NewBus(mem)
	cpu := internal.NewCpuWithInternal(mem, b)

	pushInterruptState(cpu, 0x00, 0x00, 0xAB)
	cpu.ReturnFromInterrupt()

	if cpu.GetProgramCounter() != 0xAB00 {
		t.Errorf("Expected program counter to be 0xAB00, got 0x%04x", cpu.GetProgramCounter())
	}
}

func TestReturnFromInterrupt_RestoresCarryFlag(t *testing.T) {
	mem := memory.NewMemory()
	b := bus.NewBus(mem)
	cpu := internal.NewCpuWithInternal(mem, b)

	pushInterruptState(cpu, 0b00000001, 0x00, 0x00)
	cpu.ReturnFromInterrupt()

	if !cpu.GetCarryFlag() {
		t.Error("Expected carry flag to be restored from bit 0 of flags byte")
	}
}

func TestReturnFromInterrupt_RestoresZeroFlag(t *testing.T) {
	mem := memory.NewMemory()
	b := bus.NewBus(mem)
	cpu := internal.NewCpuWithInternal(mem, b)

	pushInterruptState(cpu, 0b00000010, 0x00, 0x00)
	cpu.ReturnFromInterrupt()

	if !cpu.GetZeroFlag() {
		t.Error("Expected zero flag to be restored from bit 1 of flags byte")
	}
}

func TestReturnFromInterrupt_RestoresIRQFlag(t *testing.T) {
	mem := memory.NewMemory()
	b := bus.NewBus(mem)
	cpu := internal.NewCpuWithInternal(mem, b)

	pushInterruptState(cpu, 0b00000100, 0x00, 0x00)
	cpu.ReturnFromInterrupt()

	if !cpu.GetIRQFlag() {
		t.Error("Expected IRQ flag to be restored from bit 2 of flags byte")
	}
}

func TestReturnFromInterrupt_RestoresDecimalFlag(t *testing.T) {
	mem := memory.NewMemory()
	b := bus.NewBus(mem)
	cpu := internal.NewCpuWithInternal(mem, b)

	pushInterruptState(cpu, 0b00001000, 0x00, 0x00)
	cpu.ReturnFromInterrupt()

	if !cpu.GetDecimalFlag() {
		t.Error("Expected decimal flag to be restored from bit 3 of flags byte")
	}
}

func TestReturnFromInterrupt_RestoresOverflowFlag(t *testing.T) {
	mem := memory.NewMemory()
	b := bus.NewBus(mem)
	cpu := internal.NewCpuWithInternal(mem, b)

	pushInterruptState(cpu, 0b01000000, 0x00, 0x00)
	cpu.ReturnFromInterrupt()

	if !cpu.GetOverflowFlag() {
		t.Error("Expected overflow flag to be restored from bit 6 of flags byte")
	}
}

func TestReturnFromInterrupt_RestoresNegativeFlag(t *testing.T) {
	mem := memory.NewMemory()
	b := bus.NewBus(mem)
	cpu := internal.NewCpuWithInternal(mem, b)

	pushInterruptState(cpu, 0b10000000, 0x00, 0x00)
	cpu.ReturnFromInterrupt()

	if !cpu.GetNegativeFlag() {
		t.Error("Expected negative flag to be restored from bit 7 of flags byte")
	}
}

func TestReturnFromInterrupt_Bits4And5InFlagsAreIgnored(t *testing.T) {
	mem := memory.NewMemory()
	b := bus.NewBus(mem)
	cpu := internal.NewCpuWithInternal(mem, b)

	pushInterruptState(cpu, 0b00110000, 0x00, 0x00)
	cpu.ReturnFromInterrupt()

	if cpu.GetCarryFlag() || cpu.GetZeroFlag() || cpu.GetIRQFlag() ||
		cpu.GetDecimalFlag() || cpu.GetOverflowFlag() || cpu.GetNegativeFlag() {
		t.Error("Expected all flags to remain false when only bits 4 and 5 are set in the flags byte")
	}
}

func TestReturnFromInterrupt_RestoresAllStateAtOnce(t *testing.T) {
	mem := memory.NewMemory()
	b := bus.NewBus(mem)
	cpu := internal.NewCpuWithInternal(mem, b)

	// carry(0) + IRQ(2) + overflow(6) + negative(7) set; address 0xBEEF
	flags := uint8(0b11000101)
	pushInterruptState(cpu, flags, 0xEF, 0xBE)
	cpu.ReturnFromInterrupt()

	if cpu.GetProgramCounter() != 0xBEEF {
		t.Errorf("Expected program counter to be 0xBEEF, got 0x%04x", cpu.GetProgramCounter())
	}
	if !cpu.GetCarryFlag() {
		t.Error("Expected carry flag to be set")
	}
	if !cpu.GetIRQFlag() {
		t.Error("Expected IRQ flag to be set")
	}
	if !cpu.GetOverflowFlag() {
		t.Error("Expected overflow flag to be set")
	}
	if !cpu.GetNegativeFlag() {
		t.Error("Expected negative flag to be set")
	}
	if cpu.GetZeroFlag() {
		t.Error("Expected zero flag to remain false")
	}
	if cpu.GetDecimalFlag() {
		t.Error("Expected decimal flag to remain false")
	}
}

func TestBreakInstruction(t *testing.T) {
	mem := memory.NewMemory()
	b := bus.NewBus(mem)
	mem.Write(0xFFFE, 0x12)
	cpu := internal.NewCpuWithProgramCounter(mem, 0x1234, b)

	// Default irq should be false
	if cpu.GetIRQFlag() {
		t.Error("Expected IRQ flag to be false initially")
	}

	b.ResetTickCount()
	cpu.Break()

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
