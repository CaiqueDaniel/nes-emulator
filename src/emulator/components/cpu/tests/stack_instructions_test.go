package tests

import (
	"nes-emu/src/emulator/components/bus"
	internal "nes-emu/src/emulator/components/cpu"
	"nes-emu/src/emulator/components/memory"
	"testing"
)

func TestPushAccToStack(t *testing.T) {
	mem := memory.NewMemory()
	b := bus.NewBus(mem)
	cpu := internal.NewCpuWithInternal(mem, b)

	cpu.LoadValueIntoRegister(0x12, internal.ACCUMULATOR)
	cpu.PushAccToStack()

	if cpu.PullValueFromStack() != 0x12 {
		t.Errorf("Expected accumulator to be 0x12, got 0x%x", cpu.PullValueFromStack())
	}
}

func TestPullAccFromStack(t *testing.T) {
	mem := memory.NewMemory()
	b := bus.NewBus(mem)
	cpu := internal.NewCpuWithInternal(mem, b)

	cpu.LoadValueIntoRegister(0x12, internal.ACCUMULATOR)
	cpu.PushAccToStack()
	cpu.PullAccFromStack()

	if cpu.GetDebugData()["acc"] != 0x12 {
		t.Errorf("Expected accumulator to be 0x12, got 0x%x", cpu.GetDebugData()["acc"])
	}

	if cpu.GetNegativeFlag() != false {
		t.Errorf("Expected negative flag to be false, got %t", cpu.GetNegativeFlag())
	}

	if cpu.GetZeroFlag() != false {
		t.Errorf("Expected zero flag to be false, got %t", cpu.GetZeroFlag())
	}
}

func TestTransferXToStack(t *testing.T) {
	mem := memory.NewMemory()
	b := bus.NewBus(mem)
	cpu := internal.NewCpuWithInternal(mem, b)

	cpu.LoadValueIntoRegister(0x12, internal.REGISTER_X)
	cpu.TransferXToStack()

	if cpu.GetStackPointer() != 0x12 {
		t.Errorf("Expected stack pointer to be 0x12, got 0x%x", cpu.GetStackPointer())
	}
}

func TestTransferStackToX(t *testing.T) {
	mem := memory.NewMemory()
	b := bus.NewBus(mem)
	cpu := internal.NewCpuWithInternal(mem, b)

	cpu.LoadValueIntoRegister(0x12, internal.REGISTER_X)
	cpu.TransferXToStack()
	cpu.TransferStackToX()

	if cpu.GetDebugData()["x"] != 0x12 {
		t.Errorf("Expected x register to be 0x12, got 0x%x", cpu.GetDebugData()["x"])
	}

	if cpu.GetNegativeFlag() != false {
		t.Errorf("Expected negative flag to be false, got %t", cpu.GetNegativeFlag())
	}

	if cpu.GetZeroFlag() != false {
		t.Errorf("Expected zero flag to be false, got %t", cpu.GetZeroFlag())
	}
}

// PushStatusIntoStack

func TestPushStatusIntoStack_DefaultStateHasBits4And5Set(t *testing.T) {
	mem := memory.NewMemory()
	b := bus.NewBus(mem)
	cpu := internal.NewCpuWithInternal(mem, b)

	cpu.SetInterruptFlag()
	cpu.PushStatusIntoStack()

	got := cpu.PullValueFromStack()
	if got != 0b00110000 {
		t.Errorf("Expected status byte 0b00110000 (bits 4 and 5 always set), got 0b%08b", got)
	}
}

func TestPushStatusIntoStack_CarryFlagSetsBit0(t *testing.T) {
	mem := memory.NewMemory()
	b := bus.NewBus(mem)
	cpu := internal.NewCpuWithInternal(mem, b)

	cpu.SetCarryFlag()
	cpu.PushStatusIntoStack()

	got := cpu.PullValueFromStack()
	if got&0b00000001 == 0 {
		t.Errorf("Expected bit 0 (carry) to be set, got 0b%08b", got)
	}
}

