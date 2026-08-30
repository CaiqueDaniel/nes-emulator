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
	bus := bus.NewBusWithWorkMemory(mem)
	mockPipeline := &PixelPipelineFixture{}
	ppu := ppu.NewPPU(bus, vMemory, mockPipeline)

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
	bus := bus.NewBusWithWorkMemory(mem)
	mockPipeline := &PixelPipelineFixture{}
	ppu := ppu.NewPPU(bus, vMemory, mockPipeline)

	for i := 0; i < 336*262; i++ {
		ppu.Render()
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
	bus := bus.NewBusWithWorkMemory(mem)
	mockCpu := &fixtures.MockCPU{}
	mockPipeline := &PixelPipelineFixture{}

	bus.AttachNMI(mockCpu)

	ppu := ppu.NewPPU(bus, vMemory, mockPipeline)

	for i := 0; i < 336*240; i++ {
		ppu.Render()
	}

	if ppu.GetCurrentScanline() != 240 {
		t.Errorf("scanline expected to be 240, got %d", ppu.GetCurrentScanline())
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
	bus := bus.NewBusWithWorkMemory(mem)
	mockCpu := &fixtures.MockCPU{}
	mockPipeline := &PixelPipelineFixture{}

	bus.AttachNMI(mockCpu)

	ppu := ppu.NewPPU(bus, vMemory, mockPipeline)

	mem.Write(0x2000, 0b10000000)

	for i := 0; i < 336*240; i++ {
		ppu.Render()
	}

	if ppu.GetCurrentScanline() != 240 {
		t.Errorf("scanline expected to be 240, got %d", ppu.GetCurrentScanline())
	}

	if ppu.GetCurrentScanlinePixel() != 0 {
		t.Errorf("expected to be at pixel 0, got %d", ppu.GetCurrentScanlinePixel())
	}

	if mockCpu.NmiCalled == 0 {
		t.Errorf("did expect NMI to be called")
	}
}

func TestPPURender_ShouldResetFlagsOnStatusRegister_OnPreRender(t *testing.T) {
	mem := memory.NewMemory()
	vMemory := memory.NewMemory()
	bus := bus.NewBusWithWorkMemory(mem)
	mockCpu := &fixtures.MockCPU{}
	mockPipeline := &PixelPipelineFixture{}

	bus.AttachNMI(mockCpu)

	ppu := ppu.NewPPU(bus, vMemory, mockPipeline)

	for i := 0; i < 336*261; i++ {
		ppu.Render()
	}

	if ppu.GetCurrentScanline() != 261 {
		t.Errorf("scanline expected to be 261, got %d", ppu.GetCurrentScanline())
	}

	if ppu.GetCurrentScanlinePixel() != 0 {
		t.Errorf("expected to be at pixel 0, got %d", ppu.GetCurrentScanlinePixel())
	}

	if mem.Read(0x2002) != 0 {
		t.Errorf("did expect flags to be cleared")
	}
}

func TestPPURender_ShouldResetFlagsOnStatusRegister_WithoutChangingOtherBits_OnPreRender(t *testing.T) {
	mem := memory.NewMemory()
	vMemory := memory.NewMemory()
	bus := bus.NewBusWithWorkMemory(mem)
	mockCpu := &fixtures.MockCPU{}

	bus.AttachNMI(mockCpu)

	mockPipeline := &PixelPipelineFixture{}
	ppu := ppu.NewPPU(bus, vMemory, mockPipeline)

	mem.Write(0x2002, 0b1111_1111)

	for i := 0; i < 336*261; i++ {
		ppu.Render()
	}

	if ppu.GetCurrentScanline() != 261 {
		t.Errorf("scanline expected to be 0, got %d", ppu.GetCurrentScanline())
	}

	if ppu.GetCurrentScanlinePixel() != 0 {
		t.Errorf("expected to be at pixel 0, got %d", ppu.GetCurrentScanlinePixel())
	}

	if mem.Read(0x2002) != 0b0001_1111 {
		t.Errorf("did expect value to be 0b0001_1111. got %b", mem.Read(0x2002))
	}
}

func TestPPURender_ShouldSetVBlankFlagOnStatusRegister_OnVBlank(t *testing.T) {
	mem := memory.NewMemory()
	vMemory := memory.NewMemory()
	bus := bus.NewBusWithWorkMemory(mem)
	mockCpu := &fixtures.MockCPU{}
	mockPipeline := &PixelPipelineFixture{}

	bus.AttachNMI(mockCpu)

	ppu := ppu.NewPPU(bus, vMemory, mockPipeline)

	for i := 0; i < 336*240; i++ {
		ppu.Render()
	}

	if ppu.GetCurrentScanline() != 240 {
		t.Errorf("scanline expected to be 0, got %d", ppu.GetCurrentScanline())
	}

	if ppu.GetCurrentScanlinePixel() != 0 {
		t.Errorf("expected to be at pixel 0, got %d", ppu.GetCurrentScanlinePixel())
	}

	if mem.Read(0x2002) != 0b1000_0000 {
		t.Errorf("did expect v-blank to be setted")
	}
}

func TestPPURender_ShouldShiftRegisters_OnVisibleScanlines(t *testing.T) {
	mem := memory.NewMemory()
	vMemory := memory.NewMemory()
	bus := bus.NewBusWithWorkMemory(mem)
	mockPipeline := &PixelPipelineFixture{}
	sut := ppu.NewPPU(bus, vMemory, mockPipeline)

	for i := 0; i <= 16; i++ {
		sut.Render()
	}

	lowPatternShiftRegister, highPatternShiftRegister, lowAttributeShiftRegister, highAttributeShiftRegister := sut.GetShiftRegisters()

	if lowPatternShiftRegister != 0b10_00000010 {
		t.Errorf("expected low pattern shift register to be 0b00000010_00000010, got %b", lowPatternShiftRegister)
	}

	if highPatternShiftRegister != 0b11_00000011 {
		t.Errorf("expected high pattern shift register to be 0b00000011_00000011, got %b", highPatternShiftRegister)
	}

	if lowAttributeShiftRegister != 0b1_00000001 {
		t.Errorf("expected low attribute shift register to be 0b00000001_00000001, got %b", lowAttributeShiftRegister)
	}

	if highAttributeShiftRegister != 0b1_00000001 {
		t.Errorf("expected high attribute shift register to be 0b00000001_00000001, got %b", highAttributeShiftRegister)
	}
}

func TestPPURender_ShouldShiftRegisters_OnHBlank(t *testing.T) {
	mem := memory.NewMemory()
	vMemory := memory.NewMemory()
	bus := bus.NewBusWithWorkMemory(mem)
	mockPipeline := &PixelPipelineFixture{}
	sut := ppu.NewPPU(bus, vMemory, mockPipeline)

	for i := 0; i <= 256; i++ {
		sut.Render()
	}

	lowPatternShiftRegister, highPatternShiftRegister, lowAttributeShiftRegister, highAttributeShiftRegister := sut.GetShiftRegisters()

	if lowPatternShiftRegister != 0b1000_00010 {
		t.Errorf("expected low pattern shift register to be 0b1000_00010, got %b", lowPatternShiftRegister)
	}

	if highPatternShiftRegister != 0b1000000110000011 {
		t.Errorf("expected high pattern shift register to be 0b1000000110000011, got %b", highPatternShiftRegister)
	}

	if lowAttributeShiftRegister != 0b1000000010000001 {
		t.Errorf("expected low attribute shift register to be 0b1000000010000001, got %b", lowAttributeShiftRegister)
	}

	if highAttributeShiftRegister != 0b1000000010000001 {
		t.Errorf("expected high attribute shift register to be 0b1000000010000001, got %b", highAttributeShiftRegister)
	}

	sut.Render()
	lowPatternShiftRegister, highPatternShiftRegister, lowAttributeShiftRegister, highAttributeShiftRegister = sut.GetShiftRegisters()

	if lowPatternShiftRegister != 0b1000_00010 {
		t.Errorf("expected low pattern shift register to be 0b1000_00010, got %b", lowPatternShiftRegister)
	}

	if highPatternShiftRegister != 0b1000000110000011 {
		t.Errorf("expected high pattern shift register to be 0b1000000110000011, got %b", highPatternShiftRegister)
	}

	if lowAttributeShiftRegister != 0b1000000010000001 {
		t.Errorf("expected low attribute shift register to be 0b1000000010000001, got %b", lowAttributeShiftRegister)
	}

	if highAttributeShiftRegister != 0b1000000010000001 {
		t.Errorf("expected high attribute shift register to be 0b1000000010000001, got %b", highAttributeShiftRegister)
	}
}

func TestPPURender_ShouldShiftRegisters_OnVBlank(t *testing.T) {
	mem := memory.NewMemory()
	vMemory := memory.NewMemory()
	bus := bus.NewBusWithWorkMemory(mem)
	mockPipeline := &PixelPipelineFixture{}
	sut := ppu.NewPPU(bus, vMemory, mockPipeline)

	for i := 0; i <= 336*240; i++ {
		sut.Render()
	}

	lowPatternShiftRegister, highPatternShiftRegister, lowAttributeShiftRegister, highAttributeShiftRegister := sut.GetShiftRegisters()

	if lowPatternShiftRegister != 0b100000010 {
		t.Errorf("expected low pattern shift register to be 0b100000010, got %b", lowPatternShiftRegister)
	}

	if highPatternShiftRegister != 0b1000000110000011 {
		t.Errorf("expected high pattern shift register to be 0b1000000110000011, got %b", highPatternShiftRegister)
	}

	if lowAttributeShiftRegister != 0b1000000010000001 {
		t.Errorf("expected low attribute shift register to be 0b1000000010000001, got %b", lowAttributeShiftRegister)
	}

	if highAttributeShiftRegister != 0b1000000010000001 {
		t.Errorf("expected high attribute shift register to be 0b1000000010000001, got %b", highAttributeShiftRegister)
	}

	sut.Render()
	lowPatternShiftRegister, highPatternShiftRegister, lowAttributeShiftRegister, highAttributeShiftRegister = sut.GetShiftRegisters()

	if lowPatternShiftRegister != 0b100000010 {
		t.Errorf("expected low pattern shift register to be 0b100000010, got %b", lowPatternShiftRegister)
	}

	if highPatternShiftRegister != 0b1000000110000011 {
		t.Errorf("expected high pattern shift register to be 0b1000000110000011, got %b", highPatternShiftRegister)
	}

	if lowAttributeShiftRegister != 0b1000000010000001 {
		t.Errorf("expected low attribute shift register to be 0b1000000010000001, got %b", lowAttributeShiftRegister)
	}

	if highAttributeShiftRegister != 0b1000000010000001 {
		t.Errorf("expected high attribute shift register to be 0b1000000010000001, got %b", highAttributeShiftRegister)
	}
}
