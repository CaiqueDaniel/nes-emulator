package tests

import (
	"nes-emu/src/emulator/components/bus"
	internal "nes-emu/src/emulator/components/cpu"
	"nes-emu/src/emulator/components/memory"
	"testing"
)

func TestAnd(t *testing.T) {
	tests := []struct {
		name        string
		initialAcc  uint8
		value       uint8
		expectedAcc uint8
		expectedZ   bool
		expectedN   bool
	}{
		{
			name:        "AND with same value keeps bits",
			initialAcc:  0b11001100,
			value:       0b11001100,
			expectedAcc: 0b11001100,
			expectedZ:   false,
			expectedN:   true,
		},
		{
			name:        "AND with zero clears accumulator",
			initialAcc:  0xFF,
			value:       0x00,
			expectedAcc: 0x00,
			expectedZ:   true,
			expectedN:   false,
		},
		{
			name:        "AND masks out lower nibble",
			initialAcc:  0b11111111,
			value:       0b11110000,
			expectedAcc: 0b11110000,
			expectedZ:   false,
			expectedN:   true,
		},
		{
			name:        "AND result is positive (bit 7 clear)",
			initialAcc:  0b11111111,
			value:       0b01111111,
			expectedAcc: 0b01111111,
			expectedZ:   false,
			expectedN:   false,
		},
		{
			name:        "AND with 0xFF keeps accumulator unchanged",
			initialAcc:  0b10101010,
			value:       0xFF,
			expectedAcc: 0b10101010,
			expectedZ:   false,
			expectedN:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mem := memory.NewMemory()
			b := bus.NewBusWithWorkMemory(mem)
			cpu := internal.NewCpuWithInternal(b)

			cpu.LoadValueIntoRegister(tt.initialAcc, internal.ACCUMULATOR)
			cpu.And(tt.value)

			debugData := cpu.GetDebugData()
			if debugData["acc"] != tt.expectedAcc {
				t.Errorf("Expected accumulator %08b, got %08b", tt.expectedAcc, debugData["acc"])
			}
			if cpu.GetZeroFlag() != tt.expectedZ {
				t.Errorf("Expected Zero flag %t, got %t", tt.expectedZ, cpu.GetZeroFlag())
			}
			if cpu.GetNegativeFlag() != tt.expectedN {
				t.Errorf("Expected Negative flag %t, got %t", tt.expectedN, cpu.GetNegativeFlag())
			}
		})
	}
}

func TestOr(t *testing.T) {
	tests := []struct {
		name        string
		initialAcc  uint8
		value       uint8
		expectedAcc uint8
		expectedZ   bool
		expectedN   bool
	}{
		{
			name:        "OR with zero keeps accumulator",
			initialAcc:  0b10101010,
			value:       0x00,
			expectedAcc: 0b10101010,
			expectedZ:   false,
			expectedN:   true,
		},
		{
			name:        "OR with 0xFF sets all bits",
			initialAcc:  0x00,
			value:       0xFF,
			expectedAcc: 0xFF,
			expectedZ:   false,
			expectedN:   true,
		},
		{
			name:        "OR zero with zero stays zero",
			initialAcc:  0x00,
			value:       0x00,
			expectedAcc: 0x00,
			expectedZ:   true,
			expectedN:   false,
		},
		{
			name:        "OR merges two nibbles",
			initialAcc:  0b11110000,
			value:       0b00001111,
			expectedAcc: 0b11111111,
			expectedZ:   false,
			expectedN:   true,
		},
		{
			name:        "OR result is positive (bit 7 clear)",
			initialAcc:  0b00001111,
			value:       0b00110000,
			expectedAcc: 0b00111111,
			expectedZ:   false,
			expectedN:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mem := memory.NewMemory()
			b := bus.NewBusWithWorkMemory(mem)
			cpu := internal.NewCpuWithInternal(b)

			cpu.LoadValueIntoRegister(tt.initialAcc, internal.ACCUMULATOR)
			cpu.Or(tt.value)

			debugData := cpu.GetDebugData()
			if debugData["acc"] != tt.expectedAcc {
				t.Errorf("Expected accumulator %08b, got %08b", tt.expectedAcc, debugData["acc"])
			}
			if cpu.GetZeroFlag() != tt.expectedZ {
				t.Errorf("Expected Zero flag %t, got %t", tt.expectedZ, cpu.GetZeroFlag())
			}
			if cpu.GetNegativeFlag() != tt.expectedN {
				t.Errorf("Expected Negative flag %t, got %t", tt.expectedN, cpu.GetNegativeFlag())
			}
		})
	}
}

