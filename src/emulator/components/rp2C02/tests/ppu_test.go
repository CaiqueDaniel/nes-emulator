package tests

import (
	"nes-emu/src/emulator/components/bus"
	"nes-emu/src/emulator/components/memory"
	ppu "nes-emu/src/emulator/components/rp2C02"
	"nes-emu/test/fixtures"
	"testing"
)

func TestPPURender_ShouldDrawAPixel(t *testing.T) {
	mem := memory.NewMemory()
	vMemory := memory.NewMemory()
	bus := bus.NewBusWithWorkMemory(mem)
	bus.AtatchVideoMemory(vMemory)
	mockPipeline := &PixelPipelineFixture{}
	screenFixture := NewScreenFixture()
	ppu := ppu.NewRp2C02(bus, mockPipeline, screenFixture)

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
	bus.AtatchVideoMemory(vMemory)
	mockPipeline := &PixelPipelineFixture{}
	screenFixture := NewScreenFixture()
	ppu := ppu.NewRp2C02(bus, mockPipeline, screenFixture)

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
	bus.AtatchVideoMemory(vMemory)
	mockCpu := &fixtures.MockCPU{}
	mockPipeline := &PixelPipelineFixture{}

	bus.AttachNMI(mockCpu)

	screenFixture := NewScreenFixture()
	ppu := ppu.NewRp2C02(bus, mockPipeline, screenFixture)

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
	bus.AtatchVideoMemory(vMemory)
	mockCpu := &fixtures.MockCPU{}
	mockPipeline := &PixelPipelineFixture{}

	bus.AttachNMI(mockCpu)

	screenFixture := NewScreenFixture()
	ppu := ppu.NewRp2C02(bus, mockPipeline, screenFixture)

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
	bus.AtatchVideoMemory(vMemory)
	mockCpu := &fixtures.MockCPU{}
	mockPipeline := &PixelPipelineFixture{}

	bus.AttachNMI(mockCpu)

	screenFixture := NewScreenFixture()
	ppu := ppu.NewRp2C02(bus, mockPipeline, screenFixture)

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
	bus.AtatchVideoMemory(vMemory)
	mockCpu := &fixtures.MockCPU{}

	bus.AttachNMI(mockCpu)

	mockPipeline := &PixelPipelineFixture{}
	screenFixture := NewScreenFixture()
	ppu := ppu.NewRp2C02(bus, mockPipeline, screenFixture)

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
	bus.AtatchVideoMemory(vMemory)
	mockCpu := &fixtures.MockCPU{}
	mockPipeline := &PixelPipelineFixture{}

	bus.AttachNMI(mockCpu)

	screenFixture := NewScreenFixture()
	ppu := ppu.NewRp2C02(bus, mockPipeline, screenFixture)

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
	bus.AtatchVideoMemory(vMemory)
	mockPipeline := &PixelPipelineFixture{}
	screenFixture := NewScreenFixture()
	sut := ppu.NewRp2C02(bus, mockPipeline, screenFixture)

	for i := 0; i <= 16; i++ {
		sut.Render()
	}

	lowPatternShiftRegister, highPatternShiftRegister, lowAttributeShiftRegister, highAttributeShiftRegister := mockPipeline.GetShiftRegisters()

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
	bus.AtatchVideoMemory(vMemory)
	mockPipeline := &PixelPipelineFixture{}
	screenFixture := NewScreenFixture()
	sut := ppu.NewRp2C02(bus, mockPipeline, screenFixture)

	for i := 0; i <= 256; i++ {
		sut.Render()
	}

	lowPatternShiftRegister, highPatternShiftRegister, lowAttributeShiftRegister, highAttributeShiftRegister := mockPipeline.GetShiftRegisters()

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
	lowPatternShiftRegister, highPatternShiftRegister, lowAttributeShiftRegister, highAttributeShiftRegister = mockPipeline.GetShiftRegisters()

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
	bus.AtatchVideoMemory(vMemory)
	mockPipeline := &PixelPipelineFixture{}
	screen := NewScreenFixture()
	sut := ppu.NewRp2C02(bus, mockPipeline, screen)

	for i := 0; i <= 336*240; i++ {
		sut.Render()
	}

	lowPatternShiftRegister, highPatternShiftRegister, lowAttributeShiftRegister, highAttributeShiftRegister := mockPipeline.GetShiftRegisters()

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
	lowPatternShiftRegister, highPatternShiftRegister, lowAttributeShiftRegister, highAttributeShiftRegister = mockPipeline.GetShiftRegisters()

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

func TestPPUTriggerLatchWithWriteSignal_ShouldUpdateVRAMAddress_OnFirstAndSecondWriteToPPUAddress(t *testing.T) {
	mem := memory.NewMemory()
	vMemory := memory.NewMemory()
	bus := bus.NewBusWithWorkMemory(mem)
	bus.AtatchVideoMemory(vMemory)
	mockPipeline := &PixelPipelineFixture{}
	screenFixture := NewScreenFixture()
	ppu := ppu.NewRp2C02(bus, mockPipeline, screenFixture)

	// First write: High byte
	mem.Write(0x2006, 0x21)
	ppu.TriggerLatchWithWriteSignal(0x2006)

	// Second write: Low byte
	mem.Write(0x2006, 0x08)
	ppu.TriggerLatchWithWriteSignal(0x2006)

	// Write data to VRAM via 0x2007
	mem.Write(0x2007, 0x42)
	ppu.TriggerLatchWithWriteSignal(0x2007)

	if vMemory.Read(0x2108) != 0x42 {
		t.Errorf("expected video memory at 0x2108 to be 0x42, got 0x%X", vMemory.Read(0x2108))
	}
}

func TestPPUTriggerLatchWithWriteSignal_ShouldApplyMaskToHighByte_OnFirstWriteToPPUAddress(t *testing.T) {
	mem := memory.NewMemory()
	vMemory := memory.NewMemory()
	bus := bus.NewBusWithWorkMemory(mem)
	bus.AtatchVideoMemory(vMemory)
	mockPipeline := &PixelPipelineFixture{}
	screenFixture := NewScreenFixture()
	ppu := ppu.NewRp2C02(bus, mockPipeline, screenFixture)

	// First write: 0xFF should be masked with 0x3F00 -> high byte is 0x3F
	mem.Write(0x2006, 0xFF)
	ppu.TriggerLatchWithWriteSignal(0x2006)

	// Second write: 0x50
	mem.Write(0x2006, 0x50)
	ppu.TriggerLatchWithWriteSignal(0x2006)

	mem.Write(0x2007, 0x99)
	ppu.TriggerLatchWithWriteSignal(0x2007)

	if vMemory.Read(0x3F50) != 0x99 {
		t.Errorf("expected video memory at 0x3F50 to be 0x99, got 0x%X", vMemory.Read(0x3F50))
	}
}

func TestPPUTriggerLatchWithWriteSignal_ShouldToggleWriteLatch_OnConsecutiveWritesToPPUAddress(t *testing.T) {
	mem := memory.NewMemory()
	vMemory := memory.NewMemory()
	bus := bus.NewBusWithWorkMemory(mem)
	bus.AtatchVideoMemory(vMemory)
	mockPipeline := &PixelPipelineFixture{}
	screenFixture := NewScreenFixture()
	ppu := ppu.NewRp2C02(bus, mockPipeline, screenFixture)

	// 1st write (high byte)
	mem.Write(0x2006, 0x20)
	ppu.TriggerLatchWithWriteSignal(0x2006)
	// 2nd write (low byte) -> v becomes 0x2000
	mem.Write(0x2006, 0x00)
	ppu.TriggerLatchWithWriteSignal(0x2006)

	// 3rd write (high byte again because w toggled back to false)
	mem.Write(0x2006, 0x28)
	ppu.TriggerLatchWithWriteSignal(0x2006)
	// 4th write (low byte) -> v becomes 0x2815
	mem.Write(0x2006, 0x15)
	ppu.TriggerLatchWithWriteSignal(0x2006)

	mem.Write(0x2007, 0xAA)
	ppu.TriggerLatchWithWriteSignal(0x2007)

	if vMemory.Read(0x2815) != 0xAA {
		t.Errorf("expected video memory at 0x2815 to be 0xAA, got 0x%X", vMemory.Read(0x2815))
	}
}

func TestPPUTriggerLatchWithWriteSignal_ShouldWriteDataToVideoMemoryAndIncrementVBy1_WhenOffsetBitIsNotSet(t *testing.T) {
	mem := memory.NewMemory()
	vMemory := memory.NewMemory()
	bus := bus.NewBusWithWorkMemory(mem)
	bus.AtatchVideoMemory(vMemory)
	mockPipeline := &PixelPipelineFixture{}
	screenFixture := NewScreenFixture()
	ppu := ppu.NewRp2C02(bus, mockPipeline, screenFixture)

	// PPU_CONTROL bit 2 = 0 (increment by 1)
	mem.Write(0x2000, 0x00)

	// Set address to 0x2100
	mem.Write(0x2006, 0x21)
	ppu.TriggerLatchWithWriteSignal(0x2006)
	mem.Write(0x2006, 0x00)
	ppu.TriggerLatchWithWriteSignal(0x2006)

	mem.Write(0x2007, 0x11)
	ppu.TriggerLatchWithWriteSignal(0x2007)

	mem.Write(0x2007, 0x22)
	ppu.TriggerLatchWithWriteSignal(0x2007)

	mem.Write(0x2007, 0x33)
	ppu.TriggerLatchWithWriteSignal(0x2007)

	if vMemory.Read(0x2100) != 0x11 {
		t.Errorf("expected vMemory[0x2100] == 0x11, got 0x%X", vMemory.Read(0x2100))
	}
	if vMemory.Read(0x2101) != 0x22 {
		t.Errorf("expected vMemory[0x2101] == 0x22, got 0x%X", vMemory.Read(0x2101))
	}
	if vMemory.Read(0x2102) != 0x33 {
		t.Errorf("expected vMemory[0x2102] == 0x33, got 0x%X", vMemory.Read(0x2102))
	}
}

func TestPPUTriggerLatchWithWriteSignal_ShouldWriteDataToVideoMemoryAndIncrementVBy32_WhenOffsetBitIsSet(t *testing.T) {
	mem := memory.NewMemory()
	vMemory := memory.NewMemory()
	bus := bus.NewBusWithWorkMemory(mem)
	bus.AtatchVideoMemory(vMemory)
	mockPipeline := &PixelPipelineFixture{}
	screenFixture := NewScreenFixture()
	ppu := ppu.NewRp2C02(bus, mockPipeline, screenFixture)

	// PPU_CONTROL bit 2 = 1 (increment by 32)
	mem.Write(0x2000, 0b0000_0100)

	// Set address to 0x2100
	mem.Write(0x2006, 0x21)
	ppu.TriggerLatchWithWriteSignal(0x2006)
	mem.Write(0x2006, 0x00)
	ppu.TriggerLatchWithWriteSignal(0x2006)

	mem.Write(0x2007, 0x11)
	ppu.TriggerLatchWithWriteSignal(0x2007)

	mem.Write(0x2007, 0x22)
	ppu.TriggerLatchWithWriteSignal(0x2007)

	mem.Write(0x2007, 0x33)
	ppu.TriggerLatchWithWriteSignal(0x2007)

	if vMemory.Read(0x2100) != 0x11 {
		t.Errorf("expected vMemory[0x2100] == 0x11, got 0x%X", vMemory.Read(0x2100))
	}
	if vMemory.Read(0x2120) != 0x22 {
		t.Errorf("expected vMemory[0x2120] == 0x22, got 0x%X", vMemory.Read(0x2120))
	}
	if vMemory.Read(0x2140) != 0x33 {
		t.Errorf("expected vMemory[0x2140] == 0x33, got 0x%X", vMemory.Read(0x2140))
	}
}

func TestPPUTriggerLatchWithWriteSignal_ShouldWrapVAddress_WhenIncrementExceedsMax15Bits(t *testing.T) {
	mem := memory.NewMemory()
	vMemory := memory.NewMemory()
	bus := bus.NewBusWithWorkMemory(mem)
	bus.AtatchVideoMemory(vMemory)
	mockPipeline := &PixelPipelineFixture{}
	screenFixture := NewScreenFixture()
	ppu := ppu.NewRp2C02(bus, mockPipeline, screenFixture)

	// Increment by 32
	mem.Write(0x2000, 0b0000_0100)

	// Set address to 0x3FFF (16383)
	mem.Write(0x2006, 0x3F)
	ppu.TriggerLatchWithWriteSignal(0x2006)
	mem.Write(0x2006, 0xFF)
	ppu.TriggerLatchWithWriteSignal(0x2006)

	// Perform 512 writes with increment 32:
	// 16383 + 512 * 32 = 32767. 32767 % 32767 = 0.
	for i := 0; i < 512; i++ {
		mem.Write(0x2007, byte(i))
		ppu.TriggerLatchWithWriteSignal(0x2007)
	}

	// The 513th write should be at address 0
	mem.Write(0x2007, 0x77)
	ppu.TriggerLatchWithWriteSignal(0x2007)

	if vMemory.Read(0x0000) != 0x77 {
		t.Errorf("expected vMemory[0x0000] == 0x77 after wrapping, got 0x%X", vMemory.Read(0x0000))
	}
}

func TestPPUTriggerLatchWithWriteSignal_ShouldDoNothing_OnUnhandledAddress(t *testing.T) {
	mem := memory.NewMemory()
	vMemory := memory.NewMemory()
	bus := bus.NewBusWithWorkMemory(mem)
	bus.AtatchVideoMemory(vMemory)
	mockPipeline := &PixelPipelineFixture{}
	screenFixture := NewScreenFixture()
	ppu := ppu.NewRp2C02(bus, mockPipeline, screenFixture)

	// Addresses other than 0x2006 and 0x2007
	ppu.TriggerLatchWithWriteSignal(0x2000)
	ppu.TriggerLatchWithWriteSignal(0x2001)
	ppu.TriggerLatchWithWriteSignal(0x2002)
	ppu.TriggerLatchWithWriteSignal(0x2003)
	ppu.TriggerLatchWithWriteSignal(0x2004)
	ppu.TriggerLatchWithWriteSignal(0x2005)
	ppu.TriggerLatchWithWriteSignal(0x4014)
}

func TestPPUTriggerLatchWithReadSignal_ShouldClearVBlankFlag_OnPPUStatusRead(t *testing.T) {
	mem := memory.NewMemory()
	vMemory := memory.NewMemory()
	bus := bus.NewBusWithWorkMemory(mem)
	bus.AtatchVideoMemory(vMemory)
	mockPipeline := &PixelPipelineFixture{}
	screenFixture := NewScreenFixture()
	ppu := ppu.NewRp2C02(bus, mockPipeline, screenFixture)

	mem.Write(0x2002, 0b1000_0000)

	ppu.TriggerLatchWithReadSignal(0x2002)

	if mem.Read(0x2002) != 0b0000_0000 {
		t.Errorf("expected VBlank flag (bit 7) to be cleared, got %b", mem.Read(0x2002))
	}
}

func TestPPUTriggerLatchWithReadSignal_ShouldPreserveOtherStatusBits_OnPPUStatusRead(t *testing.T) {
	mem := memory.NewMemory()
	vMemory := memory.NewMemory()
	bus := bus.NewBusWithWorkMemory(mem)
	bus.AtatchVideoMemory(vMemory)
	mockPipeline := &PixelPipelineFixture{}
	screenFixture := NewScreenFixture()
	ppu := ppu.NewRp2C02(bus, mockPipeline, screenFixture)

	mem.Write(0x2002, 0b1111_1111)

	ppu.TriggerLatchWithReadSignal(0x2002)

	if mem.Read(0x2002) != 0b0111_1111 {
		t.Errorf("expected bits 0-6 to be preserved (0b0111_1111), got %b", mem.Read(0x2002))
	}
}

func TestPPUTriggerLatchWithReadSignal_ShouldResetWriteLatchW_OnPPUStatusRead(t *testing.T) {
	mem := memory.NewMemory()
	vMemory := memory.NewMemory()
	bus := bus.NewBusWithWorkMemory(mem)
	bus.AtatchVideoMemory(vMemory)
	mockPipeline := &PixelPipelineFixture{}
	screenFixture := NewScreenFixture()
	ppu := ppu.NewRp2C02(bus, mockPipeline, screenFixture)

	// First write to 0x2006 sets w = true
	mem.Write(0x2006, 0x21)
	ppu.TriggerLatchWithWriteSignal(0x2006)

	// Reading 0x2002 resets w back to false
	ppu.TriggerLatchWithReadSignal(0x2002)

	// Since w was reset to false, the next write is treated as the FIRST write (high byte)
	mem.Write(0x2006, 0x30)
	ppu.TriggerLatchWithWriteSignal(0x2006)

	// Second write (low byte)
	mem.Write(0x2006, 0x55)
	ppu.TriggerLatchWithWriteSignal(0x2006)

	mem.Write(0x2007, 0xEE)
	ppu.TriggerLatchWithWriteSignal(0x2007)

	if vMemory.Read(0x3055) != 0xEE {
		t.Errorf("expected vMemory[0x3055] == 0xEE, got 0x%X (vMemory[0x2130] = 0x%X)",
			vMemory.Read(0x3055), vMemory.Read(0x2130))
	}
}

func TestPPUTriggerLatchWithReadSignal_ShouldDoNothing_OnUnhandledAddress(t *testing.T) {
	mem := memory.NewMemory()
	vMemory := memory.NewMemory()
	bus := bus.NewBusWithWorkMemory(mem)
	bus.AtatchVideoMemory(vMemory)
	mockPipeline := &PixelPipelineFixture{}
	screenFixture := NewScreenFixture()
	ppu := ppu.NewRp2C02(bus, mockPipeline, screenFixture)

	mem.Write(0x2002, 0b1000_0000)

	// Addresses other than 0x2002 should not reset VBlank or affect state
	ppu.TriggerLatchWithReadSignal(0x2000)
	ppu.TriggerLatchWithReadSignal(0x2006)
	ppu.TriggerLatchWithReadSignal(0x2007)

	if mem.Read(0x2002) != 0b1000_0000 {
		t.Errorf("expected 0x2002 to remain unchanged, got %b", mem.Read(0x2002))
	}
}

func TestPPURender_ShouldPassVAndFineYToPixelPipeline(t *testing.T) {
	mem := memory.NewMemory()
	vMemory := memory.NewMemory()
	bus := bus.NewBusWithWorkMemory(mem)
	bus.AtatchVideoMemory(vMemory)
	mockPipeline := &PixelPipelineFixture{}
	screenFixture := NewScreenFixture()
	ppu := ppu.NewRp2C02(bus, mockPipeline, screenFixture)

	// Set v to 0x3123:
	// fineY is (0x3123 & 0x7000) >> 12 = 3
	mem.Write(0x2006, 0x31)
	ppu.TriggerLatchWithWriteSignal(0x2006)
	mem.Write(0x2006, 0x23)
	ppu.TriggerLatchWithWriteSignal(0x2006)

	ppu.Render()

	if mockPipeline.LastVValue != 0x3123 {
		t.Errorf("expected pipeline to receive v 0x3123, got 0x%X", mockPipeline.LastVValue)
	}

	if mockPipeline.LastFineY != 3 {
		t.Errorf("expected pipeline to receive fineY 3, got %d", mockPipeline.LastFineY)
	}
}

func TestPPURender_ShouldShowImageOnScreen_OnVBlankStart(t *testing.T) {
	mem := memory.NewMemory()
	vMemory := memory.NewMemory()
	bus := bus.NewBusWithWorkMemory(mem)
	bus.AtatchVideoMemory(vMemory)
	mockPipeline := &PixelPipelineFixture{}
	screenFixture := NewScreenFixture()
	ppu := ppu.NewRp2C02(bus, mockPipeline, screenFixture)

	// Render until start of VBlank (scanline 240)
	for i := 0; i < 336*240; i++ {
		ppu.Render()
	}

	if screenFixture.ShowImageCalls == 0 {
		t.Errorf("expected ShowImage to be called on VBlank start, got %d", screenFixture.ShowImageCalls)
	}

	if screenFixture.LastBuffer == nil {
		t.Error("expected LastBuffer to not be nil")
	}
}

