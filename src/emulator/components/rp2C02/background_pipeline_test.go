package rp2C02

import (
	"nes-emu/src/emulator/application"
	"testing"
)

type mockBus struct {
	workMemory  map[uint16]byte
	videoMemory map[uint16]byte
}

func newMockBus() *mockBus {
	return &mockBus{
		workMemory:  make(map[uint16]byte),
		videoMemory: make(map[uint16]byte),
	}
}

func (b *mockBus) Tick() {}
func (b *mockBus) ReadFromMemory(address uint16) uint8 {
	if v, ok := b.workMemory[address]; ok {
		return v
	}
	return 0
}
func (b *mockBus) WriteToMemory(address uint16, value uint8) {
	b.workMemory[address] = value
}
func (b *mockBus) CallNMIHandler() {}
func (b *mockBus) ReadFromVideoMemory(address uint16) uint8 {
	if v, ok := b.videoMemory[address]; ok {
		return v
	}
	return 0
}
func (b *mockBus) WriteToVideoMemory(address uint16, value uint8) {
	b.videoMemory[address] = value
}

func (b *mockBus) AtatchWorkMemory(memory application.Memory)  {}
func (b *mockBus) AtatchVideoMemory(memory application.Memory) {}

func TestNewPipeline(t *testing.T) {
	bus := newMockBus()
	p := NewPipeline(bus)

	if p == nil {
		t.Fatal("expected NewPipeline to return non-nil instance")
	}

	pipelineImpl, ok := p.(*pipeline)
	if !ok {
		t.Fatal("expected instance to be *pipeline")
	}

	lowPat, highPat, lowAttr, highAttr := pipelineImpl.GetShiftRegisters()
	if lowPat != 0 || highPat != 0 || lowAttr != 0 || highAttr != 0 {
		t.Fatalf("expected all shift registers to be initialized to 0, got (%d, %d, %d, %d)",
			lowPat, highPat, lowAttr, highAttr)
	}
}

func TestStepUpPipelineSetsTileIndex(t *testing.T) {
	bus := newMockBus()
	p := NewPipeline(bus).(*pipeline)

	vValue := uint16(0x0005)
	expectedTile := byte(0xAA)
	addr := base_nametable_address | (vValue & 0x0FFF)
	bus.WriteToVideoMemory(addr, expectedTile)

	// dot 2 is fetch_tile_index_dot
	result := p.StepUpPipeline(2, vValue, 0)

	if result != false {
		t.Errorf("expected StepUpPipeline on dot 2 to return false, got %v", result)
	}
	if p.tileIndex != expectedTile {
		t.Fatalf("tileIndex expected 0x%X, got 0x%X", expectedTile, p.tileIndex)
	}
}

func TestStepUpPipelineExtractPaletteBits_AllQuadrants(t *testing.T) {
	tests := []struct {
		name         string
		coarseX      uint16
		coarseY      uint16
		attrByte     byte
		expectedLow  byte
		expectedHigh byte
	}{
		{
			name:         "top-left quadrant (shift 0), palette 0 (0b00)",
			coarseX:      0,
			coarseY:      0,
			attrByte:     0b00,
			expectedLow:  0,
			expectedHigh: 0,
		},
		{
			name:         "top-left quadrant (shift 0), palette 1 (0b01)",
			coarseX:      0,
			coarseY:      0,
			attrByte:     0b01,
			expectedLow:  255,
			expectedHigh: 0,
		},
		{
			name:         "top-left quadrant (shift 0), palette 2 (0b10)",
			coarseX:      0,
			coarseY:      0,
			attrByte:     0b10,
			expectedLow:  0,
			expectedHigh: 255,
		},
		{
			name:         "top-left quadrant (shift 0), palette 3 (0b11)",
			coarseX:      0,
			coarseY:      0,
			attrByte:     0b11,
			expectedLow:  255,
			expectedHigh: 255,
		},
		{
			name:         "top-right quadrant (shift 2), palette 1 (0b01)",
			coarseX:      2,
			coarseY:      0,
			attrByte:     0b0100,
			expectedLow:  255,
			expectedHigh: 0,
		},
		{
			name:         "bottom-left quadrant (shift 4), palette 2 (0b10)",
			coarseX:      0,
			coarseY:      2,
			attrByte:     0b100000,
			expectedLow:  0,
			expectedHigh: 255,
		},
		{
			name:         "bottom-right quadrant (shift 6), palette 3 (0b11)",
			coarseX:      2,
			coarseY:      2,
			attrByte:     0b11000000,
			expectedLow:  255,
			expectedHigh: 255,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bus := newMockBus()
			p := NewPipeline(bus).(*pipeline)

			vValue := (tt.coarseY << 5) | tt.coarseX
			attrAddr := p.getAttrTableAddress(vValue)
			bus.WriteToVideoMemory(attrAddr, tt.attrByte)

			// dot 4 is fetch_attr_table_dot
			result := p.StepUpPipeline(4, vValue, 0)

			if result != false {
				t.Errorf("expected StepUpPipeline on dot 4 to return false, got %v", result)
			}
			if p.lowPallette != tt.expectedLow {
				t.Errorf("low palette expected %d, got %d", tt.expectedLow, p.lowPallette)
			}
			if p.highPallette != tt.expectedHigh {
				t.Errorf("high palette expected %d, got %d", tt.expectedHigh, p.highPallette)
			}
		})
	}
}