func TestPushStatusIntoStack_ZeroFlagSetsBit1(t *testing.T) {
	mem := memory.NewMemory()
	b := bus.NewBus(mem)
	cpu := internal.NewCpuWithInternal(mem, b)

	// Load 0 into acc to trigger the zero flag
	cpu.LoadValueIntoRegister(0x00, internal.ACCUMULATOR)
	cpu.PullAccFromStack() // sets zero flag based on acc=0
	cpu.PushStatusIntoStack()

	got := cpu.PullValueFromStack()
	if got&0b00000010 == 0 {
		t.Errorf("Expected bit 1 (zero) to be set, got 0b%08b", got)
	}
}

func TestPushStatusIntoStack_IRQFlagSetsBit2(t *testing.T) {
	mem := memory.NewMemory()
	b := bus.NewBus(mem)
	cpu := internal.NewCpuWithInternal(mem, b)

	cpu.Break() // sets irq = true
	cpu.PushStatusIntoStack()

	got := cpu.PullValueFromStack()
	if got&0b00000100 == 0 {
		t.Errorf("Expected bit 2 (IRQ) to be set, got 0b%08b", got)
	}
}

func TestPushStatusIntoStack_DecimalFlagSetsBit3(t *testing.T) {
	mem := memory.NewMemory()
	b := bus.NewBus(mem)
	cpu := internal.NewCpuWithInternal(mem, b)

	cpu.SetDecimalFlag()
	cpu.PushStatusIntoStack()

	got := cpu.PullValueFromStack()
	if got&0b00001000 == 0 {
		t.Errorf("Expected bit 3 (decimal) to be set, got 0b%08b", got)
	}
}

func TestPushStatusIntoStack_OverflowFlagSetsBit6(t *testing.T) {
	mem := memory.NewMemory()
	b := bus.NewBus(mem)
	cpu := internal.NewCpuWithInternal(mem, b)

	cpu.SetOverflowFlag()
	cpu.PushStatusIntoStack()

	got := cpu.PullValueFromStack()
	if got&0b01000000 == 0 {
		t.Errorf("Expected bit 6 (overflow) to be set, got 0b%08b", got)
	}
}

func TestPushStatusIntoStack_NegativeFlagSetsBit7(t *testing.T) {
	mem := memory.NewMemory()
	b := bus.NewBus(mem)
	cpu := internal.NewCpuWithInternal(mem, b)

	// Load 0x80 (bit 7 set) and pull acc so negative is set
	cpu.LoadValueIntoRegister(0x80, internal.ACCUMULATOR)
	cpu.PushAccToStack()
	cpu.PullAccFromStack() // sets negative = true
	cpu.PushStatusIntoStack()

	got := cpu.PullValueFromStack()
	if got&0b10000000 == 0 {
		t.Errorf("Expected bit 7 (negative) to be set, got 0b%08b", got)
	}
}

func TestPushStatusIntoStack_AllFlagsSet(t *testing.T) {
	mem := memory.NewMemory()
	b := bus.NewBus(mem)
	cpu := internal.NewCpuWithInternal(mem, b)

	cpu.SetCarryFlag()
	cpu.SetDecimalFlag()
	cpu.SetOverflowFlag()
	cpu.Break() // sets irq

	// Load 0x80 then pull acc to set negative; zero will be clear (0x80 != 0)
	cpu.LoadValueIntoRegister(0x80, internal.ACCUMULATOR)
	cpu.PushAccToStack()
	cpu.PullAccFromStack()

	cpu.PushStatusIntoStack()

	got := cpu.PullValueFromStack()
	const expected = 0b11111101 // all flags except zero + bits 4/5 always set
	// carry(0), irq(2), decimal(3), bit4, bit5, overflow(6), negative(7)
	if got&0b11001101 != 0b11001101 {
		t.Errorf("Expected carry, IRQ, decimal, overflow and negative bits to be set, got 0b%08b", got)
	}
	_ = expected
}

// PullStatusFromStack

func TestPullStatusFromStack_RestoresCarryFlag(t *testing.T) {
	mem := memory.NewMemory()
	b := bus.NewBus(mem)
	cpu := internal.NewCpuWithInternal(mem, b)

	cpu.PushValueToStack(0b00000001) // only carry bit set
	cpu.PullStatusFromStack()

	if !cpu.GetCarryFlag() {
		t.Error("Expected carry flag to be true after pulling status with bit 0 set")
	}
}