func TestXor(t *testing.T) {
	tests := []struct {
		name        string
		initialAcc  uint8
		value       uint8
		expectedAcc uint8
		expectedZ   bool
		expectedN   bool
	}{
		{
			name:        "XOR with same value clears accumulator (zero flag)",
			initialAcc:  0b10101010,
			value:       0b10101010,
			expectedAcc: 0x00,
			expectedZ:   true,
			expectedN:   false,
		},
		{
			name:        "XOR with 0xFF inverts all bits",
			initialAcc:  0b01010101,
			value:       0xFF,
			expectedAcc: 0b10101010,
			expectedZ:   false,
			expectedN:   true,
		},
		{
			name:        "XOR zero with zero stays zero",
			initialAcc:  0x00,
			value:       0x00,
			expectedAcc: 0x00,
			expectedZ:   true,
			expectedN:   false,
		},
		{
			name:        "XOR toggles specific bits",
			initialAcc:  0b11110000,
			value:       0b11001100,
			expectedAcc: 0b00111100,
			expectedZ:   false,
			expectedN:   false,
		},
		{
			name:        "XOR result is negative (bit 7 set)",
			initialAcc:  0b00001111,
			value:       0b10000000,
			expectedAcc: 0b10001111,
			expectedZ:   false,
			expectedN:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mem := memory.NewMemory()
			b := bus.NewBusWithWorkMemory(mem)
			cpu := internal.NewCpuWithInternal(b)

			cpu.LoadValueIntoRegister(tt.initialAcc, internal.ACCUMULATOR)
			cpu.Xor(tt.value)

			debugData := cpu.GetDebugData()
			if debugData["acc"] != tt.expectedAcc {
				t.Errorf("Expected accumulator %08b, got %08b", tt.expectedAcc, debugData["acc"])
			}
			if cpu.GetZeroFlag() != tt.expectedZ {
				t.Errorf("Expected Zero flag %t, got %t", tt.expectedZ, cpu.GetZeroFlag())
			}
			if cpu.GetNegativeFlag() != tt.expectedN {
				t.Errorf("Expected Negative flag %t, got %t", tt.expectedN, cpu.GetNegativeFlag())
			}
		})
	}
}

func TestBit(t *testing.T) {
	tests := []struct {
		name       string
		initialAcc uint8
		value      uint8
		expectedZ  bool
		expectedN  bool
		expectedV  bool
	}{
		{
			name:       "BIT: result is zero (no shared bits)",
			initialAcc: 0b00001111,
			value:      0b11110000,
			expectedZ:  true,
			expectedN:  false,
			expectedV:  false,
		},
		{
			name:       "BIT: result has bit 7 set (negative flag)",
			initialAcc: 0xFF,
			value:      0b10000000,
			expectedZ:  false,
			expectedN:  true,
			expectedV:  false,
		},
		{
			name:       "BIT: result has bit 6 set (overflow flag)",
			initialAcc: 0xFF,
			value:      0b01000000,
			expectedZ:  false,
			expectedN:  false,
			expectedV:  true,
		},
		{
			name:       "BIT: result has bits 6 and 7 set",
			initialAcc: 0xFF,
			value:      0b11000000,
			expectedZ:  false,
			expectedN:  true,
			expectedV:  true,
		},
		{
			name:       "BIT: AND with zero clears all flags (only zero flag)",
			initialAcc: 0x00,
			value:      0xFF,
			expectedZ:  true,
			expectedN:  false,
			expectedV:  false,
		},
		{
			name:       "BIT: sets both negative and overflow flags",
			initialAcc: 0b11001100,
			value:      0b11000000,
			expectedZ:  false,
			expectedN:  true,
			expectedV:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mem := memory.NewMemory()
			b := bus.NewBusWithWorkMemory(mem)
			cpu := internal.NewCpuWithInternal(b)

			cpu.LoadValueIntoRegister(tt.initialAcc, internal.ACCUMULATOR)
			cpu.Bit(tt.value)

			// BIT must NOT modify the accumulator
			debugData := cpu.GetDebugData()
			if debugData["acc"] != tt.initialAcc {
				t.Errorf("BIT should not modify accumulator: expected %08b, got %08b", tt.initialAcc, debugData["acc"])
			}
			if cpu.GetZeroFlag() != tt.expectedZ {
				t.Errorf("Expected Zero flag %t, got %t", tt.expectedZ, cpu.GetZeroFlag())
			}
			if cpu.GetNegativeFlag() != tt.expectedN {
				t.Errorf("Expected Negative flag %t, got %t", tt.expectedN, cpu.GetNegativeFlag())
			}
			if cpu.GetOverflowFlag() != tt.expectedV {
				t.Errorf("Expected Overflow flag %t, got %t", tt.expectedV, cpu.GetOverflowFlag())
			}
		})
	}
}