func TestStepUpPipelineFetchLowerPatternByte(t *testing.T) {
	tests := []struct {
		name                 string
		ppuControlBit4       byte
		tileIndex            byte
		fineY                uint16
		expectedLowByte      byte
		expectedPatternTable uint16
	}{
		{
			name:                 "pattern table 0 (ppu_control bit 4 is 0)",
			ppuControlBit4:       0x00,
			tileIndex:            0x05,
			fineY:                2,
			expectedLowByte:      0x42,
			expectedPatternTable: 0x0000,
		},
		{
			name:                 "pattern table 1 (ppu_control bit 4 is 1)",
			ppuControlBit4:       0x10,
			tileIndex:            0x0A,
			fineY:                5,
			expectedLowByte:      0x99,
			expectedPatternTable: 0x1000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bus := newMockBus()
			p := NewPipeline(bus).(*pipeline)

			bus.WriteToMemory(ppu_control, tt.ppuControlBit4)
			p.tileIndex = tt.tileIndex

			// getByteFromPatternTable computes address: patternIndex | (tileIndex << 4) | high_byte_offset | fineY
			addr := tt.expectedPatternTable | (uint16(tt.tileIndex) << 4) | 0b1000 | tt.fineY
			bus.WriteToVideoMemory(addr, tt.expectedLowByte)

			// dot 6 is fetch_lower_pattern_table_dot
			result := p.StepUpPipeline(6, 0, tt.fineY)

			if result != false {
				t.Errorf("expected StepUpPipeline on dot 6 to return false, got %v", result)
			}
			if p.lowPatternByte != tt.expectedLowByte {
				t.Fatalf("lowPatternByte expected 0x%X, got 0x%X", tt.expectedLowByte, p.lowPatternByte)
			}
		})
	}
}

func TestStepUpPipelineFetchHigherPatternByte_AndFillsShiftRegisters(t *testing.T) {
	bus := newMockBus()
	p := NewPipeline(bus).(*pipeline)

	vValue := uint16(0x0005)
	tile := byte(0x10)
	fineY := uint16(3)

	bus.WriteToVideoMemory(base_nametable_address|(vValue&0x0FFF), tile)
	bus.WriteToVideoMemory(base_nametable_address|0x03C0, 0b01) // palette 1 -> lowPalette: 255, highPalette: 0
	bus.WriteToMemory(ppu_control, 0x10)                       // pattern table 1 (0x1000)

	highByteOffset := uint16(0b1000)
	patternAddr := (uint16(1) << 12) | (uint16(tile) << 4) | highByteOffset | fineY
	patternVal := byte(0x5A)
	bus.WriteToVideoMemory(patternAddr, patternVal)

	// Step through dot 2 (tile), 4 (attr), 6 (low pattern byte)
	p.StepUpPipeline(2, vValue, fineY)
	p.StepUpPipeline(4, vValue, fineY)
	p.StepUpPipeline(6, vValue, fineY)

	// Dot 0 (or dot 8): fetch higher pattern byte and fill shift registers
	result := p.StepUpPipeline(0, vValue, fineY)

	if result != true {
		t.Fatalf("expected StepUpPipeline on dot 0 to return true, got %v", result)
	}

	if p.highPatternByte != patternVal {
		t.Errorf("expected highPatternByte 0x%X, got 0x%X", patternVal, p.highPatternByte)
	}

	lowPat, highPat, lowAttr, highAttr := p.GetShiftRegisters()

	if lowPat != uint16(patternVal) {
		t.Errorf("expected lowPatternShiftRegister 0x%X, got 0x%X", patternVal, lowPat)
	}
	if highPat != uint16(patternVal) {
		t.Errorf("expected highPatternShiftRegister 0x%X, got 0x%X", patternVal, highPat)
	}
	if lowAttr != 255 {
		t.Errorf("expected lowAttributeShiftRegister 255, got %d", lowAttr)
	}
	if highAttr != 0 {
		t.Errorf("expected highAttributeShiftRegister 0, got %d", highAttr)
	}
}