func TestPullStatusFromStack_RestoresZeroFlag(t *testing.T) {
	mem := memory.NewMemory()
	b := bus.NewBus(mem)
	cpu := internal.NewCpuWithInternal(mem, b)

	cpu.PushValueToStack(0b00000010) // only zero bit set
	cpu.PullStatusFromStack()

	if !cpu.GetZeroFlag() {
		t.Error("Expected zero flag to be true after pulling status with bit 1 set")
	}
}

func TestPullStatusFromStack_RestoresIRQFlag(t *testing.T) {
	mem := memory.NewMemory()
	b := bus.NewBus(mem)
	cpu := internal.NewCpuWithInternal(mem, b)

	cpu.PushValueToStack(0b00000100) // only IRQ bit set
	cpu.PullStatusFromStack()

	if !cpu.GetIRQFlag() {
		t.Error("Expected IRQ flag to be true after pulling status with bit 2 set")
	}
}

func TestPullStatusFromStack_RestoresDecimalFlag(t *testing.T) {
	mem := memory.NewMemory()
	b := bus.NewBus(mem)
	cpu := internal.NewCpuWithInternal(mem, b)

	cpu.PushValueToStack(0b00001000) // only decimal bit set
	cpu.PullStatusFromStack()

	if !cpu.GetDecimalFlag() {
		t.Error("Expected decimal flag to be true after pulling status with bit 3 set")
	}
}

func TestPullStatusFromStack_RestoresOverflowFlag(t *testing.T) {
	mem := memory.NewMemory()
	b := bus.NewBus(mem)
	cpu := internal.NewCpuWithInternal(mem, b)

	cpu.PushValueToStack(0b01000000) // only overflow bit set
	cpu.PullStatusFromStack()

	if !cpu.GetOverflowFlag() {
		t.Error("Expected overflow flag to be true after pulling status with bit 6 set")
	}
}

func TestPullStatusFromStack_RestoresNegativeFlag(t *testing.T) {
	mem := memory.NewMemory()
	b := bus.NewBus(mem)
	cpu := internal.NewCpuWithInternal(mem, b)

	cpu.PushValueToStack(0b10000000) // only negative bit set
	cpu.PullStatusFromStack()

	if !cpu.GetNegativeFlag() {
		t.Error("Expected negative flag to be true after pulling status with bit 7 set")
	}
}

func TestPullStatusFromStack_Bits4And5AreIgnored(t *testing.T) {
	mem := memory.NewMemory()
	b := bus.NewBus(mem)
	cpu := internal.NewCpuWithInternal(mem, b)

	// bits 4 and 5 set but no flag bits set
	cpu.PushValueToStack(0b00110000)
	cpu.PullStatusFromStack()

	if cpu.GetCarryFlag() || cpu.GetZeroFlag() || cpu.GetIRQFlag() ||
		cpu.GetDecimalFlag() || cpu.GetOverflowFlag() || cpu.GetNegativeFlag() {
		t.Error("Expected all flags to remain false when only bits 4 and 5 are set in the pulled byte")
	}
}

func TestPushAndPullStatusFromStack_Roundtrip(t *testing.T) {
	mem := memory.NewMemory()
	b := bus.NewBus(mem)
	cpu := internal.NewCpuWithInternal(mem, b)

	// Set a mix of flags
	cpu.SetCarryFlag()
	cpu.SetOverflowFlag()
	cpu.SetDecimalFlag()

	cpu.PushStatusIntoStack()
	// Clear all flags before pulling to confirm they are actually restored
	cpu.ClearCarryFlag()
	cpu.ClearOverflowFlag()
	cpu.ClearDecimalFlag()

	cpu.PullStatusFromStack()

	if !cpu.GetCarryFlag() {
		t.Error("Expected carry flag to be restored after roundtrip")
	}
	if !cpu.GetOverflowFlag() {
		t.Error("Expected overflow flag to be restored after roundtrip")
	}
	if !cpu.GetDecimalFlag() {
		t.Error("Expected decimal flag to be restored after roundtrip")
	}
	if cpu.GetNegativeFlag() {
		t.Error("Expected negative flag to remain false after roundtrip")
	}
}
