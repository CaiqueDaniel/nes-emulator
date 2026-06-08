package tests

import (
	"nes-emu/src/components/cpu/internal"
	"nes-emu/src/components/memory"
	"testing"
)

func TestAddWithCarry(t *testing.T) {
	tests := []struct {
		name         string
		initialAcc   uint8
		valueToAdd   uint8
		initialCarry bool
		expectedAcc  uint8
		expectedC    bool
		expectedZ    bool
		expectedN    bool
		expectedV    bool
	}{
		{
			name:         "Simple positive addition, no carry or overflow",
			initialAcc:   10,
			valueToAdd:   20,
			initialCarry: false,
			expectedAcc:  30,
			expectedC:    false,
			expectedZ:    false,
			expectedN:    false,
			expectedV:    false,
		},
		{
			name:         "Simple positive addition with carry-in, no overflow",
			initialAcc:   10,
			valueToAdd:   20,
			initialCarry: true,
			expectedAcc:  31,
			expectedC:    false,
			expectedZ:    false,
			expectedN:    false,
			expectedV:    false,
		},
		{
			name:         "Carry-out set (exceeds 255), no overflow",
			initialAcc:   255,
			valueToAdd:   1,
			initialCarry: false,
			expectedAcc:  0,
			expectedC:    true,
			expectedZ:    true,
			expectedN:    false,
			expectedV:    false,
		},
		{
			name:         "Carry-out set with carry-in, no overflow",
			initialAcc:   255,
			valueToAdd:   5,
			initialCarry: true,
			expectedAcc:  5,
			expectedC:    true,
			expectedZ:    false,
			expectedN:    false,
			expectedV:    false,
		},
		{
			name:         "Signed overflow (Positive + Positive = Negative)",
			initialAcc:   127,
			valueToAdd:   1,
			initialCarry: false,
			expectedAcc:  128,
			expectedC:    false,
			expectedZ:    false,
			expectedN:    true,
			expectedV:    true,
		},
		{
			name:         "Signed overflow (Negative + Negative = Positive)",
			initialAcc:   128,
			valueToAdd:   128,
			initialCarry: false,
			expectedAcc:  0,
			expectedC:    true,
			expectedZ:    true,
			expectedN:    false,
			expectedV:    true,
		},
		{
			name:         "Large negative inputs, no overflow",
			initialAcc:   255,
			valueToAdd:   255,
			initialCarry: false,
			expectedAcc:  254,
			expectedC:    true,
			expectedZ:    false,
			expectedN:    true,
			expectedV:    false,
		},
		{
			name:         "Positive + Negative, negative result, no overflow",
			initialAcc:   1,
			valueToAdd:   254,
			initialCarry: false,
			expectedAcc:  255,
			expectedC:    false,
			expectedZ:    false,
			expectedN:    true,
			expectedV:    false,
		},
		{
			name:         "Zero flag set",
			initialAcc:   0,
			valueToAdd:   0,
			initialCarry: false,
			expectedAcc:  0,
			expectedC:    false,
			expectedZ:    true,
			expectedN:    false,
			expectedV:    false,
		},
		{
			name:         "Negative flag set without overflow",
			initialAcc:   0,
			valueToAdd:   130,
			initialCarry: false,
			expectedAcc:  130,
			expectedC:    false,
			expectedZ:    false,
			expectedN:    true,
			expectedV:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mem := memory.NewMemory()
			cpu := internal.NewCpu(mem)

			cpu.LoadValueIntoRegister(tt.initialAcc, internal.ACCUMULATOR)
			if tt.initialCarry {
				cpu.SetCarryFlag()
			} else {
				cpu.ClearCarryFlag()
			}

			cpu.AddWithCarry(tt.valueToAdd)

			debugData := cpu.GetDebugData()
			actualAcc := debugData["acc"]

			if actualAcc != tt.expectedAcc {
				t.Errorf("Expected accumulator value %d, got %d", tt.expectedAcc, actualAcc)
			}
			if cpu.GetCarryFlag() != tt.expectedC {
				t.Errorf("Expected Carry flag %t, got %t", tt.expectedC, cpu.GetCarryFlag())
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

func TestSubtractWithCarry(t *testing.T) {
	tests := []struct {
		name         string
		initialAcc   uint8
		valueToSub   uint8
		initialCarry bool
		expectedAcc  uint8
		expectedC    bool
		expectedZ    bool
		expectedN    bool
		expectedV    bool
	}{
		{
			name:         "Simple subtraction, no borrow (carry-in true)",
			initialAcc:   50,
			valueToSub:   20,
			initialCarry: true,
			expectedAcc:  30,
			expectedC:    true,
			expectedZ:    false,
			expectedN:    false,
			expectedV:    false,
		},
		{
			name:         "Subtraction with borrow-in (carry-in false)",
			initialAcc:   50,
			valueToSub:   20,
			initialCarry: false,
			expectedAcc:  29,
			expectedC:    true,
			expectedZ:    false,
			expectedN:    false,
			expectedV:    false,
		},
		{
			name:         "Subtraction causing borrow-out (result < 0)",
			initialAcc:   10,
			valueToSub:   20,
			initialCarry: true,
			expectedAcc:  246,
			expectedC:    false,
			expectedZ:    false,
			expectedN:    true,
			expectedV:    false,
		},
		{
			name:         "Subtraction causing borrow-out with borrow-in",
			initialAcc:   10,
			valueToSub:   20,
			initialCarry: false,
			expectedAcc:  245,
			expectedC:    false,
			expectedZ:    false,
			expectedN:    true,
			expectedV:    false,
		},
		{
			name:         "Signed overflow (Positive - Negative = Negative)",
			initialAcc:   127,
			valueToSub:   255, // -1
			initialCarry: true,
			expectedAcc:  128,
			expectedC:    false,
			expectedZ:    false,
			expectedN:    true,
			expectedV:    true,
		},
		{
			name:         "Signed overflow (Negative - Positive = Positive)",
			initialAcc:   128,
			valueToSub:   1,
			initialCarry: true,
			expectedAcc:  127,
			expectedC:    true,
			expectedZ:    false,
			expectedN:    false,
			expectedV:    true,
		},
		{
			name:         "Zero result",
			initialAcc:   20,
			valueToSub:   20,
			initialCarry: true,
			expectedAcc:  0,
			expectedC:    true,
			expectedZ:    true,
			expectedN:    false,
			expectedV:    false,
		},
		{
			name:         "Negative result without overflow",
			initialAcc:   10,
			valueToSub:   15,
			initialCarry: true,
			expectedAcc:  251,
			expectedC:    false,
			expectedZ:    false,
			expectedN:    true,
			expectedV:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mem := memory.NewMemory()
			cpu := internal.NewCpu(mem)

			cpu.LoadValueIntoRegister(tt.initialAcc, internal.ACCUMULATOR)
			if tt.initialCarry {
				cpu.SetCarryFlag()
			} else {
				cpu.ClearCarryFlag()
			}

			cpu.SubtractWithCarry(tt.valueToSub)

			debugData := cpu.GetDebugData()
			actualAcc := debugData["acc"]

			if actualAcc != tt.expectedAcc {
				t.Errorf("Expected accumulator value %d, got %d", tt.expectedAcc, actualAcc)
			}
			if cpu.GetCarryFlag() != tt.expectedC {
				t.Errorf("Expected Carry flag %t, got %t", tt.expectedC, cpu.GetCarryFlag())
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
			cpu := internal.NewCpu(mem)

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
			cpu := internal.NewCpu(mem)

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
			cpu := internal.NewCpu(mem)

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
			cpu := internal.NewCpu(mem)

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

func TestCompareWithRegister(t *testing.T) {
	tests := []struct {
		name       string
		initialAcc uint8
		value      uint8
		expectedZ  bool
		expectedC  bool
		expectedN  bool
	}{
		{
			name:       "Equal values: zero and carry set, no negative",
			initialAcc: 0x10,
			value:      0x10,
			expectedZ:  true,
			expectedC:  true,
			expectedN:  false,
		},
		{
			name:       "ACC greater than value: carry set, no zero or negative",
			initialAcc: 0x20,
			value:      0x10,
			expectedZ:  false,
			expectedC:  true,
			expectedN:  false,
		},
		{
			name:       "ACC less than value: no zero or carry",
			initialAcc: 0x10,
			value:      0x20,
			expectedZ:  false,
			expectedC:  false,
			expectedN:  false,
		},
		{
			name:       "Both zero: zero and carry set",
			initialAcc: 0x00,
			value:      0x00,
			expectedZ:  true,
			expectedC:  true,
			expectedN:  false,
		},
		{
			name:       "Negative flag set when bit 7 of (acc AND value) is 1",
			initialAcc: 0b10000001,
			value:      0b10000010,
			expectedZ:  false,
			expectedC:  false,
			expectedN:  true,
		},
		{
			name:       "Negative flag clear when bit 7 of (acc AND value) is 0",
			initialAcc: 0b11000000,
			value:      0b01000000,
			expectedZ:  false,
			expectedC:  true,
			expectedN:  false,
		},
		{
			name:       "ACC equals 0xFF, value equals 0xFF",
			initialAcc: 0xFF,
			value:      0xFF,
			expectedZ:  true,
			expectedC:  true,
			expectedN:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mem := memory.NewMemory()
			cpu := internal.NewCpu(mem)

			cpu.LoadValueIntoRegister(tt.initialAcc, internal.ACCUMULATOR)
			cpu.CompareWithRegister(tt.value)

			// CompareWithRegister must NOT modify the accumulator
			debugData := cpu.GetDebugData()
			if debugData["acc"] != tt.initialAcc {
				t.Errorf("CompareWithRegister should not modify accumulator: expected %d, got %d", tt.initialAcc, debugData["acc"])
			}
			if cpu.GetZeroFlag() != tt.expectedZ {
				t.Errorf("Expected Zero flag %t, got %t", tt.expectedZ, cpu.GetZeroFlag())
			}
			if cpu.GetCarryFlag() != tt.expectedC {
				t.Errorf("Expected Carry flag %t, got %t", tt.expectedC, cpu.GetCarryFlag())
			}
			if cpu.GetNegativeFlag() != tt.expectedN {
				t.Errorf("Expected Negative flag %t, got %t", tt.expectedN, cpu.GetNegativeFlag())
			}
		})
	}
}
