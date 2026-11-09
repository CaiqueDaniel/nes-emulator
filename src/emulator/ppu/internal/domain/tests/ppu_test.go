package tests

import (
	"testing"

	"nes-emu/src/emulator/ppu/internal/domain"
)

func TestNewPPUInitialState(t *testing.T) {
	ppu := domain.NewPPU()

	if ppu == nil {
		t.Fatal("NewPPU() returned nil")
	}

	if got := ppu.GetScanline(); got != 0 {
		t.Errorf("GetScanline() = %d, want 0", got)
	}

	if got := ppu.GetDot(); got != 0 {
		t.Errorf("GetDot() = %d, want 0", got)
	}

	if got := ppu.IsBackgroundRenderingEnabled(); got {
		t.Error("IsBackgroundRenderingEnabled() = true, want false")
	}

	if got := ppu.IsSpriteRenderingEnabled(); got {
		t.Error("IsSpriteRenderingEnabled() = true, want false")
	}

	if got := ppu.IsRenderingEnabled(); got {
		t.Error("IsRenderingEnabled() = true, want false")
	}
}

func TestPPUMoveToNextDot(t *testing.T) {
	ppu := domain.NewPPU()
	ppu.MoveToNextDot()

	if got := ppu.GetDot(); got != 1 {
		t.Errorf("GetDot() after MoveToNextDot() = %d, want 1", got)
	}
}

func TestPPUMoveToNextScanline(t *testing.T) {
	ppu := domain.NewPPU()

	for i := 0; i < domain.MAX_DOTS_PER_SCANLINE; i++ {
		ppu.MoveToNextDot()
	}

	if got := ppu.GetDot(); got != 0 {
		t.Errorf("GetDot() after MoveToNextDot() = %d, want 0", got)
	}

	if got := ppu.GetScanline(); got != 1 {
		t.Errorf("GetScanline() after MoveToNextDot() = %d, want 1", got)
	}
}

func TestPPUShouldWrapAroundScanlines(t *testing.T) {
	ppu := domain.NewPPU()

	for scanline := 0; scanline < domain.MAX_SCANLINES; scanline++ {
		for i := 0; i < domain.MAX_DOTS_PER_SCANLINE; i++ {
			ppu.MoveToNextDot()
		}
	}

	if got := ppu.GetDot(); got != 0 {
		t.Errorf("GetDot() after MoveToNextDot() = %d, want 0", got)
	}

	if got := ppu.GetScanline(); got != 0 {
		t.Errorf("GetScanline() after MoveToNextDot() = %d, want 0", got)
	}
}

func TestPPUShouldIdentifyVBlank(t *testing.T) {
	ppu := domain.NewPPU()

	for scanline := 0; scanline < domain.MAX_VISIBLE_SCANLINES; scanline++ {
		for i := 0; i < domain.MAX_DOTS_PER_SCANLINE; i++ {
			ppu.MoveToNextDot()
		}
	}

	if ppu.GetDot() != 0 {
		t.Errorf("Expected dot to be 0; got %d", ppu.GetDot())
	}

	if !ppu.IsVBlankStarted() {
		t.Errorf("Expected to be at v-blank")
	}
}

func TestPPUShouldShiftBitsOnBgShiftRegisters(t *testing.T) {
	ppu := domain.NewPPU()

	ppu.SetCurrentBgTileIndex(0)
	ppu.SetCurrentBgAttribute(0b01)
	ppu.SetLowBgPatternLatch(0x10)
	ppu.SetHighBgPatternLatch(0x11)
	ppu.LoadBgShiftRegistersFromLatches()

	lowPattern, highPattern, lowAttr, highAttr := ppu.GetShiftRegisters()
	if lowPattern != 0x10 || highPattern != 0x11 || lowAttr != 255 || highAttr != 0 {
		t.Errorf("Expected shift registers to be 0x10 and 0x11 and 0xFF and 0x0; got %X and %X and %X and %X", lowPattern, highPattern, lowAttr, highAttr)
	}

	ppu.MoveToNextDot()

	lowPattern, highPattern, lowAttr, highAttr = ppu.GetShiftRegisters()
	if lowPattern != (0x10<<1) || highPattern != (0x11<<1) || lowAttr != 254 || highAttr != (0<<1) {
		t.Errorf("Expected shift registers not to be 0x10 and 0x11 and 0xFF and 0x0; got %X and %X and %X and %X", lowPattern, highPattern, lowAttr, highAttr)
	}
}

func TestGetNameTableAddressGivenVZero(t *testing.T) {
	ppu := domain.NewPPU()

	if ppu.GetV() != 0 {
		t.Errorf("Expected v to be 0 at start; got %d", ppu.GetV())
	}

	if ppu.GetNameTableAddress() != 0x2000 {
		t.Errorf("Expected nametable address to be 0x2000 (at start) given v is 0; got %X", ppu.GetNameTableAddress())
	}
}

