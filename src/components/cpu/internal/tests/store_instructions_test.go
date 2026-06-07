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
