package tests

import (
	"testing"

	"nes-emu/src/components/cpu/internal"
	"nes-emu/src/components/memory"
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
