package tests

import (
	"nes-emu/src/emulator/components/bus"
	internal "nes-emu/src/emulator/components/cpu"
	"nes-emu/src/emulator/components/memory"
	"testing"
)

func TestInstructionCycles(t *testing.T) {
	// Map of opcode to expected base cycles.
	// Omitted branching and jumping instructions to avoid breaking the Test RunProgram execution loop.
	expectedCycles := map[uint8]uint{
		// Flags
		0x18: 2, 0x38: 2, 0x58: 2, 0x78: 2, 0xD8: 2, 0xF8: 2, 0xB8: 2, 0xEA: 2,

		// Load
		0xA9: 2, 0xA5: 3, 0xB5: 4, 0xAD: 4, 0xBD: 4, 0xB9: 4, 0xA1: 6, 0xB1: 5,
		0xA2: 2, 0xA6: 3, 0xB6: 4, 0xAE: 4, 0xBE: 4,
		0xA0: 2, 0xA4: 3, 0xB4: 4, 0xAC: 4, 0xBC: 4,

		// Store
		0x85: 3, 0x95: 4, 0x8D: 4, 0x9D: 5, 0x99: 5, 0x81: 6, 0x91: 6,
		0x86: 3, 0x96: 4, 0x8E: 4,
		0x84: 3, 0x94: 4, 0x8C: 4,

		// Transfers
		0xAA: 2, 0xA8: 2, 0xBA: 2, 0x8A: 2, 0x98: 2, 0x9A: 2,

		// Arithmetic (ADC, SBC)
		0x69: 2, 0x65: 3, 0x75: 4, 0x6D: 4, 0x7D: 4, 0x79: 4, 0x61: 6, 0x71: 5,
		0xE9: 2, 0xE5: 3, 0xF5: 4, 0xED: 4, 0xFD: 4, 0xF9: 4, 0xE1: 6, 0xF1: 5,

		// INC/DEC
		0xE6: 5, 0xF6: 6, 0xEE: 6, 0xFE: 7,
		0xC6: 5, 0xD6: 6, 0xCE: 6, 0xDE: 7,
		0xE8: 2, 0xC8: 2, 0xCA: 2, 0x88: 2,

		// Shift
		0x0A: 2, 0x06: 5, 0x16: 6, 0x0E: 6, 0x1E: 7,
		0x4A: 2, 0x46: 5, 0x56: 6, 0x4E: 6, 0x5E: 7,
		0x2A: 2, 0x26: 5, 0x36: 6, 0x2E: 6, 0x3E: 7,
		0x6A: 2, 0x66: 5, 0x76: 6, 0x6E: 6, 0x7E: 7,

		// Bitwise
		0x29: 2, 0x25: 3, 0x35: 4, 0x2D: 4, 0x3D: 4, 0x39: 4, 0x21: 6, 0x31: 5,
		0x09: 2, 0x05: 3, 0x15: 4, 0x0D: 4, 0x1D: 4, 0x19: 4, 0x01: 6, 0x11: 5,
		0x49: 2, 0x45: 3, 0x55: 4, 0x4D: 4, 0x5D: 4, 0x59: 4, 0x41: 6, 0x51: 5,
		0x24: 3, 0x2C: 4,

		// Compare
		0xC9: 2, 0xC5: 3, 0xD5: 4, 0xCD: 4, 0xDD: 4, 0xD9: 4, 0xC1: 6, 0xD1: 5,
		0xE0: 2, 0xE4: 3, 0xEC: 4,
		0xC0: 2, 0xC4: 3, 0xCC: 4,

		// Stack
		0x48: 3, 0x68: 4, 0x08: 3, 0x28: 4,
	}

	for opcode, expected := range expectedCycles {
		mem := memory.NewMemory()
		b := bus.NewBus(mem)

		// Set the Reset vector to 0x8000
		mem.Write(0xFFFC, 0x00)
		mem.Write(0xFFFD, 0x80)

		// Write the instruction to be tested
		mem.Write(0x8000, opcode)

		// Write dummy arguments to avoid reading uninitialized memory
		mem.Write(0x8001, 0x00)
		mem.Write(0x8002, 0x00)

		// Stop right after the instruction is fetched (PC increments)
		// RunProgram stops when programCounter >= stopPcAt
		cpu := internal.NewCpuWithStopAt(b, 0x8001)

		b.ResetTickCount()
		cpu.RunProgram()

		// RunProgram executes the following memory reads before running the opcode:
		// 1. Read 0xFFFC (1 cycle)
		// 2. Read 0xFFFD (1 cycle)
		// Thus, we subtract 2 from the total ticks.
		actualCycles := b.GetTickCount() - 2

		if actualCycles != expected {
			t.Errorf("Instruction 0x%02X: expected %d cycles, got %d", opcode, expected, actualCycles)
		}
	}
}