func TestStepUpPipelineNonFetchingDots(t *testing.T) {
	bus := newMockBus()
	p := NewPipeline(bus).(*pipeline)

	nonFetchingDots := []uint{1, 3, 5, 7, 9, 11, 13, 15}
	for _, dot := range nonFetchingDots {
		result := p.StepUpPipeline(dot, 0x0000, 0)
		if result != false {
			t.Errorf("expected dot %d to return false, got %v", dot, result)
		}
	}
}

func TestRenderPixel_PaletteLookup(t *testing.T) {
	// Tests that pixelData is correctly assembled from:
	// bit 0: lowPatternBit
	// bit 1: highPatternBit
	// bit 2: lowAttrBit
	// bit 3: highAttrBit
	// and returns colorPallet[pixelData].
	tests := []struct {
		name              string
		lowPatternBit     uint16
		highPatternBit    uint16
		lowAttrBit        uint16
		highAttrBit       uint16
		expectedPixelData uint8
	}{
		{
			name:              "all zero",
			lowPatternBit:     0,
			highPatternBit:    0,
			lowAttrBit:        0,
			highAttrBit:       0,
			expectedPixelData: 0,
		},
		{
			name:              "low pattern bit set",
			lowPatternBit:     1,
			highPatternBit:    0,
			lowAttrBit:        0,
			highAttrBit:       0,
			expectedPixelData: 0b0001,
		},
		{
			name:              "high pattern bit set",
			lowPatternBit:     0,
			highPatternBit:    1,
			lowAttrBit:        0,
			highAttrBit:       0,
			expectedPixelData: 0b0010,
		},
		{
			name:              "low attribute bit set",
			lowPatternBit:     0,
			highPatternBit:    0,
			lowAttrBit:        1,
			highAttrBit:       0,
			expectedPixelData: 0b0100,
		},
		{
			name:              "high attribute bit set",
			lowPatternBit:     0,
			highPatternBit:    0,
			lowAttrBit:        0,
			highAttrBit:       1,
			expectedPixelData: 0b1000,
		},
		{
			name:              "all bits set (pixel data 15)",
			lowPatternBit:     1,
			highPatternBit:    1,
			lowAttrBit:        1,
			highAttrBit:       1,
			expectedPixelData: 15,
		},
		{
			name:              "mixed pattern and attribute (pixel data 10: 0b1010)",
			lowPatternBit:     0,
			highPatternBit:    1,
			lowAttrBit:        0,
			highAttrBit:       1,
			expectedPixelData: 0b1010,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bus := newMockBus()
			p := NewPipeline(bus).(*pipeline)

			fineX := byte(0)
			bitPos := 15 - fineX

			p.lowPatternShiftRegister = tt.lowPatternBit << bitPos
			p.highPatternShiftRegister = tt.highPatternBit << bitPos
			p.lowAttributeShiftRegister = tt.lowAttrBit << bitPos
			p.highAttributeShiftRegister = tt.highAttrBit << bitPos

			expectedColor := colorPallet[tt.expectedPixelData]
			color := p.RenderPixel(fineX)

			if color != expectedColor {
				t.Fatalf("expected color 0x%06X for pixel data %d, got 0x%06X",
					expectedColor, tt.expectedPixelData, color)
			}
		})
	}
}

func TestRenderPixel_FineXBitSelection(t *testing.T) {
	// Tests that fineX selects bit (15 - fineX)
	for fineX := byte(0); fineX <= 7; fineX++ {
		t.Run("fineX bit position", func(t *testing.T) {
			bus := newMockBus()
			p := NewPipeline(bus).(*pipeline)

			bitPos := 15 - fineX
			p.lowPatternShiftRegister = 1 << bitPos // only this bit set

			color := p.RenderPixel(fineX)
			expectedColor := colorPallet[1] // lowPatternBit = 1 -> pixelData = 1

			if color != expectedColor {
				t.Fatalf("for fineX=%d, expected color 0x%06X, got 0x%06X",
					fineX, expectedColor, color)
			}
		})
	}
}

