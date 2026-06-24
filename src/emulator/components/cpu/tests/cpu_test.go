package tests

import (
	"nes-emu/src/emulator/components/bus"
	"nes-emu/src/emulator/components/cpu"
	internal "nes-emu/src/emulator/components/cpu"
	"nes-emu/src/emulator/components/memory"
	"testing"
)

// Register Constants Tests
func TestRegisterConstantsValues(t *testing.T) {
	tests := []struct {
		name     string
		constant string
		expected string
	}{
		{
			name:     "REGISTER_X constant value",
			constant: internal.REGISTER_X,
			expected: "X",
		},
		{
			name:     "REGISTER_Y constant value",
			constant: internal.REGISTER_Y,
			expected: "Y",
		},
		{
			name:     "ACCUMULATOR constant value",
			constant: internal.ACCUMULATOR,
			expected: "ACC",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.constant != tt.expected {
				t.Errorf("Expected %s to be %s, got %s", tt.name, tt.expected, tt.constant)
			}
		})
	}
}

// CPU Initialization and Reset Tests
func TestNewCpuInitialization(t *testing.T) {
	b := bus.NewBus()
	memory := memory.NewMemory(b)
	cpu := internal.NewCpu(memory, b)

	debugData := cpu.GetDebugData()

	if debugData["acc"] != 0 {
		t.Errorf("Expected acc to be 0 on initialization, got %d", debugData["acc"])
	}
	if debugData["x"] != 0 {
		t.Errorf("Expected x to be 0 on initialization, got %d", debugData["x"])
	}
	if debugData["y"] != 0 {
		t.Errorf("Expected y to be 0 on initialization, got %d", debugData["y"])
	}
}

func TestCpuResetClearsAllRegisters(t *testing.T) {
	b := bus.NewBus()
	memory := memory.NewMemory(b)
	cpu := internal.NewCpu(memory, b)

	// Set some values to non-zero
	cpu.LoadValueIntoRegister(0xFF, internal.ACCUMULATOR)
	cpu.LoadValueIntoRegister(0xAA, internal.REGISTER_X)
	cpu.LoadValueIntoRegister(0xBB, internal.REGISTER_Y)

	// Verify values were set
	debugData := cpu.GetDebugData()
	if debugData["acc"] != 0xFF {
		t.Errorf("Expected acc to be 0xFF before reset, got %d", debugData["acc"])
	}
	if debugData["x"] != 0xAA {
		t.Errorf("Expected x to be 0xAA before reset, got %d", debugData["x"])
	}
	if debugData["y"] != 0xBB {
		t.Errorf("Expected y to be 0xBB before reset, got %d", debugData["y"])
	}

	// Reset CPU
	cpu.Reset()

	// Verify all registers are cleared
	debugData = cpu.GetDebugData()
	if debugData["acc"] != 0 {
		t.Errorf("Expected acc to be 0 after reset, got %d", debugData["acc"])
	}
	if debugData["x"] != 0 {
		t.Errorf("Expected x to be 0 after reset, got %d", debugData["x"])
	}
	if debugData["y"] != 0 {
		t.Errorf("Expected y to be 0 after reset, got %d", debugData["y"])
	}
}

func TestGetDebugDataReturnsCurrentRegisterState(t *testing.T) {
	b := bus.NewBus()
	memory := memory.NewMemory(b)
	cpu := internal.NewCpu(memory, b)

	// Set specific values
	cpu.LoadValueIntoRegister(0x42, internal.ACCUMULATOR)
	cpu.LoadValueIntoRegister(0x33, internal.REGISTER_X)
	cpu.LoadValueIntoRegister(0x55, internal.REGISTER_Y)

	// Get debug data
	debugData := cpu.GetDebugData()

	// Verify all values are present and correct
	if len(debugData) != 3 {
		t.Errorf("Expected debug data to have 3 entries, got %d", len(debugData))
	}

	if debugData["acc"] != 0x42 {
		t.Errorf("Expected acc to be 0x42, got %d", debugData["acc"])
	}
	if debugData["x"] != 0x33 {
		t.Errorf("Expected x to be 0x33, got %d", debugData["x"])
	}
	if debugData["y"] != 0x55 {
		t.Errorf("Expected y to be 0x55, got %d", debugData["y"])
	}
}

