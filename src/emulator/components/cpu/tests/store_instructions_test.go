package tests

import (
	internal "nes-emu/src/emulator/components/cpu"
	"nes-emu/src/emulator/components/memory"
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
