package tests

import (
	"nes-emu/src/emulator/components/bus"
	internal "nes-emu/src/emulator/components/cpu"
	"nes-emu/src/emulator/components/memory"
	"strings"
	"testing"
)

func TestTransferFromAccumulatorToRegisters(t *testing.T) {
	tests := []struct {
		name     string
		constant string
		expected string
	}{
		{
			name:     "transfer value from accumulator to register X",
			constant: internal.REGISTER_X,
			expected: "X",
		},
		{
			name:     "transfer value from accumulator to register Y",
			constant: internal.REGISTER_Y,
			expected: "Y",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mem := memory.NewMemory()
			b := bus.NewBus(mem)
			cpu := internal.NewCpuWithInternal(mem, b)

			cpu.LoadValueIntoRegister(0x42, internal.ACCUMULATOR)
			cpu.TransferFromAccumulatorToRegister(tt.constant)

			debugData := cpu.GetDebugData()

			if debugData[strings.ToLower(tt.constant)] != 0x42 {
				t.Errorf("Expected %s to be 0x42, got %d", tt.constant, debugData[tt.constant])
			}

			if debugData["acc"] != 0x42 {
				t.Errorf("Expected acc to be 0x42, got %d", debugData["acc"])
			}

			if cpu.GetZeroFlag() {
				t.Errorf("Expected zero flag to be false, got true")
			}
			if cpu.GetNegativeFlag() {
				t.Errorf("Expected negative flag to be false, got true")
			}
		})
	}
}

func TestTransferFromAccumulatorToRegistersWithFlaggableValues(t *testing.T) {
	tests := []struct {
		name     string
		constant string
		expected string
	}{
		{
			name:     "transfer value from accumulator to register X",
			constant: internal.REGISTER_X,
			expected: "X",
		},
		{
			name:     "transfer value from accumulator to register Y",
			constant: internal.REGISTER_Y,
			expected: "Y",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mem := memory.NewMemory()
			b := bus.NewBus(mem)
			cpu := internal.NewCpuWithInternal(mem, b)

			cpu.LoadValueIntoRegister(0, internal.ACCUMULATOR)
			cpu.TransferFromAccumulatorToRegister(tt.constant)

			if !cpu.GetZeroFlag() {
				t.Errorf("Expected zero flag to be true, got false")
			}

			if cpu.GetNegativeFlag() {
				t.Errorf("Expected negative flag to be false, got true")
			}

			cpu.LoadValueIntoRegister(254, internal.ACCUMULATOR)
			cpu.TransferFromAccumulatorToRegister(tt.constant)

			if cpu.GetZeroFlag() {
				t.Errorf("Expected zero flag to be false, got true")
			}

			if !cpu.GetNegativeFlag() {
				t.Errorf("Expected negative flag to be true, got false")
			}
		})
	}
}

func TestTransferFromRegisterToAccumulator(t *testing.T) {
	tests := []struct {
		name     string
		constant string
		expected string
	}{
		{
			name:     "transfer value from register X to accumulator",
			constant: internal.REGISTER_X,
			expected: "X",
		},
		{
			name:     "transfer value from register Y to accumulator",
			constant: internal.REGISTER_Y,
			expected: "Y",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mem := memory.NewMemory()
			b := bus.NewBus(mem)
			cpu := internal.NewCpuWithInternal(mem, b)

			cpu.LoadValueIntoRegister(0x42, tt.constant)
			cpu.TransferFromRegisterToAccumulator(tt.constant)

			debugData := cpu.GetDebugData()

			if debugData[strings.ToLower(tt.constant)] != 0x42 {
				t.Errorf("Expected %s to be 0x42, got %d", tt.constant, debugData[tt.constant])
			}

			if debugData["acc"] != 0x42 {
				t.Errorf("Expected acc to be 0x42, got %d", debugData["acc"])
			}

			if cpu.GetZeroFlag() {
				t.Errorf("Expected zero flag to be false, got true")
			}
			if cpu.GetNegativeFlag() {
				t.Errorf("Expected negative flag to be false, got true")
			}
		})
	}
}
