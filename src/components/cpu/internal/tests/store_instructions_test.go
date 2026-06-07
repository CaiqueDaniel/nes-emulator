package tests

import (
	"nes-emu/src/components/cpu/internal"
	"nes-emu/src/components/memory"
	"testing"
)

func TestShouldStoreRegisterIntoAbsoluteMemory(t *testing.T) {
	tests := []struct {
		name     string
		register string
		value    uint8
	}{
		{
			name:     "Store ACC into absolute memory",
			register: internal.ACCUMULATOR,
			value:    0x42,
		},
		{
			name:     "Store REGISTER_X into absolute memory",
			register: internal.REGISTER_X,
			value:    0x33,
		},
		{
			name:     "Store REGISTER_Y into absolute memory",
			register: internal.REGISTER_Y,
			value:    0x55,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mem := memory.NewMemory()
			cpu := internal.NewCpu(mem)
			
			cpu.LoadValueIntoRegister(tt.value, tt.register)
			cpu.StoreRegisterIntoAbsoluteMemory(0x1234, tt.register)

			if mem.Read(0x1234) != tt.value {
				t.Errorf("Expected memory at 0x1234 to be 0x%02X, got 0x%02X", tt.value, mem.Read(0x1234))
			}
		})
	}
}

func TestShouldStoreRegisterIntoZeroPageMemory(t *testing.T) {
	tests := []struct {
		name     string
		register string
		value    uint8
	}{
		{
			name:     "Store ACC into ZeroPage memory",
			register: internal.ACCUMULATOR,
			value:    0x42,
		},
		{
			name:     "Store REGISTER_X into ZeroPage memory",
			register: internal.REGISTER_X,
			value:    0x33,
		},
		{
			name:     "Store REGISTER_Y into ZeroPage memory",
			register: internal.REGISTER_Y,
			value:    0x55,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mem := memory.NewMemory()
			cpu := internal.NewCpu(mem)

			cpu.LoadValueIntoRegister(tt.value, tt.register)
			cpu.StoreRegisterIntoZeroPageMemory(0x34, tt.register)

			if mem.Read(0x0034) != tt.value {
				t.Errorf("Expected memory at 0x0034 to be 0x%02X, got 0x%02X", tt.value, mem.Read(0x0034))
			}
		})
	}
}

func TestShouldStoreAccumulatorIntoMemoryWithIndexAbsolute(t *testing.T) {
	tests := []struct {
		name          string
		indexRegister string
		indexValue    uint8
		baseAddress   uint16
		expectedAddr  uint16
	}{
		{
			name:          "Store ACC with Absolute,X indexing",
			indexRegister: internal.REGISTER_X,
			indexValue:    0x05,
			baseAddress:   0x1200,
			expectedAddr:  0x1205,
		},
		{
			name:          "Store ACC with Absolute,Y indexing",
			indexRegister: internal.REGISTER_Y,
			indexValue:    0x0A,
			baseAddress:   0x1200,
			expectedAddr:  0x120A,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mem := memory.NewMemory()
			cpu := internal.NewCpu(mem)

			accValue := uint8(0x77)
			cpu.LoadValueIntoRegister(accValue, internal.ACCUMULATOR)
			cpu.LoadValueIntoRegister(tt.indexValue, tt.indexRegister)

			cpu.StoreAccumulatorIntoMemoryWithIndexAbsolute(tt.baseAddress, tt.indexRegister)

			if mem.Read(tt.expectedAddr) != accValue {
				t.Errorf("Expected memory at 0x%04X to be 0x%02X, got 0x%02X", tt.expectedAddr, accValue, mem.Read(tt.expectedAddr))
			}
		})
	}
}

func TestShouldStoreAccumulatorIntoMemoryWithIndexZeroPage(t *testing.T) {
	tests := []struct {
		name          string
		indexRegister string
		indexValue    uint8
		baseAddress   uint8
		expectedAddr  uint16
	}{
		{
			name:          "Store ACC with ZeroPage,X indexing",
			indexRegister: internal.REGISTER_X,
			indexValue:    0x05,
			baseAddress:   0x20,
			expectedAddr:  0x0025,
		},
		{
			name:          "Store ACC with ZeroPage,Y indexing",
			indexRegister: internal.REGISTER_Y,
			indexValue:    0x0A,
			baseAddress:   0x20,
			expectedAddr:  0x002A,
		},
		{
			name:          "Store ACC with ZeroPage,X indexing page wrap-around",
			indexRegister: internal.REGISTER_X,
			indexValue:    0xFF,
			baseAddress:   0x01,
			expectedAddr:  0x0000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mem := memory.NewMemory()
			cpu := internal.NewCpu(mem)

			accValue := uint8(0x88)
			cpu.LoadValueIntoRegister(accValue, internal.ACCUMULATOR)
			cpu.LoadValueIntoRegister(tt.indexValue, tt.indexRegister)

			cpu.StoreAccumulatorIntoMemoryWithIndexZeroPage(tt.baseAddress, tt.indexRegister)

			if mem.Read(tt.expectedAddr) != accValue {
				t.Errorf("Expected memory at 0x%04X to be 0x%02X, got 0x%02X", tt.expectedAddr, accValue, mem.Read(tt.expectedAddr))
			}
		})
	}
}

