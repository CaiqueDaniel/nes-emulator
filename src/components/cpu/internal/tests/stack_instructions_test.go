package tests

import (
	"nes-emu/src/components/cpu/internal"
	"nes-emu/src/components/memory"
	"testing"
)

func TestPushAccToStack(t *testing.T) {
	mem := memory.NewMemory()
	cpu := internal.NewCpu(mem)

	cpu.LoadValueIntoRegister(0x12, internal.ACCUMULATOR)
	cpu.PushAccToStack()

	if cpu.PullValueFromStack() != 0x12 {
		t.Errorf("Expected accumulator to be 0x12, got 0x%x", cpu.PullValueFromStack())
	}
}

func TestPullAccFromStack(t *testing.T) {
	mem := memory.NewMemory()
	cpu := internal.NewCpu(mem)

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
	cpu := internal.NewCpu(mem)

	cpu.LoadValueIntoRegister(0x12, internal.REGISTER_X)
	cpu.TransferXToStack()

	if cpu.GetStackPointer() != 0x12 {
		t.Errorf("Expected stack pointer to be 0x12, got 0x%x", cpu.GetStackPointer())
	}
}

func TestTransferStackToX(t *testing.T) {
	mem := memory.NewMemory()
	cpu := internal.NewCpu(mem)

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
