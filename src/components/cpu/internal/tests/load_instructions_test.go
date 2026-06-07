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

// Indexed Indirect (Indirect,X) addressing mode
func TestShouldLoadRegisterFromMemoryWithIndexedXIndirectAddress(t *testing.T) {
	tests := []struct {
		name     string
		register string
		value    uint8
	}{
		{
			name:     "Load ACC with Indexed X Indirect Addressing",
			register: internal.ACCUMULATOR,
			value:    0x42,
		},
		{
			name:     "Load X with Indexed X Indirect Addressing",
			register: internal.REGISTER_X,
			value:    0x33,
		},
		{
			name:     "Load Y with Indexed X Indirect Addressing",
			register: internal.REGISTER_Y,
			value:    0x55,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mem := memory.NewMemory()
			cpu := internal.NewCpu(mem)

			// Setup X register
			var xVal uint8 = 0x04
			cpu.LoadValueIntoRegister(xVal, internal.REGISTER_X)

			// Pivot address is 0x20 + 0x04 = 0x24
			var targetAddr uint16 = 0x5678
			mem.Write(0x24, uint8(targetAddr&0xFF))
			mem.Write(0x25, uint8(targetAddr>>8))

			// Setup memory value to be loaded
			mem.Write(targetAddr, tt.value)

			cpu.LoadValueFromMemoryIntoRegisterWithIndexedXIndirectAddress(0x20, tt.register)

			regKey := ""
			switch tt.register {
			case internal.ACCUMULATOR:
				regKey = "acc"
			case internal.REGISTER_X:
				regKey = "x"
			case internal.REGISTER_Y:
				regKey = "y"
			}

			if cpu.GetDebugData()[regKey] != tt.value {
				t.Errorf("Expected %s to be 0x%02X, got 0x%02X", regKey, tt.value, cpu.GetDebugData()[regKey])
			}
		})
	}
}

// Indirect Indexed (Indirect),Y addressing mode
func TestShouldLoadRegisterFromMemoryWithIndirectIndexedYAddress(t *testing.T) {
	tests := []struct {
		name     string
		register string
		value    uint8
	}{
		{
			name:     "Load ACC with Indirect Indexed Y Addressing",
			register: internal.ACCUMULATOR,
			value:    0x42,
		},
		{
			name:     "Load X with Indirect Indexed Y Addressing",
			register: internal.REGISTER_X,
			value:    0x33,
		},
		{
			name:     "Load Y with Indirect Indexed Y Addressing",
			register: internal.REGISTER_Y,
			value:    0x55,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mem := memory.NewMemory()
			cpu := internal.NewCpu(mem)

			// Setup Y register
			var yVal uint8 = 0x05
			cpu.LoadValueIntoRegister(yVal, internal.REGISTER_Y)

			// Base address at 0x20 -> 0x5670
			var baseAddr uint16 = 0x5670
			mem.Write(0x20, uint8(baseAddr&0xFF))
			mem.Write(0x21, uint8(baseAddr>>8))

			// Final address is 0x5670 + 0x05 = 0x5675
			finalAddress := baseAddr + uint16(yVal)

			// Setup memory value to be loaded
			mem.Write(finalAddress, tt.value)

			cpu.LoadValueFromMemoryIntoRegisterWithIndirectIndexedYAddress(0x20, tt.register)

			regKey := ""
			switch tt.register {
			case internal.ACCUMULATOR:
				regKey = "acc"
			case internal.REGISTER_X:
				regKey = "x"
			case internal.REGISTER_Y:
				regKey = "y"
			}

			if cpu.GetDebugData()[regKey] != tt.value {
				t.Errorf("Expected %s to be 0x%02X, got 0x%02X", regKey, tt.value, cpu.GetDebugData()[regKey])
			}
		})
	}
}
