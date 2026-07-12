package tests

import (
	"nes-emu/src/emulator/components/bus"
	"nes-emu/src/emulator/components/memory"
	"nes-emu/src/emulator/components/ppu"
	"nes-emu/test/fixtures"
	"testing"
)

func TestPPURender_ShouldDrawAPixel(t *testing.T) {
	mem := memory.NewMemory()
	vMemory := memory.NewMemory()
	bus := bus.NewBus(mem)
	ppu := ppu.NewPPU(bus, vMemory)

	ppu.Render()

	if ppu.GetCurrentScanline() != 0 {
		t.Fatal("scanline expected to be 0")
	}

	if ppu.GetCurrentScanlinePixel() != 1 {
		t.Fatal("expected to draw a pixel")
	}
}

func TestPPURender_ShouldWrapScanlineToStart(t *testing.T) {
	mem := memory.NewMemory()
	vMemory := memory.NewMemory()
	bus := bus.NewBus(mem)
	ppu := ppu.NewPPU(bus, vMemory)
	i := 0

	for {
		ppu.Render()
		i++

		if i == 256*262 {
			break
		}
	}

	if ppu.GetCurrentScanline() != 0 {
		t.Errorf("scanline expected to be 0, got %d", ppu.GetCurrentScanline())
	}

	if ppu.GetCurrentScanlinePixel() != 0 {
		t.Errorf("expected to be at pixel 0, got %d", ppu.GetCurrentScanlinePixel())
	}
}

func TestPPURender_ShouldNotTriggerAnNMIOnVBlank_WhenNMIFlagDisabled(t *testing.T) {
	mem := memory.NewMemory()
	vMemory := memory.NewMemory()
	bus := bus.NewBus(mem)
	mockCpu := &fixtures.MockCPU{}

	bus.AttachNMI(mockCpu)

	ppu := ppu.NewPPU(bus, vMemory)
	i := 0

	for {
		ppu.Render()
		i++

		if i == 256*240 {
			break
		}
	}

	if ppu.GetCurrentScanline() != 240 {
		t.Errorf("scanline expected to be 0, got %d", ppu.GetCurrentScanline())
	}

	if ppu.GetCurrentScanlinePixel() != 0 {
		t.Errorf("expected to be at pixel 0, got %d", ppu.GetCurrentScanlinePixel())
	}

	if mockCpu.NmiCalled > 0 {
		t.Errorf("did not expect NMI to be called")
	}
}

func TestPPURender_ShouldTriggerAnNMIOnVBlank_WhenNMIFlagEnabled(t *testing.T) {
	mem := memory.NewMemory()
	vMemory := memory.NewMemory()
	bus := bus.NewBus(mem)
	mockCpu := &fixtures.MockCPU{}

	bus.AttachNMI(mockCpu)

	ppu := ppu.NewPPU(bus, vMemory)
	i := 0

	mem.Write(0x2000, 0b10000000)

	for {
		ppu.Render()
		i++

		if i == 256*240 {
			break
		}
	}

	if ppu.GetCurrentScanline() != 240 {
		t.Errorf("scanline expected to be 0, got %d", ppu.GetCurrentScanline())
	}

	if ppu.GetCurrentScanlinePixel() != 0 {
		t.Errorf("expected to be at pixel 0, got %d", ppu.GetCurrentScanlinePixel())
	}

	if mockCpu.NmiCalled == 0 {
		t.Errorf("did expect NMI to be called")
	}
}
