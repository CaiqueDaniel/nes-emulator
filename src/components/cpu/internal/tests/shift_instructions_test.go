package tests

import (
	"nes-emu/src/components/cpu/internal"
	"nes-emu/src/components/memory"
	"testing"
)

func TestArithmeticShiftLeftOnAccWithPositiveValue(t *testing.T) {
	mem := memory.NewMemory()
	cpu := internal.NewCpu(mem)

	cpu.AddWithCarry(1)
	cpu.ArithmeticShiftLeft()

	if cpu.GetDebugData()["acc"] != 2 {
		t.Errorf("Expected accumulator to be 2, got %d", cpu.GetDebugData()["acc"])
	}

	if cpu.GetZeroFlag() {
		t.Error("Expected zero flag to be false")
	}

	if cpu.GetNegativeFlag() {
		t.Error("Expected negative flag to be false")
	}

	if cpu.GetOverflowFlag() {
		t.Error("Expected overflow flag to be false")
	}

	if cpu.GetCarryFlag() {
		t.Error("Expected carry flag to be false")
	}
}

func TestArithmeticShiftLeftOnAccWithNegativeValue(t *testing.T) {
	mem := memory.NewMemory()
	cpu := internal.NewCpu(mem)

	cpu.LoadValueIntoRegister(0, internal.ACCUMULATOR)
	cpu.AddWithCarry(127)
	cpu.ArithmeticShiftLeft()

	if cpu.GetDebugData()["acc"] != 254 {
		t.Errorf("Expected accumulator to be 254, got %d", cpu.GetDebugData()["acc"])
	}

	if cpu.GetZeroFlag() {
		t.Error("Expected zero flag to be false")
	}

	if !cpu.GetNegativeFlag() {
		t.Error("Expected negative flag to be true")
	}

	if cpu.GetOverflowFlag() {
		t.Error("Expected overflow flag to be false")
	}

	if cpu.GetCarryFlag() {
		t.Error("Expected carry flag to be false")
	}
}

func TestArithmeticShiftLeftOnAccWithOverflowValue(t *testing.T) {
	mem := memory.NewMemory()
	cpu := internal.NewCpu(mem)

	cpu.LoadValueIntoRegister(0, internal.ACCUMULATOR)
	cpu.AddWithCarry(129)
	cpu.ArithmeticShiftLeft()

	if cpu.GetDebugData()["acc"] != 2 {
		t.Errorf("Expected accumulator to be 2, got %d", cpu.GetDebugData()["acc"])
	}

	if cpu.GetZeroFlag() {
		t.Error("Expected zero flag to be false")
	}

	if cpu.GetNegativeFlag() {
		t.Error("Expected negative flag to be false")
	}

	if cpu.GetOverflowFlag() {
		t.Error("Expected overflow flag to be true")
	}

	if !cpu.GetCarryFlag() {
		t.Error("Expected carry flag to be false")
	}
}

func TestArithmeticShiftLeftOnAbsoluteAddressWithPositiveValue(t *testing.T) {
	mem := memory.NewMemory()
	cpu := internal.NewCpu(mem)

	mem.Write(0, 1)
	cpu.ArithmeticShiftLeftAbsolute(0x0000, false)

	if mem.Read(0) != 2 {
		t.Errorf("Expected accumulator to be 2, got %d", mem.Read(0))
	}

	if cpu.GetZeroFlag() {
		t.Error("Expected zero flag to be false")
	}

	if cpu.GetNegativeFlag() {
		t.Error("Expected negative flag to be false")
	}

	if cpu.GetOverflowFlag() {
		t.Error("Expected overflow flag to be false")
	}

	if cpu.GetCarryFlag() {
		t.Error("Expected carry flag to be false")
	}
}

func TestArithmeticShiftLeftOnAbsoluteAddressWithPositiveValueAndXIndexedMode(t *testing.T) {
	mem := memory.NewMemory()
	cpu := internal.NewCpu(mem)

	cpu.LoadValueIntoRegister(1, internal.REGISTER_X)
	mem.Write(1, 1)
	cpu.ArithmeticShiftLeftAbsolute(0x0000, true)

	if mem.Read(1) != 2 {
		t.Errorf("Expected accumulator to be 2, got %d", mem.Read(1))
	}

	if cpu.GetZeroFlag() {
		t.Error("Expected zero flag to be false")
	}

	if cpu.GetNegativeFlag() {
		t.Error("Expected negative flag to be false")
	}

	if cpu.GetOverflowFlag() {
		t.Error("Expected overflow flag to be false")
	}

	if cpu.GetCarryFlag() {
		t.Error("Expected carry flag to be false")
	}
}