func TestGetAttributeTableAddress(t *testing.T) {
	tests := []struct {
		name            string
		coarseX         uint16
		coarseY         uint16
		nametable       uint16
		expectedAddress uint16
	}{
		{
			name:            "given nametable 0",
			coarseX:         0,
			coarseY:         0,
			nametable:       0,
			expectedAddress: 0x23C0,
		},
		{
			name:            "given nametable 1",
			coarseX:         0,
			coarseY:         0,
			nametable:       0b00_01_00000_00000,
			expectedAddress: 0x27C0,
		},
		{
			name:            "given nametable 2",
			coarseX:         0,
			coarseY:         0,
			nametable:       0b00_10_00000_00000,
			expectedAddress: 0x2BC0,
		},
		{
			name:            "given nametable 3",
			coarseX:         0,
			coarseY:         0,
			nametable:       0b00_11_00000_00000,
			expectedAddress: 0x2FC0,
		},
		{
			name:            "given nametable 0 and coarseX within same attribute byte (coarseX = 3)",
			coarseX:         3,
			coarseY:         0,
			nametable:       0,
			expectedAddress: 0x23C0,
		},
		{
			name:            "given nametable 0 and coarseX in next attribute byte (coarseX = 4)",
			coarseX:         4,
			coarseY:         0,
			nametable:       0,
			expectedAddress: 0x23C1,
		},
		{
			name:            "given nametable 0 and max coarseX (coarseX = 31)",
			coarseX:         31,
			coarseY:         0,
			nametable:       0,
			expectedAddress: 0x23C7,
		},
		{
			name:            "given nametable 0 and coarseY within same attribute row (coarseY = 3)",
			coarseX:         0,
			coarseY:         3,
			nametable:       0,
			expectedAddress: 0x23C0,
		},
		{
			name:            "given nametable 0 and coarseY in next attribute row (coarseY = 4)",
			coarseX:         0,
			coarseY:         4,
			nametable:       0,
			expectedAddress: 0x23C8,
		},
		{
			name:            "given nametable 0 and max coarseY (coarseY = 31)",
			coarseX:         0,
			coarseY:         31,
			nametable:       0,
			expectedAddress: 0x23F8,
		},
		{
			name:            "given nametable 0 with coarseX = 4 and coarseY = 4",
			coarseX:         4,
			coarseY:         4,
			nametable:       0,
			expectedAddress: 0x23C9,
		},
		{
			name:            "given nametable 0 with max coarseX and coarseY (coarseX = 31, coarseY = 31)",
			coarseX:         31,
			coarseY:         31,
			nametable:       0,
			expectedAddress: 0x23FF,
		},
		{
			name:            "given nametable 1 with coarseX = 8 and coarseY = 12",
			coarseX:         8,
			coarseY:         12,
			nametable:       0b00_01_00000_00000,
			expectedAddress: 0x27DA,
		},
		{
			name:            "given nametable 2 with coarseX = 20 and coarseY = 16",
			coarseX:         20,
			coarseY:         16,
			nametable:       0b00_10_00000_00000,
			expectedAddress: 0x2BE5,
		},
		{
			name:            "given nametable 3 with max coarseX and coarseY (coarseX = 31, coarseY = 31)",
			coarseX:         31,
			coarseY:         31,
			nametable:       0b00_11_00000_00000,
			expectedAddress: 0x2FFF,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ppu := domain.NewPPU()

			v := tt.nametable | (tt.coarseY << 5) | tt.coarseX
			ppu.SetVDirectly(v)

			if ppu.GetV() != v {
				t.Errorf("Expected v to be %d; got %d", v, ppu.GetV())
			}

			if ppu.GetAttributeTableAddress() != tt.expectedAddress {
				t.Errorf("Expected attribute table address to be %X; got %X", tt.expectedAddress, ppu.GetAttributeTableAddress())
			}
		})
	}
}

func TestGetNameTableAddressWithCoarseXAndCoarseY(t *testing.T) {
	tests := []struct {
		name            string
		coarseX         uint16
		coarseY         uint16
		nametable       uint16
		expectedAddress uint16
	}{
		{
			name:            "given nametable 0 with coarseX = 5 and coarseY = 2",
			coarseX:         5,
			coarseY:         2,
			nametable:       0,
			expectedAddress: 0x2000 | (2 << 5) | 5,
		},
		{
			name:            "given nametable 1 with coarseX = 31 and coarseY = 29",
			coarseX:         31,
			coarseY:         29,
			nametable:       0b00_01_00000_00000,
			expectedAddress: 0x2400 | (29 << 5) | 31,
		},
		{
			name:            "given nametable 3 with max coarseX and coarseY (coarseX = 31, coarseY = 31)",
			coarseX:         31,
			coarseY:         31,
			nametable:       0b00_11_00000_00000,
			expectedAddress: 0x2C00 | (31 << 5) | 31,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ppu := domain.NewPPU()

			v := tt.nametable | (tt.coarseY << 5) | tt.coarseX
			ppu.SetVDirectly(v)

			if ppu.GetV() != v {
				t.Errorf("Expected v to be %d; got %d", v, ppu.GetV())
			}

			if ppu.GetNameTableAddress() != tt.expectedAddress {
				t.Errorf("Expected nametable address to be %X; got %X", tt.expectedAddress, ppu.GetNameTableAddress())
			}
		})
	}
}

