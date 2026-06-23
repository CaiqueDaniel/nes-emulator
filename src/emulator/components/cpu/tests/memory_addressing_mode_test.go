package tests

import (
	"testing"

	internal "nes-emu/src/emulator/components/cpu"
	"nes-emu/src/emulator/components/memory"
)

func TestGetValueByAbsoluteMode(t *testing.T) {
	memory := memory.NewMemory()
	cpu := internal.NewCpu(memory)

	memory.Write(0x1000, 0xFF)

	value := cpu.GetValueByAbsoluteMode(0x1000)

	if value != 0xFF {
		t.Errorf("Expected %d, got %d", 0xFF, value)
	}
}

func TestGetValueByZeroPageMode(t *testing.T) {
	memory := memory.NewMemory()
	cpu := internal.NewCpu(memory)

	memory.Write(0x0010, 0xFF)

	value := cpu.GetValueByZeroPageMode(0x10)

	if value != 0xFF {
		t.Errorf("Expected %d, got %d", 0xFF, value)
	}
}

func TestGetValueByIndexedAbsoluteMode(t *testing.T) {
	memory := memory.NewMemory()
	cpu := internal.NewCpu(memory)

	memory.Write(0x1010, 0xFF)

	value := cpu.GetValueByIndexedAbsoluteMode(0x1000, 0x10)

	if value != 0xFF {
		t.Errorf("Expected %d, got %d", 0xFF, value)
	}
}

func TestGetValueByZeroPageIndexedMode(t *testing.T) {
	memory := memory.NewMemory()
	cpu := internal.NewCpu(memory)

	memory.Write(0x000F, 0xFF)

	value := cpu.GetValueByZeroPageIndexedMode(0x000E, 0x1)

	if value != 0xFF {
		t.Errorf("Expected %d, got %d", 0xFF, value)
	}
}

func TestGetValueByIndirectAbsoluteMode(t *testing.T) {
	memory := memory.NewMemory()
	cpu := internal.NewCpu(memory)

	memory.Write(0xFAFF, 0x80)
	memory.Write(0x1000, 0xFF)
	memory.Write(0x1001, 0xFA)

	value := cpu.GetValueByIndirectAbsoluteMode(0x1000)

	if value != 0x80 {
		t.Errorf("Expected %d, got %d", 0x80, value)
	}
}

func TestGetValueByIndexedIndirectXMode(t *testing.T) {
	memory := memory.NewMemory()
	cpu := internal.NewCpu(memory)

	memory.Write(0xFAFF, 0x80)
	memory.Write(0x01, 0xFF)
	memory.Write(0x02, 0xFA)
	cpu.LoadValueIntoRegister(0x1, internal.REGISTER_X)

	value := cpu.GetValueByIndexedIndirectXMode(0x0)

	if value != 0x80 {
		t.Errorf("Expected %d, got %d", 0x80, value)
	}
}

func TestGetValueByIndirectIndexedYMode(t *testing.T) {
	memory := memory.NewMemory()
	cpu := internal.NewCpu(memory)

	memory.Write(0xFAFF, 0x80)
	memory.Write(0x01, 0xFF)
	memory.Write(0x02, 0xFA)
	cpu.LoadValueIntoRegister(0x1, internal.REGISTER_X)

	value := cpu.GetValueByIndexedIndirectXMode(0x0)

	if value != 0x80 {
		t.Errorf("Expected %d, got %d", 0x80, value)
	}
}