func TestShouldStoreRegisterIntoMemoryWithIndirectAddress(t *testing.T) {
	tests := []struct {
		name     string
		register string
		value    uint8
	}{
		{
			name:     "Store ACC with Indirect Addressing",
			register: internal.ACCUMULATOR,
			value:    0x42,
		},
		{
			name:     "Store X with Indirect Addressing",
			register: internal.REGISTER_X,
			value:    0x33,
		},
		{
			name:     "Store Y with Indirect Addressing",
			register: internal.REGISTER_Y,
			value:    0x55,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mem := memory.NewMemory()
			cpu := internal.NewCpu(mem)

			// target address 0x5678
			mem.Write(0x1000, 0x78)
			mem.Write(0x1001, 0x56)

			cpu.LoadValueIntoRegister(tt.value, tt.register)
			cpu.StoreRegisterIntoMemoryWithIndirectAddress(0x1000, tt.register)

			if mem.Read(0x5678) != tt.value {
				t.Errorf("Expected memory at 0x5678 to be 0x%02X, got 0x%02X", tt.value, mem.Read(0x5678))
			}
		})
	}
}

func TestShouldStoreRegisterIntoMemoryWithIndexedXIndirectAddress(t *testing.T) {
	tests := []struct {
		name     string
		register string
		value    uint8
	}{
		{
			name:     "Store ACC with Indexed X Indirect Addressing",
			register: internal.ACCUMULATOR,
			value:    0x42,
		},
		{
			name:     "Store X with Indexed X Indirect Addressing",
			register: internal.REGISTER_X,
			value:    0x33,
		},
		{
			name:     "Store Y with Indexed X Indirect Addressing",
			register: internal.REGISTER_Y,
			value:    0x55,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mem := memory.NewMemory()
			cpu := internal.NewCpu(mem)

			var xVal uint8 = 0x04
			var targetAddr uint16 = 0x5678
			var expectedValue uint8 = tt.value

			if tt.register == internal.REGISTER_X {
				xVal = tt.value
				pivotAddress := uint16(0x20) + uint16(tt.value)
				mem.Write(pivotAddress, uint8(targetAddr&0xFF))
				mem.Write(pivotAddress+1, uint8(targetAddr>>8))
			} else {
				mem.Write(0x24, 0x78)
				mem.Write(0x25, 0x56)
			}

			cpu.LoadValueIntoRegister(xVal, internal.REGISTER_X)
			if tt.register != internal.REGISTER_X {
				cpu.LoadValueIntoRegister(tt.value, tt.register)
			}

			cpu.StoreRegisterIntoMemoryWithIndexedXIndirectAddress(0x20, tt.register)

			if mem.Read(targetAddr) != expectedValue {
				t.Errorf("Expected memory at 0x%04X to be 0x%02X, got 0x%02X", targetAddr, expectedValue, mem.Read(targetAddr))
			}
		})
	}
}

func TestShouldStoreRegisterIntoMemoryWithIndirectIndexedYAddress(t *testing.T) {
	tests := []struct {
		name     string
		register string
		value    uint8
	}{
		{
			name:     "Store ACC with Indirect Indexed Y Addressing",
			register: internal.ACCUMULATOR,
			value:    0x42,
		},
		{
			name:     "Store X with Indirect Indexed Y Addressing",
			register: internal.REGISTER_X,
			value:    0x33,
		},
		{
			name:     "Store Y with Indirect Indexed Y Addressing",
			register: internal.REGISTER_Y,
			value:    0x55,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mem := memory.NewMemory()
			cpu := internal.NewCpu(mem)

			var yVal uint8 = 0x05
			var baseAddr uint16 = 0x5670
			var expectedValue uint8 = tt.value

			if tt.register == internal.REGISTER_Y {
				yVal = tt.value
			}

			mem.Write(0x20, uint8(baseAddr&0xFF))
			mem.Write(0x21, uint8(baseAddr>>8))

			finalAddress := baseAddr + uint16(yVal)

			cpu.LoadValueIntoRegister(yVal, internal.REGISTER_Y)
			if tt.register != internal.REGISTER_Y {
				cpu.LoadValueIntoRegister(tt.value, tt.register)
			}

			cpu.StoreRegisterIntoMemoryWithIndirectIndexedYAddress(0x20, tt.register)

			if mem.Read(finalAddress) != expectedValue {
				t.Errorf("Expected memory at 0x%04X to be 0x%02X, got 0x%02X", finalAddress, expectedValue, mem.Read(finalAddress))
			}
		})
	}
}
