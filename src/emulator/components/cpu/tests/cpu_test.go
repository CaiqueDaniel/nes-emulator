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
	memory := memory.NewMemory()
	b := bus.NewBusWithWorkMemory(memory)
	cpu := internal.NewCpuWithInternal(b)

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
	if cpu.GetCarryFlag() {
		t.Errorf("Expected carry flag to be disabled")
	}
	if cpu.GetDecimalFlag() {
		t.Errorf("Expected decimal flag to be disabled")
	}
	if cpu.GetNegativeFlag() {
		t.Errorf("Expected negative flag to be disabled")
	}
	if cpu.GetOverflowFlag() {
		t.Errorf("Expected overflow flag to be disabled")
	}
	if cpu.GetZeroFlag() {
		t.Errorf("Expected zero flag to be disabled")
	}
	if !cpu.GetIRQFlag() {
		t.Errorf("Expected interrupt flag to be enabled")
	}
}

func TestCpuResetClearsAllRegisters(t *testing.T) {
	memory := memory.NewMemory()
	b := bus.NewBusWithWorkMemory(memory)
	cpu := internal.NewCpuWithInternal(b)

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
	memory := memory.NewMemory()
	b := bus.NewBusWithWorkMemory(memory)
	cpu := internal.NewCpuWithInternal(b)

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
	memory := memory.NewMemory()
	b := bus.NewBusWithWorkMemory(memory)
	cpu := internal.NewCpuWithInternal(b)

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
	memory := memory.NewMemory()
	b := bus.NewBusWithWorkMemory(memory)
	cpu := internal.NewCpuWithInternal(b)

	cpu.PushValueToStack(0x42)

	if cpu.GetStackPointer() != 0xFE {
		t.Errorf("Expected stack pointer to be 0xFE, got %d", cpu.GetStackPointer())
	}

	if memory.Read(0x1FF) != 0x42 {
		t.Errorf("Expected value 0x42 to be pushed to stack, got %d", memory.Read(0x1FF))
	}
}

func TestPushValueToStackMultipleTimesDecrementsStackPointerCorrectly(t *testing.T) {
	memory := memory.NewMemory()
	b := bus.NewBusWithWorkMemory(memory)
	cpu := internal.NewCpuWithInternal(b)

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
	memory := memory.NewMemory()
	b := bus.NewBusWithWorkMemory(memory)
	cpu := internal.NewCpuWithInternal(b)

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
	memory := memory.NewMemory()
	bus := bus.NewBusWithWorkMemory(memory)
	sut := cpu.NewCpuWithStopAt(bus, 0xC002)

	memory.Write(0xC000, 0xA9) //LDA #
	memory.Write(0xC001, 34)   //34

	//Pointer to start on 0xC000
	memory.Write(0xFFFC, 0x00)
	memory.Write(0xFFFD, 0xC0)

	bus.ResetTickCount()
	sut.RunProgram()

	acc := sut.GetDebugData()["acc"]
	pc := sut.GetProgramCounter()
	clockTicks := bus.GetTickCount()

	if acc != 34 {
		t.Error("Accumulator should have value")
		return
	}

	if pc != 0xC002 {
		t.Error("Program Counter should have value")
		return
	}

	if clockTicks != 4 {
		t.Errorf("Bus should have 4 ticks. It got: %d", clockTicks)
		return
	}
}

func TestNumberOfInstructions(t *testing.T) {
	memory := memory.NewMemory()
	bus := bus.NewBusWithWorkMemory(memory)
	cpu := cpu.NewCpuWithInternal(bus)

	if cpu.GetNumberOfInstructions() != 151 {
		t.Errorf("Incorrect number of instructions loaded (%d)", cpu.GetNumberOfInstructions())
	}
}

func TestSetNMI(t *testing.T) {
	memory := memory.NewMemory()
	b := bus.NewBusWithWorkMemory(memory)
	cpu := internal.NewCpuWithInternal(b)

	if cpu.GetNMIFlag() {
		t.Errorf("Expected NMI flag to be false initially")
	}

	cpu.SetNMI()

	if !cpu.GetNMIFlag() {
		t.Errorf("Expected NMI flag to be true after calling SetNMI")
	}
}

func TestHandleNMI(t *testing.T) {
	memory := memory.NewMemory()
	b := bus.NewBusWithWorkMemory(memory)
	cpu := internal.NewCpuWithProgramCounter(0x1234, b)

	// Configure NMI vector in memory (0xFFFA and 0xFFFB)
	memory.Write(0xFFFA, 0x00)
	memory.Write(0xFFFB, 0x80) // Vector address: 0x8000

	// Handle NMI
	cpu.HandleNMI()

	// Verify Program Counter is set to the vector address
	if cpu.GetProgramCounter() != 0x8000 {
		t.Errorf("Expected Program Counter to be 0x8000, got 0x%X", cpu.GetProgramCounter())
	}

	// Verify stack pointer decreased by 1 (PushStatusIntoStack)
	if cpu.GetStackPointer() != 0xFE {
		t.Errorf("Expected Stack Pointer to be 0xFE, got 0x%X", cpu.GetStackPointer())
	}

	// Verify the stack contents for pushed status
	// Default status flag is 0x30
	if memory.Read(0x01FF) != 0x30 {
		t.Errorf("Expected Stack at 0x01FF to have Status 0x30, got 0x%X", memory.Read(0x01FF))
	}
}

func TestRunProgramHandlesNMI(t *testing.T) {
	memory := memory.NewMemory()
	b := bus.NewBusWithWorkMemory(memory)
	cpu := internal.NewCpuWithStopAt(b, 0x8001)

	// Set Start Pointer (0xFFFC) to 0xC000
	memory.Write(0xFFFC, 0x00)
	memory.Write(0xFFFD, 0xC0)

	// Normal program instruction at 0xC000
	memory.Write(0xC000, 0xEA) // NOP

	// Configure NMI vector in memory (0xFFFA and 0xFFFB)
	memory.Write(0xFFFA, 0x00)
	memory.Write(0xFFFB, 0x80) // Vector address: 0x8000

	// NMI handler instruction at 0x8000
	memory.Write(0x8000, 0xEA) // NOP

	// Trigger NMI before running the program
	cpu.SetNMI()

	// Run program
	cpu.RunProgram()

	// Validate NMI flag is cleared
	if cpu.GetNMIFlag() {
		t.Errorf("Expected NMI flag to be cleared after handling")
	}

	// Validate it executed the NMI handler instead of normal flow
	if cpu.GetProgramCounter() != 0x8002 {
		t.Errorf("Expected Program Counter to be 0x8002 (0x8001 inside NMI handler + 1), got 0x%X", cpu.GetProgramCounter())
	}
}