func TestGetDebugDataIsIndependent(t *testing.T) {
	b := bus.NewBus()
	memory := memory.NewMemory(b)
	cpu := internal.NewCpu(memory, b)

	// Get initial debug data
	debugData1 := cpu.GetDebugData()
	debugData1["acc"] = 0xFF // Modify the returned map

	// Get debug data again
	debugData2 := cpu.GetDebugData()

	// Original register should not be affected by map modification
	if debugData2["acc"] != 0 {
		t.Errorf("Expected acc to still be 0, got %d", debugData2["acc"])
	}
}

func TestPushValueToStackDecrementsStackPointer(t *testing.T) {
	b := bus.NewBus()
	memory := memory.NewMemory(b)
	cpu := internal.NewCpu(memory, b)

	cpu.PushValueToStack(0x42)

	if cpu.GetStackPointer() != 0xFE {
		t.Errorf("Expected stack pointer to be 0xFE, got %d", cpu.GetStackPointer())
	}

	if memory.Read(0x1FF) != 0x42 {
		t.Errorf("Expected value 0x42 to be pushed to stack, got %d", memory.Read(0x1FF))
	}
}

func TestPushValueToStackMultipleTimesDecrementsStackPointerCorrectly(t *testing.T) {
	b := bus.NewBus()
	memory := memory.NewMemory(b)
	cpu := internal.NewCpu(memory, b)

	cpu.PushValueToStack(0x42)
	cpu.PushValueToStack(0x43)
	cpu.PushValueToStack(0x44)

	if cpu.GetStackPointer() != 0xFC {
		t.Errorf("Expected stack pointer to be 0xFC, got %d", cpu.GetStackPointer())
	}

	if memory.Read(0x1FF) != 0x42 {
		t.Errorf("Expected value 0x42 to be pushed to stack, got %d", memory.Read(0x1FF))
	}

	if memory.Read(0x1FE) != 0x43 {
		t.Errorf("Expected value 0x43 to be pushed to stack, got %d", memory.Read(0x1FE))
	}

	if memory.Read(0x1FD) != 0x44 {
		t.Errorf("Expected value 0x44 to be pushed to stack, got %d", memory.Read(0x1FD))
	}
}

func TestPullValueFromStackIncrementsStackPointer(t *testing.T) {
	b := bus.NewBus()
	memory := memory.NewMemory(b)
	cpu := internal.NewCpu(memory, b)

	cpu.PushValueToStack(0x42)
	cpu.PushValueToStack(0x43)
	cpu.PushValueToStack(0x44)

	value := cpu.PullValueFromStack()

	if value != 0x44 {
		t.Errorf("Expected value 0x44 to be pulled from stack, got %d", value)
	}

	if cpu.GetStackPointer() != 0xFD {
		t.Errorf("Expected stack pointer to be 0xFD, got %d", cpu.GetStackPointer())
	}

	if memory.Read(0x1FD) != 0x0 {
		t.Errorf("Expected value 0x0 to be pulled from stack, got %d", memory.Read(0x1FD))
	}
}

func TestRunProgram(t *testing.T) {
	bus := bus.NewBus()
	memory := memory.NewMemory(bus)
	sut := cpu.NewCpuWithStopAt(memory, bus, 0xC002)

	memory.Write(0xC000, 0xA9) //LDA #
	memory.Write(0xC001, 34)   //34

	//Pointer to start on 0xC000
	memory.Write(0xFFFC, 0x00)
	memory.Write(0xFFFC, 0xC0)

	sut.RunProgram()

	acc := sut.GetDebugData()["acc"]
	pc := sut.GetProgramCounter()

	if acc != 34 {
		t.Error("Accumulator should have value")
		return
	}

	if pc != 0xC002 {
		t.Error("Program Counter should have value")
		return
	}
}

func TestRunProgramWithoutCartidge(t *testing.T) {
	bus := bus.NewBus()
	memory := memory.NewMemory(bus)
	sut := cpu.NewCpuWithStopAt(memory, bus, 0)

	memory.Write(0xC000, 0xA9) //LDA #
	memory.Write(0xC001, 34)   //34
	memory.Write(0xFFFC, 0)    //Start

	sut.RunProgram()

	acc := sut.GetDebugData()["acc"]
	pc := sut.GetProgramCounter()

	if acc != 0 {
		t.Error("Accumulator should not have value")
		return
	}

	if pc != 0 {
		t.Error("Program Counter should not have value")
		return
	}
}