func TestInstructionCyclesPageCross(t *testing.T) {
	// Instructions that get +1 cycle on page cross
	absXOps := []uint8{0xBD, 0x7D, 0xFD, 0xDD, 0x3D, 0x1D, 0x5D, 0xBC}
	absYOps := []uint8{0xB9, 0xBE, 0x79, 0xF9, 0xD9, 0x39, 0x19, 0x59}
	indYOps := []uint8{0xB1, 0x71, 0xF1, 0xD1, 0x31, 0x11, 0x51}

	type pageCrossTest struct {
		opcode         uint8
		name           string
		register       string
		registerValue  uint8
		args           []uint8
		expectedCycles uint
	}

	var tests []pageCrossTest

	for _, op := range absXOps {
		tests = append(tests, pageCrossTest{
			opcode:         op,
			name:           "AbsX Page Cross",
			register:       internal.REGISTER_X,
			registerValue:  0x01,
			args:           []uint8{0xFF, 0x01}, // Address 0x01FF
			expectedCycles: 5,                   // 4 + 1
		})
	}

	for _, op := range absYOps {
		tests = append(tests, pageCrossTest{
			opcode:         op,
			name:           "AbsY Page Cross",
			register:       internal.REGISTER_Y,
			registerValue:  0x01,
			args:           []uint8{0xFF, 0x01}, // Address 0x01FF
			expectedCycles: 5,                   // 4 + 1
		})
	}

	for _, op := range indYOps {
		tests = append(tests, pageCrossTest{
			opcode:         op,
			name:           "IndY Page Cross",
			register:       internal.REGISTER_Y,
			registerValue:  0x01,
			args:           []uint8{0x80}, // Zero page pointer
			expectedCycles: 6,             // 5 + 1
		})
	}

	for _, tt := range tests {
		mem := memory.NewMemory()
		b := bus.NewBus(mem)

		// Set the Reset vector to 0x8000
		mem.Write(0xFFFC, 0x00)
		mem.Write(0xFFFD, 0x80)

		if tt.name == "IndY Page Cross" {
			mem.Write(0x80, 0xFF)
			mem.Write(0x81, 0x01)
		}

		pc := uint16(0x8000)

		mem.Write(pc, tt.opcode)
		pc++

		for _, arg := range tt.args {
			mem.Write(pc, arg)
			pc++
		}

		cpu := internal.NewCpuWithStopAt(b, int(pc))
		cpu.LoadValueIntoRegister(tt.registerValue, tt.register)

		b.ResetTickCount()
		cpu.RunProgram()

		// Subtract Reset Vector read cycles (2 cycles)
		actualCycles := b.GetTickCount() - 2

		if actualCycles != tt.expectedCycles {
			t.Errorf("Instruction 0x%02X (%s): expected %d cycles, got %d", tt.opcode, tt.name, tt.expectedCycles, actualCycles)
		}
	}
}

func TestBreakInstructionCycles(t *testing.T) {
	mem := memory.NewMemory()
	b := bus.NewBus(mem)

	// Set Reset vector to 0x8000
	mem.Write(0xFFFC, 0x00)
	mem.Write(0xFFFD, 0x80)

	// BRK instruction at 0x8000
	mem.Write(0x8000, 0x00)

	// Dummy byte for BRK padding
	mem.Write(0x8001, 0x00)

	// Set Interrupt vector to 0x8005
	mem.Write(0xFFFE, 0x05)
	mem.Write(0xFFFF, 0x80)

	// Stop at 0x8005 (where the interrupt vector points)
	sut := internal.NewCpuWithStopAt(b, 0x8005)

	b.ResetTickCount()
	sut.RunProgram()

	// RunProgram reads the reset vector (2 cycles)
	// BRK should take 7 cycles
	actualCycles := b.GetTickCount() - 2

	if actualCycles != 7 {
		t.Errorf("Expected 7 cycles for BRK instruction, got %d", actualCycles)
	}
}