func TestRenderPixel_ShiftsRegisters(t *testing.T) {
	bus := newMockBus()
	p := NewPipeline(bus).(*pipeline)

	p.lowPatternShiftRegister = 0b00000000_00000011
	p.highPatternShiftRegister = 0b00000000_00000101
	p.lowAttributeShiftRegister = 0b00000000_00000111
	p.highAttributeShiftRegister = 0b00000000_00001001

	p.RenderPixel(0)

	lowPat, highPat, lowAttr, highAttr := p.GetShiftRegisters()

	if lowPat != 0b00000000_00000110 {
		t.Errorf("expected lowPatternShiftRegister 0b00000000_00000110, got %b", lowPat)
	}
	if highPat != 0b00000000_00001010 {
		t.Errorf("expected highPatternShiftRegister 0b00000000_00001010, got %b", highPat)
	}
	if lowAttr != 0b00000000_00001110 {
		t.Errorf("expected lowAttributeShiftRegister 0b00000000_00001110, got %b", lowAttr)
	}
	if highAttr != 0b00000000_00010010 {
		t.Errorf("expected highAttributeShiftRegister 0b00000000_00010010, got %b", highAttr)
	}
}

func TestRenderPixel_ConsecutivePixels(t *testing.T) {
	bus := newMockBus()
	p := NewPipeline(bus).(*pipeline)

	// Preload registers so that at fineX=0, MSB (bit 15) represents the current pixel,
	// and each call to RenderPixel(0) shifts registers left to reveal the next pixel.
	// Low pattern:  1 0 1 0 1 0 1 0 (0xAA00)
	// High pattern: 0 1 0 1 0 1 0 1 (0x5500)
	// Attributes:   all 0
	p.lowPatternShiftRegister = 0xAA00
	p.highPatternShiftRegister = 0x5500

	expectedPixelData := []uint8{
		0b01, // bit 15: low=1, high=0 -> 1
		0b10, // bit 14: low=0, high=1 -> 2
		0b01, // bit 13: low=1, high=0 -> 1
		0b10, // bit 12: low=0, high=1 -> 2
		0b01, // bit 11: low=1, high=0 -> 1
		0b10, // bit 10: low=0, high=1 -> 2
		0b01, // bit 9:  low=1, high=0 -> 1
		0b10, // bit 8:  low=0, high=1 -> 2
	}

	for i, expectedData := range expectedPixelData {
		color := p.RenderPixel(0)
		expectedColor := colorPallet[expectedData]
		if color != expectedColor {
			t.Fatalf("pixel %d: expected color 0x%06X (data %d), got 0x%06X",
				i, expectedColor, expectedData, color)
		}
	}
}

func TestGetAttrTableAddress(t *testing.T) {
	bus := newMockBus()
	p := NewPipeline(bus).(*pipeline)

	tests := []struct {
		name        string
		vValue      uint16
		expectedAddr uint16
	}{
		{
			name:        "nametable 0, coarseX=0, coarseY=0",
			vValue:      0x0000,
			expectedAddr: 0x23C0,
		},
		{
			name:        "nametable 1, coarseX=0, coarseY=0",
			vValue:      0x0400,
			expectedAddr: 0x27C0,
		},
		{
			name:        "nametable 2, coarseX=0, coarseY=0",
			vValue:      0x0800,
			expectedAddr: 0x2BC0,
		},
		{
			name:        "nametable 3, coarseX=0, coarseY=0",
			vValue:      0x0C00,
			expectedAddr: 0x2FC0,
		},
		{
			name:        "nametable 0 with coarseX and coarseY offset",
			// coarseX = 4 (0b00100), coarseY = 4 (0b00100 << 5 = 0b00100_00000 = 0x0080)
			vValue:      0x0084,
			// offsetX = (0x0084 & 0xE0) >> 2 = 0x80 >> 2 = 0x20
			// offsetY = (0x0084 >> 2) & 0x03E0 = 0x21 & 0x03E0 = 0x00
			expectedAddr: 0x23C0 | 0x20,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addr := p.getAttrTableAddress(tt.vValue)
			if addr != tt.expectedAddr {
				t.Fatalf("expected address 0x%04X, got 0x%04X", tt.expectedAddr, addr)
			}
		})
	}
}
