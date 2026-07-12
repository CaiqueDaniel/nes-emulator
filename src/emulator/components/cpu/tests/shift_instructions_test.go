package tests

import (
	"nes-emu/src/emulator/components/bus"
	internal "nes-emu/src/emulator/components/cpu"
	"nes-emu/src/emulator/components/memory"
	"testing"
)

func TestArithmeticShiftLeftOnAccWithPositiveValue(t *testing.T) {
	mem := memory.NewMemory()
	b := bus.NewBus(mem)
	cpu := internal.NewCpuWithInternal(mem, b)

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
	b := bus.NewBus(mem)
	cpu := internal.NewCpuWithInternal(mem, b)

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
	b := bus.NewBus(mem)
	cpu := internal.NewCpuWithInternal(mem, b)

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
	b := bus.NewBus(mem)
	cpu := internal.NewCpuWithInternal(mem, b)

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
	b := bus.NewBus(mem)
	cpu := internal.NewCpuWithInternal(mem, b)

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

func TestLogicalShiftRightOnAccWithPositiveValue(t *testing.T) {
	mem := memory.NewMemory()
	b := bus.NewBus(mem)
	cpu := internal.NewCpuWithInternal(mem, b)

	cpu.LoadValueIntoRegister(0, internal.ACCUMULATOR)
	cpu.AddWithCarry(1)
	cpu.LogicalShiftRight()

	if cpu.GetDebugData()["acc"] != 0 {
		t.Errorf("Expected accumulator to be 0, got %d", cpu.GetDebugData()["acc"])
	}

	if !cpu.GetZeroFlag() {
		t.Error("Expected zero flag to be true")
	}

	if cpu.GetNegativeFlag() {
		t.Error("Expected negative flag to be false")
	}

	if cpu.GetOverflowFlag() {
		t.Error("Expected overflow flag to be false")
	}

	if !cpu.GetCarryFlag() {
		t.Error("Expected carry flag to be true")
	}
}

func TestLogicalShiftRightOnAccWithNegativeValue(t *testing.T) {
	mem := memory.NewMemory()
	b := bus.NewBus(mem)
	cpu := internal.NewCpuWithInternal(mem, b)

	cpu.LoadValueIntoRegister(254, internal.ACCUMULATOR)
	cpu.LogicalShiftRight()

	if cpu.GetDebugData()["acc"] != 127 {
		t.Errorf("Expected accumulator to be 127, got %d", cpu.GetDebugData()["acc"])
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

func TestLogicalShiftRightOnAbsoluteAddressWithPositiveValue(t *testing.T) {
	mem := memory.NewMemory()
	b := bus.NewBus(mem)
	cpu := internal.NewCpuWithInternal(mem, b)

	mem.Write(0, 1)
	cpu.LogicalShiftRightAbsolute(0x0000, false)

	if mem.Read(0) != 0 {
		t.Errorf("Expected accumulator to be 0, got %d", mem.Read(0))
	}

	if !cpu.GetZeroFlag() {
		t.Error("Expected zero flag to be true")
	}

	if cpu.GetNegativeFlag() {
		t.Error("Expected negative flag to be false")
	}

	if cpu.GetOverflowFlag() {
		t.Error("Expected overflow flag to be false")
	}

	if !cpu.GetCarryFlag() {
		t.Error("Expected carry flag to be true")
	}
}

// ─── RotateLeft ──────────────────────────────────────────────────────────────

func TestRotateLeftOnAccWithNoCarry(t *testing.T) {
	mem := memory.NewMemory()
	b := bus.NewBus(mem)
	cpu := internal.NewCpuWithInternal(mem, b)

	// acc = 0b00000011 (3), carry = false → result = 0b00000110 (6)
	cpu.LoadValueIntoRegister(0b00000011, internal.ACCUMULATOR)
	cpu.RotateLeft()

	if cpu.GetDebugData()["acc"] != 0b00000110 {
		t.Errorf("Expected accumulator to be 6, got %d", cpu.GetDebugData()["acc"])
	}

	// MSB of original value was 0 → carry = false
	if cpu.GetCarryFlag() {
		t.Error("Expected carry flag to be false")
	}

	if cpu.GetNegativeFlag() {
		t.Error("Expected negative flag to be false")
	}
}

func TestRotateLeftOnAccWithPriorCarrySet(t *testing.T) {
	mem := memory.NewMemory()
	b := bus.NewBus(mem)
	cpu := internal.NewCpuWithInternal(mem, b)

	// acc = 0b00000001 (1), carry = true → result = 0b00000011 (3)
	// old carry (1) is shifted into bit 0
	cpu.LoadValueIntoRegister(0b00000001, internal.ACCUMULATOR)
	cpu.SetCarryFlag()
	cpu.RotateLeft()

	if cpu.GetDebugData()["acc"] != 0b00000011 {
		t.Errorf("Expected accumulator to be 3, got %d", cpu.GetDebugData()["acc"])
	}

	// MSB of original value was 0 → carry = false
	if cpu.GetCarryFlag() {
		t.Error("Expected carry flag to be false")
	}

	if cpu.GetNegativeFlag() {
		t.Error("Expected negative flag to be false")
	}
}

func TestRotateLeftOnAccWithMSBSet(t *testing.T) {
	mem := memory.NewMemory()
	b := bus.NewBus(mem)
	cpu := internal.NewCpuWithInternal(mem, b)

	cpu.LoadValueIntoRegister(0b10000001, internal.ACCUMULATOR)
	cpu.RotateLeft()

	if cpu.GetDebugData()["acc"] != 0b00000010 {
		t.Errorf("Expected accumulator to be 2, got %d", cpu.GetDebugData()["acc"])
	}

	if !cpu.GetCarryFlag() {
		t.Error("Expected carry flag to be true")
	}

	if cpu.GetNegativeFlag() {
		t.Error("Expected negative flag to be false")
	}
}

// ─── RotateLeftAbsolute ───────────────────────────────────────────────────────

func TestRotateLeftAbsoluteWithNoCarry(t *testing.T) {
	mem := memory.NewMemory()
	b := bus.NewBus(mem)
	cpu := internal.NewCpuWithInternal(mem, b)

	mem.Write(0x0010, 0b00000010)
	cpu.ClearCarryFlag()
	cpu.RotateLeftAbsolute(0x0010, false)

	if mem.Read(0x0010) != 0b00000100 {
		t.Errorf("Expected memory[0x10] to be 4, got %d", mem.Read(0x0010))
	}

	if cpu.GetCarryFlag() {
		t.Error("Expected carry flag to be false")
	}

	if cpu.GetNegativeFlag() {
		t.Error("Expected negative flag to be false")
	}

	if cpu.GetZeroFlag() {
		t.Error("Expected zero flag to be false")
	}
}

func TestRotateLeftAbsoluteWithPriorCarry(t *testing.T) {
	mem := memory.NewMemory()
	b := bus.NewBus(mem)
	cpu := internal.NewCpuWithInternal(mem, b)

	// value = 0b00000010 (2), old carry = 1 → result = 0b00000101 (5)
	mem.Write(0x0010, 0b00000010)
	cpu.SetCarryFlag()
	cpu.RotateLeftAbsolute(0x0010, false)

	if mem.Read(0x0010) != 0b00000101 {
		t.Errorf("Expected memory[0x10] to be 5, got %d", mem.Read(0x0010))
	}

	if cpu.GetCarryFlag() {
		t.Error("Expected carry flag to be true")
	}
}

func TestRotateLeftAbsoluteWithXIndexedMode(t *testing.T) {
	mem := memory.NewMemory()
	b := bus.NewBus(mem)
	cpu := internal.NewCpuWithInternal(mem, b)

	// X = 2, base address = 0x0010 → effective address = 0x0012
	cpu.LoadValueIntoRegister(2, internal.REGISTER_X)
	mem.Write(0x0012, 0b00000001)
	cpu.RotateLeftAbsolute(0x0010, true)

	if mem.Read(0x0012) != 0b00000010 {
		t.Errorf("Expected memory[0x12] to be 2, got %d", mem.Read(0x0012))
	}

	if cpu.GetCarryFlag() {
		t.Error("Expected carry flag to be false")
	}
}

// ─── RotateLeftZeroPage ───────────────────────────────────────────────────────

func TestRotateLeftZeroPageWithNoCarry(t *testing.T) {
	mem := memory.NewMemory()
	b := bus.NewBus(mem)
	cpu := internal.NewCpuWithInternal(mem, b)

	mem.Write(0x0005, 0b00000100)
	cpu.ClearCarryFlag()
	cpu.RotateLeftZeroPage(0x05, false)

	if mem.Read(0x0005) != 0b00001000 {
		t.Errorf("Expected memory[0x05] to be 8, got %d", mem.Read(0x0005))
	}

	if cpu.GetCarryFlag() {
		t.Error("Expected carry flag to be false")
	}

	if cpu.GetNegativeFlag() {
		t.Error("Expected negative flag to be false")
	}
}

// ─── RotateRight ─────────────────────────────────────────────────────────────

func TestRotateRightOnAccWithNoCarry(t *testing.T) {
	mem := memory.NewMemory()
	b := bus.NewBus(mem)
	cpu := internal.NewCpuWithInternal(mem, b)

	// acc = 0b00000110 (6), carry = false → result = 0b00000011 (3)
	cpu.LoadValueIntoRegister(0b00000110, internal.ACCUMULATOR)
	cpu.ClearCarryFlag()
	cpu.RotateRight()

	if cpu.GetDebugData()["acc"] != 0b00000011 {
		t.Errorf("Expected accumulator to be 3, got %d", cpu.GetDebugData()["acc"])
	}

	// LSB of original value was 0 → carry = false
	if cpu.GetCarryFlag() {
		t.Error("Expected carry flag to be false")
	}

	if cpu.GetNegativeFlag() {
		t.Error("Expected negative flag to be false")
	}
}

func TestRotateRightOnAccWithPriorCarrySet(t *testing.T) {
	mem := memory.NewMemory()
	b := bus.NewBus(mem)
	cpu := internal.NewCpuWithInternal(mem, b)

	cpu.LoadValueIntoRegister(0b00000100, internal.ACCUMULATOR)
	cpu.SetCarryFlag()
	cpu.RotateRight()

	if cpu.GetDebugData()["acc"] != 0b10000010 {
		t.Errorf("Expected accumulator to rotate right with carry true, got %d", cpu.GetDebugData()["acc"])
	}

	if cpu.GetCarryFlag() {
		t.Error("Expected carry flag to be true")
	}
	if !cpu.GetNegativeFlag() {
		t.Error("Expected negative flag to be true")
	}
}

func TestRotateRightOnAccWithLSBSet(t *testing.T) {
	mem := memory.NewMemory()
	b := bus.NewBus(mem)
	cpu := internal.NewCpuWithInternal(mem, b)

	// acc = 0b00000011 (3), carry = false → result = 0b00000001 (1)
	// LSB of original value was 1 → carry = true
	cpu.LoadValueIntoRegister(0b00000011, internal.ACCUMULATOR)
	cpu.ClearCarryFlag()
	cpu.RotateRight()

	if cpu.GetDebugData()["acc"] != 0b00000001 {
		t.Errorf("Expected accumulator to be 1, got %d", cpu.GetDebugData()["acc"])
	}

	if !cpu.GetCarryFlag() {
		t.Error("Expected carry flag to be true because LSB of original value was set")
	}

	if cpu.GetNegativeFlag() {
		t.Error("Expected negative flag to be false")
	}
}

// ─── RotateRightAbsolute ──────────────────────────────────────────────────────

func TestRotateRightAbsoluteWithNoCarry(t *testing.T) {
	mem := memory.NewMemory()
	b := bus.NewBus(mem)
	cpu := internal.NewCpuWithInternal(mem, b)

	mem.Write(0x0010, 0b00001000)
	cpu.ClearCarryFlag()
	cpu.RotateRightAbsolute(0x0010, false)

	if mem.Read(0x0010) != 0b00000100 {
		t.Errorf("Expected memory[0x10] to be 4, got %d", mem.Read(0x0010))
	}

	if cpu.GetCarryFlag() {
		t.Error("Expected carry flag to be false")
	}

	if cpu.GetNegativeFlag() {
		t.Error("Expected negative flag to be false")
	}

	if cpu.GetZeroFlag() {
		t.Error("Expected zero flag to be false")
	}
}

func TestRotateRightAbsoluteWithPriorCarry(t *testing.T) {
	mem := memory.NewMemory()
	b := bus.NewBus(mem)
	cpu := internal.NewCpuWithInternal(mem, b)

	mem.Write(0x0010, 0b00000100)
	cpu.SetCarryFlag()
	cpu.RotateRightAbsolute(0x0010, false)

	if mem.Read(0x0010) != 0b10000010 {
		t.Errorf("Expected memory[0x10] to be 3, got %d", mem.Read(0x0010))
	}

	if cpu.GetCarryFlag() {
		t.Error("Expected carry flag to be false")
	}

	if !cpu.GetNegativeFlag() {
		t.Error("Expected negative flag to be true")
	}
}

func TestRotateRightAbsoluteWithXIndexedMode(t *testing.T) {
	mem := memory.NewMemory()
	b := bus.NewBus(mem)
	cpu := internal.NewCpuWithInternal(mem, b)

	// X = 2, base address = 0x0010 → effective address = 0x0012
	cpu.LoadValueIntoRegister(2, internal.REGISTER_X)
	mem.Write(0x0012, 0b00000010)
	cpu.ClearCarryFlag()
	cpu.RotateRightAbsolute(0x0010, true)

	if mem.Read(0x0012) != 0b00000001 {
		t.Errorf("Expected memory[0x12] to be 1, got %d", mem.Read(0x0012))
	}

	if cpu.GetCarryFlag() {
		t.Error("Expected carry flag to be false")
	}
}

func TestRotateRightAbsoluteWithLSBSet(t *testing.T) {
	mem := memory.NewMemory()
	b := bus.NewBus(mem)
	cpu := internal.NewCpuWithInternal(mem, b)

	// value = 0b00000011 (3), carry = false → result = 0b00000001 (1), carry = true
	mem.Write(0x0020, 0b00000011)
	cpu.ClearCarryFlag()
	cpu.RotateRightAbsolute(0x0020, false)

	if mem.Read(0x0020) != 0b00000001 {
		t.Errorf("Expected memory[0x20] to be 1, got %d", mem.Read(0x0020))
	}

	if !cpu.GetCarryFlag() {
		t.Error("Expected carry flag to be true because LSB of original value was set")
	}
}

// ─── RotateRightZeroPage ──────────────────────────────────────────────────────

func TestRotateRightZeroPageWithNoCarry(t *testing.T) {
	mem := memory.NewMemory()
	b := bus.NewBus(mem)
	cpu := internal.NewCpuWithInternal(mem, b)

	mem.Write(0x0005, 0b00010000)
	cpu.ClearCarryFlag()
	cpu.RotateRightZeroPage(0x05, false)

	if mem.Read(0x0005) != 0b00001000 {
		t.Errorf("Expected memory[0x05] to be 8, got %d", mem.Read(0x0005))
	}

	if cpu.GetCarryFlag() {
		t.Error("Expected carry flag to be false")
	}

	if cpu.GetNegativeFlag() {
		t.Error("Expected negative flag to be false")
	}
}
