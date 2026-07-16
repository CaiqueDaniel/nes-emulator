package tests

import (
	"nes-emu/src/emulator/components/bus"
	"nes-emu/src/emulator/components/cpu"
	"nes-emu/src/emulator/components/memory"
	"testing"
	"time"
)

const clock_in_mhz = 1.789773

func TestClockTime(t *testing.T) {
	ram := memory.NewMemory()
	board := bus.NewBus(ram)
	sut := cpu.NewCpuWithStopAt(board, 16)
	expectedCycles := uint(92)
	expectedTimeElapsed := time.Second / 60

	program := []uint8{
		0xA2, 0x00, //LDX #00
		0xE8,      //INX
		0xE0, 0xA, //CPX #FF
		0xF0, 0x0A, //BEQ #FF
		0x4C, 0x02, //JMP $02
	}

	for address, programByte := range program {
		ram.Write(uint16(address), programByte)
	}

	ram.Write(0xFFFC, 0x0)
	ram.Write(0xFFFD, 0x0)

	startTime := time.Now()

	sut.RunProgram()

	timeElapsed := time.Since(startTime)

	if board.GetTickCount() != expectedCycles {
		t.Errorf("expected 92 cycles. got %d", board.GetTickCount())
	}

	if timeElapsed < expectedTimeElapsed {
		t.Errorf("expected program to last at least 16.6666ms frame. got %dms", timeElapsed)
	}
}
