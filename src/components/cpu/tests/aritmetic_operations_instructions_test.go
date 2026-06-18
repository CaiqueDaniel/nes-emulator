package tests

import (
	internal "nes-emu/src/components/cpu"
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

func TestCompareWithRegister(t *testing.T) {
	tests := []struct {
		name            string
		register        string
		initialRegister uint8
		value           uint8
		expectedZ       bool
		expectedC       bool
		expectedN       bool
	}{
		// Accumulator (CMP)
		{
			name:            "ACC equal: zero and carry set, no negative",
			register:        internal.ACCUMULATOR,
			initialRegister: 0x10,
			value:           0x10,
			expectedZ:       true,
			expectedC:       true,
			expectedN:       false,
		},
		{
			name:            "ACC greater than value: carry set only",
			register:        internal.ACCUMULATOR,
			initialRegister: 0x20,
			value:           0x10,
			expectedZ:       false,
			expectedC:       true,
			expectedN:       false,
		},
		{
			name:            "ACC less than value: no flags set",
			register:        internal.ACCUMULATOR,
			initialRegister: 0x10,
			value:           0x20,
			expectedZ:       false,
			expectedC:       false,
			expectedN:       false,
		},
		{
			name:            "ACC both zero: zero and carry set",
			register:        internal.ACCUMULATOR,
			initialRegister: 0x00,
			value:           0x00,
			expectedZ:       true,
			expectedC:       true,
			expectedN:       false,
		},
		{
			name:            "ACC: negative flag set when bit 7 of (reg AND value) is 1",
			register:        internal.ACCUMULATOR,
			initialRegister: 0b10000001,
			value:           0b10000010,
			expectedZ:       false,
			expectedC:       false,
			expectedN:       true,
		},
		{
			name:            "ACC: negative flag clear when bit 7 of (reg AND value) is 0",
			register:        internal.ACCUMULATOR,
			initialRegister: 0b11000000,
			value:           0b01000000,
			expectedZ:       false,
			expectedC:       true,
			expectedN:       false,
		},
		{
			name:            "ACC equals 0xFF, value equals 0xFF",
			register:        internal.ACCUMULATOR,
			initialRegister: 0xFF,
			value:           0xFF,
			expectedZ:       true,
			expectedC:       true,
			expectedN:       true,
		},
		// Register X (CPX)
		{
			name:            "X equal: zero and carry set",
			register:        internal.REGISTER_X,
			initialRegister: 0x42,
			value:           0x42,
			expectedZ:       true,
			expectedC:       true,
			expectedN:       false,
		},
		{
			name:            "X greater than value: carry set only",
			register:        internal.REGISTER_X,
			initialRegister: 0x50,
			value:           0x30,
			expectedZ:       false,
			expectedC:       true,
			expectedN:       false,
		},
		{
			name:            "X less than value: no flags set",
			register:        internal.REGISTER_X,
			initialRegister: 0x10,
			value:           0x40,
			expectedZ:       false,
			expectedC:       false,
			expectedN:       false,
		},
		// Register Y (CPY)
		{
			name:            "Y equal: zero and carry set",
			register:        internal.REGISTER_Y,
			initialRegister: 0x77,
			value:           0x77,
			expectedZ:       true,
			expectedC:       true,
			expectedN:       false,
		},
		{
			name:            "Y greater than value: carry set only",
			register:        internal.REGISTER_Y,
			initialRegister: 0x80,
			value:           0x01,
			expectedZ:       false,
			expectedC:       true,
			expectedN:       false,
		},
		{
			name:            "Y less than value: no flags set",
			register:        internal.REGISTER_Y,
			initialRegister: 0x01,
			value:           0x80,
			expectedZ:       false,
			expectedC:       false,
			expectedN:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mem := memory.NewMemory()
			cpu := internal.NewCpu(mem)

			cpu.LoadValueIntoRegister(tt.initialRegister, tt.register)
			cpu.CompareWithRegister(tt.value, tt.register)

			// CompareWithRegister must NOT modify the register
			debugData := cpu.GetDebugData()
			var registerKey string
			switch tt.register {
			case internal.ACCUMULATOR:
				registerKey = "acc"
			case internal.REGISTER_X:
				registerKey = "x"
			case internal.REGISTER_Y:
				registerKey = "y"
			}
			if debugData[registerKey] != tt.initialRegister {
				t.Errorf("CompareWithRegister should not modify register %s: expected %d, got %d", tt.register, tt.initialRegister, debugData[registerKey])
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

func TestIncrementMemory(t *testing.T) {
	tests := []struct {
		name          string
		address       uint16
		initialValue  uint8
		expectedValue uint8
		expectedZ     bool
		expectedN     bool
	}{
		{
			name:          "Increment normal value",
			address:       0x0010,
			initialValue:  0x05,
			expectedValue: 0x06,
			expectedZ:     false,
			expectedN:     false,
		},
		{
			name:          "Increment wraps from 0xFF to 0x00 and sets zero flag",
			address:       0x0020,
			initialValue:  0xFF,
			expectedValue: 0x00,
			expectedZ:     true,
			expectedN:     false,
		},
		{
			name:          "Increment sets negative flag when result bit 7 is set",
			address:       0x0030,
			initialValue:  0x7F,
			expectedValue: 0x80,
			expectedZ:     false,
			expectedN:     true,
		},
		{
			name:          "Increment value already in negative range stays negative",
			address:       0x0040,
			initialValue:  0x80,
			expectedValue: 0x81,
			expectedZ:     false,
			expectedN:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mem := memory.NewMemory()
			cpu := internal.NewCpu(mem)

			mem.Write(tt.address, tt.initialValue)
			cpu.IncrementMemory(tt.address)

			actualValue := mem.Read(tt.address)
			if actualValue != tt.expectedValue {
				t.Errorf("Expected memory value %d, got %d", tt.expectedValue, actualValue)
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

func TestIncrementRegister(t *testing.T) {
	tests := []struct {
		name          string
		register      string
		initialValue  uint8
		expectedValue uint8
		expectedZ     bool
		expectedN     bool
	}{
		{
			name:          "Increment X normal value",
			register:      internal.REGISTER_X,
			initialValue:  0x05,
			expectedValue: 0x06,
			expectedZ:     false,
			expectedN:     false,
		},
		{
			name:          "Increment X wraps from 0xFF to 0x00 and sets zero flag",
			register:      internal.REGISTER_X,
			initialValue:  0xFF,
			expectedValue: 0x00,
			expectedZ:     true,
			expectedN:     false,
		},
		{
			name:          "Increment X sets negative flag",
			register:      internal.REGISTER_X,
			initialValue:  0x7F,
			expectedValue: 0x80,
			expectedZ:     false,
			expectedN:     true,
		},
		{
			name:          "Increment Y normal value",
			register:      internal.REGISTER_Y,
			initialValue:  0x10,
			expectedValue: 0x11,
			expectedZ:     false,
			expectedN:     false,
		},
		{
			name:          "Increment Y wraps from 0xFF to 0x00 and sets zero flag",
			register:      internal.REGISTER_Y,
			initialValue:  0xFF,
			expectedValue: 0x00,
			expectedZ:     true,
			expectedN:     false,
		},
		{
			name:          "Increment Y sets negative flag",
			register:      internal.REGISTER_Y,
			initialValue:  0x7F,
			expectedValue: 0x80,
			expectedZ:     false,
			expectedN:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mem := memory.NewMemory()
			cpu := internal.NewCpu(mem)

			cpu.LoadValueIntoRegister(tt.initialValue, tt.register)
			cpu.IncrementRegister(tt.register)

			debugData := cpu.GetDebugData()
			var registerKey string
			switch tt.register {
			case internal.REGISTER_X:
				registerKey = "x"
			case internal.REGISTER_Y:
				registerKey = "y"
			}

			if debugData[registerKey] != tt.expectedValue {
				t.Errorf("Expected register %s value %d, got %d", tt.register, tt.expectedValue, debugData[registerKey])
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

func TestDecrementMemory(t *testing.T) {
	tests := []struct {
		name          string
		address       uint16
		initialValue  uint8
		expectedValue uint8
		expectedZ     bool
		expectedN     bool
	}{
		{
			name:          "Decrement normal value",
			address:       0x0010,
			initialValue:  0x05,
			expectedValue: 0x04,
			expectedZ:     false,
			expectedN:     false,
		},
		{
			name:          "Decrement sets zero flag when result is 0x00",
			address:       0x0020,
			initialValue:  0x01,
			expectedValue: 0x00,
			expectedZ:     true,
			expectedN:     false,
		},
		{
			name:          "Decrement wraps from 0x00 to 0xFF and sets negative flag",
			address:       0x0030,
			initialValue:  0x00,
			expectedValue: 0xFF,
			expectedZ:     false,
			expectedN:     true,
		},
		{
			name:          "Decrement sets negative flag when result bit 7 is set",
			address:       0x0040,
			initialValue:  0x81,
			expectedValue: 0x80,
			expectedZ:     false,
			expectedN:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mem := memory.NewMemory()
			cpu := internal.NewCpu(mem)

			mem.Write(tt.address, tt.initialValue)
			cpu.DecrementMemory(tt.address)

			actualValue := mem.Read(tt.address)
			if actualValue != tt.expectedValue {
				t.Errorf("Expected memory value %d, got %d", tt.expectedValue, actualValue)
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

func TestDecrementRegister(t *testing.T) {
	tests := []struct {
		name          string
		register      string
		initialValue  uint8
		expectedValue uint8
		expectedZ     bool
		expectedN     bool
	}{
		{
			name:          "Decrement X normal value",
			register:      internal.REGISTER_X,
			initialValue:  0x05,
			expectedValue: 0x04,
			expectedZ:     false,
			expectedN:     false,
		},
		{
			name:          "Decrement X sets zero flag",
			register:      internal.REGISTER_X,
			initialValue:  0x01,
			expectedValue: 0x00,
			expectedZ:     true,
			expectedN:     false,
		},
		{
			name:          "Decrement X wraps from 0x00 to 0xFF and sets negative flag",
			register:      internal.REGISTER_X,
			initialValue:  0x00,
			expectedValue: 0xFF,
			expectedZ:     false,
			expectedN:     true,
		},
		{
			name:          "Decrement Y normal value",
			register:      internal.REGISTER_Y,
			initialValue:  0x10,
			expectedValue: 0x0F,
			expectedZ:     false,
			expectedN:     false,
		},
		{
			name:          "Decrement Y sets zero flag",
			register:      internal.REGISTER_Y,
			initialValue:  0x01,
			expectedValue: 0x00,
			expectedZ:     true,
			expectedN:     false,
		},
		{
			name:          "Decrement Y wraps from 0x00 to 0xFF and sets negative flag",
			register:      internal.REGISTER_Y,
			initialValue:  0x00,
			expectedValue: 0xFF,
			expectedZ:     false,
			expectedN:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mem := memory.NewMemory()
			cpu := internal.NewCpu(mem)

			cpu.LoadValueIntoRegister(tt.initialValue, tt.register)
			cpu.DecrementRegister(tt.register)

			debugData := cpu.GetDebugData()
			var registerKey string
			switch tt.register {
			case internal.REGISTER_X:
				registerKey = "x"
			case internal.REGISTER_Y:
				registerKey = "y"
			}

			if debugData[registerKey] != tt.expectedValue {
				t.Errorf("Expected register %s value %d, got %d", tt.register, tt.expectedValue, debugData[registerKey])
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
