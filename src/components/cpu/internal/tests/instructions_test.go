package tests

import (
	"nes-emu/src/components/cpu/internal"
	"nes-emu/src/components/memory"
	"testing"
)

// Imediate addressing mode
func TestShouldLoadAcumulator(t *testing.T) {
	memory := memory.NewMemory()
	cpu := internal.NewCpu(memory)
	cpu.LoadValueIntoRegister(0x42, internal.ACCUMULATOR)

	if cpu.GetDebugData()["acc"] != 0x42 {
		t.Errorf("Expected acc to be 0x42, got %d", cpu.GetDebugData()["acc"])
	}
}

func TestShouldLoadXRegister(t *testing.T) {
	memory := memory.NewMemory()
	cpu := internal.NewCpu(memory)
	cpu.LoadValueIntoRegister(0x42, internal.REGISTER_X)

	if cpu.GetDebugData()["x"] != 0x42 {
		t.Errorf("Expected x to be 0x42, got %d", cpu.GetDebugData()["x"])
	}
}

func TestShouldLoadYRegister(t *testing.T) {
	memory := memory.NewMemory()
	cpu := internal.NewCpu(memory)
	cpu.LoadValueIntoRegister(0x42, internal.REGISTER_Y)

	if cpu.GetDebugData()["y"] != 0x42 {
		t.Errorf("Expected y to be 0x42, got %d", cpu.GetDebugData()["y"])
	}
}

// Absolute addressing mode
func TestShouldLoadValueFromMemoryIntoAccumulator(t *testing.T) {
	memory := memory.NewMemory()
	memory.Write(0x1234, 0x42)

	cpu := internal.NewCpu(memory)
	cpu.LoadValueFromMemoryIntoRegister(0x1234, internal.ACCUMULATOR)

	if cpu.GetDebugData()["acc"] != 0x42 {
		t.Errorf("Expected acc to be 0x42, got %d", cpu.GetDebugData()["acc"])
	}
}

func TestShouldLoadValueFromMemoryIntoXRegister(t *testing.T) {
	memory := memory.NewMemory()
	memory.Write(0x1234, 0x42)

	cpu := internal.NewCpu(memory)
	cpu.LoadValueFromMemoryIntoRegister(0x1234, internal.REGISTER_X)

	if cpu.GetDebugData()["x"] != 0x42 {
		t.Errorf("Expected x to be 0x42, got %d", cpu.GetDebugData()["x"])
	}
}

func TestShouldLoadValueFromMemoryIntoYRegister(t *testing.T) {
	memory := memory.NewMemory()
	memory.Write(0x1234, 0x42)

	cpu := internal.NewCpu(memory)
	cpu.LoadValueFromMemoryIntoRegister(0x1234, internal.REGISTER_Y)

	if cpu.GetDebugData()["y"] != 0x42 {
		t.Errorf("Expected y to be 0x42, got %d", cpu.GetDebugData()["y"])
	}
}

// ZeroPage addressing mode
func TestShouldLoadValueFromMemoryIntoAccumulatorZeroPage(t *testing.T) {
	memory := memory.NewMemory()
	memory.Write(0x0034, 0x42)

	cpu := internal.NewCpu(memory)
	cpu.LoadValueFromMemoryIntoRegister(0x0034, internal.ACCUMULATOR)

	if cpu.GetDebugData()["acc"] != 0x42 {
		t.Errorf("Expected acc to be 0x42, got %d", cpu.GetDebugData()["acc"])
	}
}

func TestShouldLoadValueFromMemoryIntoXRegisterZeroPage(t *testing.T) {
	memory := memory.NewMemory()
	memory.Write(0x0034, 0x42)

	cpu := internal.NewCpu(memory)
	cpu.LoadValueFromMemoryIntoRegister(0x0034, internal.REGISTER_X)

	if cpu.GetDebugData()["x"] != 0x42 {
		t.Errorf("Expected x to be 0x42, got %d", cpu.GetDebugData()["x"])
	}
}

func TestShouldLoadValueFromMemoryIntoYRegisterZeroPage(t *testing.T) {
	memory := memory.NewMemory()
	memory.Write(0x0034, 0x42)

	cpu := internal.NewCpu(memory)
	cpu.LoadValueFromMemoryIntoRegister(0x0034, internal.REGISTER_Y)

	if cpu.GetDebugData()["y"] != 0x42 {
		t.Errorf("Expected y to be 0x42, got %d", cpu.GetDebugData()["y"])
	}
}

// ZeroPage,X addressing mode
func TestShouldLoadValueFromMemoryIntoAccumulatorZeroPageXWhenXIsZero(t *testing.T) {
	memory := memory.NewMemory()
	memory.Write(0x0034, 0x42)

	cpu := internal.NewCpu(memory)
	cpu.LoadValueFromMemoryIntoAccumulatorWithIndexZeroPage(0x34, internal.REGISTER_X)

	if cpu.GetDebugData()["acc"] != 0x42 {
		t.Errorf("Expected acc to be 0x42, got %d", cpu.GetDebugData()["acc"])
	}
}

func TestShouldLoadValueFromMemoryIntoAccumulatorZeroPageYWhenYIsZero(t *testing.T) {
	memory := memory.NewMemory()
	memory.Write(0x0034, 0x42)

	cpu := internal.NewCpu(memory)
	cpu.LoadValueFromMemoryIntoAccumulatorWithIndexZeroPage(0x34, internal.REGISTER_Y)

	if cpu.GetDebugData()["acc"] != 0x42 {
		t.Errorf("Expected acc to be 0x42, got %d", cpu.GetDebugData()["acc"])
	}
}