func TestGetFineYFromV(t *testing.T) {
	state := domain.NewPPU()

	//yyy NN YYYYY XXXXX
	state.SetVDirectly(0b111_00_00000_00000)

	if state.GetFineY() != 0b111 {
		t.Errorf("Expected fineY to be 0b111; got %d", state.GetFineY())
	}
}

func TestGetPalatteIndexFromAttribute(t *testing.T) {
	tests := []struct {
		name            string
		coarseX         uint16
		coarseY         uint16
		expectedPalette byte
	}{
		{
			name:            "given coarse is pointing to top-left palette",
			coarseX:         0,
			coarseY:         0,
			expectedPalette: 0b1,
		},
		{
			name:            "given coarse is pointing to top-right palette",
			coarseX:         0b00000_00010,
			coarseY:         0b00000_00000,
			expectedPalette: 0b10,
		},
		{
			name:            "given coarse is pointing to bottom-left palette",
			coarseX:         0b00000_00000,
			coarseY:         0b00010_00000,
			expectedPalette: 0b11,
		},
		{
			name:            "given coarse is pointing to bottom-right palette",
			coarseX:         0b00000_00010,
			coarseY:         0b00010_00000,
			expectedPalette: 0b0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			state := domain.NewPPU()

			//yyy NN YYYYY XXXXX
			state.SetVDirectly(tt.coarseY | tt.coarseX)
			state.SetCurrentBgAttribute(0b00_11_10_01)

			if state.GetPaletteIndexFromTileAttribute() != tt.expectedPalette {
				t.Errorf("Expected palette to be %d; got %d", tt.expectedPalette, state.GetPaletteIndexFromTileAttribute())
			}
		})
	}
}

func TestSetFineX_ShouldOnlyConsiderTheLower3Bits(t *testing.T) {
	ppu := domain.NewPPU()

	ppu.SetFineX(0xFF)

	if ppu.GetFineX() != 0b111 {
		t.Errorf("Expected fineX to be 0x7F; got %X", ppu.GetFineX())
	}

	ppu.SetFineX(0b1111_1011)

	if ppu.GetFineX() != 0b011 {
		t.Errorf("Expected fineX to be 0x7F; got %X", ppu.GetFineX())
	}
}

func TestShouldGetPaletteIndex_GivenAttributeIsNotZero(t *testing.T) {
	ppu := domain.NewPPU()

	ppu.SetLowBgPatternLatch(0xFF)
	ppu.SetHighBgPatternLatch(0xFF)
	ppu.SetCurrentBgAttribute(0b00_11_10_10)
	ppu.LoadBgShiftRegistersFromLatches()

	for range 8 {
		ppu.MoveToNextDot()
	}

	ppu.LoadBgShiftRegistersFromLatches()

	paletteIndex := ppu.GetPaletteIndex()

	if paletteIndex != 0b10_11 {
		t.Errorf("Expected palette index to be 0b10_11, got %b", paletteIndex)
	}
}

func TestShouldPushPixelColorToFrameBuffer(t *testing.T) {
	ppu := domain.NewPPU()
	ppu.PushPixelColorToFrameBuffer(0xFC)

	if ppu.GetFrameBuffer()[0][0] != 0xFC {
		t.Errorf("Expected pixel color to be 0xFC, got %X", ppu.GetFrameBuffer()[0][0])
	}
}

func TestShouldNotPushPixelColorForNonVisibleDots(t *testing.T) {
	ppu := domain.NewPPU()

	for range 257 {
		ppu.MoveToNextDot()
		ppu.PushPixelColorToFrameBuffer(0xFC)
	}

	if len(ppu.GetFrameBuffer()[0]) != 256 {
		t.Errorf("Expected pixel color to be 0xFC, got %d", len(ppu.GetFrameBuffer()[0]))
	}
}

func TestShouldNotPushPixelColorForNonVisibleLines(t *testing.T) {
	ppu := domain.NewPPU()

	for {
		ppu.MoveToNextDot()
		ppu.PushPixelColorToFrameBuffer(0xFC)

		if ppu.IsVBlankStarted() {
			break
		}
	}

	ppu.PushPixelColorToFrameBuffer(0xFC)

	if len(ppu.GetFrameBuffer()) != 240 {
		t.Errorf("Expected pixel color to be 0xFC, got %d", len(ppu.GetFrameBuffer()))
	}

	if ppu.GetFrameBuffer()[239][255] != 0xFC {
		t.Errorf("Expected pixel color to be 0xFC, got %X", ppu.GetFrameBuffer()[239][255])
	}
}
