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