func TestShouldLoadValueFromMemoryIntoAccumulatorZeroPageXWhenXIsNotZero(t *testing.T) {
	memory := memory.NewMemory()
	memory.Write(0x0034, 0x42)

	cpu := internal.NewCpu(memory)
	cpu.LoadValueIntoRegister(0x1, internal.REGISTER_X)

	cpu.LoadValueFromMemoryIntoAccumulatorWithIndexZeroPage(0x33, internal.REGISTER_X)

	if cpu.GetDebugData()["acc"] != 0x42 {
		t.Errorf("Expected acc to be 0x42, got %d", cpu.GetDebugData()["acc"])
	}
}

func TestShouldLoadValueFromMemoryIntoAccumulatorZeroPageYWhenYNotZero(t *testing.T) {
	memory := memory.NewMemory()
	memory.Write(0x0034, 0x42)

	cpu := internal.NewCpu(memory)
	cpu.LoadValueIntoRegister(0x1, internal.REGISTER_Y)

	cpu.LoadValueFromMemoryIntoAccumulatorWithIndexZeroPage(0x33, internal.REGISTER_Y)

	if cpu.GetDebugData()["acc"] != 0x42 {
		t.Errorf("Expected acc to be 0x42, got %d", cpu.GetDebugData()["acc"])
	}
}

func TestShouldLoadValueFromMemoryIntoAccumulatorZeroPageXWhenXOverflows(t *testing.T) {
	memory := memory.NewMemory()
	memory.Write(0x0000, 0x42)

	cpu := internal.NewCpu(memory)
	cpu.LoadValueIntoRegister(0xFF, internal.REGISTER_X)

	cpu.LoadValueFromMemoryIntoAccumulatorWithIndexZeroPage(0x01, internal.REGISTER_X)

	if cpu.GetDebugData()["acc"] != 0x42 {
		t.Errorf("Expected acc to be 0x42, got %d", cpu.GetDebugData()["acc"])
	}
}

func TestShouldLoadValueFromMemoryIntoAccumulatorZeroPageXWhenYOverflows(t *testing.T) {
	memory := memory.NewMemory()
	memory.Write(0x0000, 0x42)

	cpu := internal.NewCpu(memory)
	cpu.LoadValueIntoRegister(0xFF, internal.REGISTER_Y)

	cpu.LoadValueFromMemoryIntoAccumulatorWithIndexZeroPage(0x01, internal.REGISTER_Y)

	if cpu.GetDebugData()["acc"] != 0x42 {
		t.Errorf("Expected acc to be 0x42, got %d", cpu.GetDebugData()["acc"])
	}
}

// Absolute,X addressing mode
func TestShouldLoadValueFromMemoryIntoAccumulatorAbsoluteXWhenXIsZero(t *testing.T) {
	memory := memory.NewMemory()
	memory.Write(0x0034, 0x42)

	cpu := internal.NewCpu(memory)
	cpu.LoadValueFromMemoryIntoAccumulatorWithIndexAbsolute(0x34, internal.REGISTER_X)

	if cpu.GetDebugData()["acc"] != 0x42 {
		t.Errorf("Expected acc to be 0x42, got %d", cpu.GetDebugData()["acc"])
	}
}

func TestShouldLoadValueFromMemoryIntoAccumulatorAbsoluteYWhenYIsZero(t *testing.T) {
	memory := memory.NewMemory()
	memory.Write(0x0034, 0x42)

	cpu := internal.NewCpu(memory)
	cpu.LoadValueFromMemoryIntoAccumulatorWithIndexAbsolute(0x34, internal.REGISTER_Y)

	if cpu.GetDebugData()["acc"] != 0x42 {
		t.Errorf("Expected acc to be 0x42, got %d", cpu.GetDebugData()["acc"])
	}
}

func TestShouldLoadValueFromMemoryIntoAccumulatorAbsoluteXWhenXIsNotZero(t *testing.T) {
	memory := memory.NewMemory()
	memory.Write(0x0034, 0x42)

	cpu := internal.NewCpu(memory)
	cpu.LoadValueIntoRegister(0x1, internal.REGISTER_X)

	cpu.LoadValueFromMemoryIntoAccumulatorWithIndexAbsolute(0x33, internal.REGISTER_X)

	if cpu.GetDebugData()["acc"] != 0x42 {
		t.Errorf("Expected acc to be 0x42, got %d", cpu.GetDebugData()["acc"])
	}
}

func TestShouldLoadValueFromMemoryIntoAccumulatorAbsoluteYWhenYNotZero(t *testing.T) {
	memory := memory.NewMemory()
	memory.Write(0x0034, 0x42)

	cpu := internal.NewCpu(memory)
	cpu.LoadValueIntoRegister(0x1, internal.REGISTER_Y)

	cpu.LoadValueFromMemoryIntoAccumulatorWithIndexAbsolute(0x33, internal.REGISTER_Y)

	if cpu.GetDebugData()["acc"] != 0x42 {
		t.Errorf("Expected acc to be 0x42, got %d", cpu.GetDebugData()["acc"])
	}
}
