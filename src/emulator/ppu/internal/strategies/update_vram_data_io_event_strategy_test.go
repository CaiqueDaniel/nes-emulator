package strategies

import (
	"nes-emu/src/emulator/ppu/internal/domain"
	"testing"
)

func TestNewUpdateVRAMDataIOEventStrategy(t *testing.T) {
	bus := newMockBus()
	sut := NewUpdateVRAMDataIOEventStrategy(bus)

	if sut == nil {
		t.Fatal("expected NewUpdateVRAMDataIOEventStrategy to return a non-nil instance")
	}

	if !sut.initialized {
		t.Errorf("expected initialized to be true")
	}

	if sut.bus != bus {
		t.Errorf("expected bus to be properly set")
	}
}

func TestUpdateVRAMDataIOEventStrategy_PanicWhenNotInitialized(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic when UpdateVRAMDataIOEventStrategy is not initialized")
		}
		if r != "UpdateVRAMDataIOEventStrategy not initialized" {
			t.Errorf("expected panic 'UpdateVRAMDataIOEventStrategy not initialized', got '%v'", r)
		}
	}()

	sut := &UpdateVRAMDataIOEventStrategy{}
	state := domain.NewPPU()
	sut.Handle(state)
}

func TestUpdateVRAMDataIOEventStrategy_Handle_WritesDataToVideoMemory(t *testing.T) {
	bus := newMockBus()
	sut := NewUpdateVRAMDataIOEventStrategy(bus)
	state := domain.NewPPU()

	targetVRAMAddress := uint16(0x2150)
	state.SetV(targetVRAMAddress)

	// Set data value in PPU_DATA register (0x2007)
	expectedValue := uint8(0x7F)
	bus.WriteToMemory(domain.PPU_DATA, expectedValue)

	sut.Handle(state)

	actualValue := bus.ReadFromVideoMemory(targetVRAMAddress)
	if actualValue != expectedValue {
		t.Errorf("expected video memory at 0x%04X to be 0x%02X, got 0x%02X", targetVRAMAddress, expectedValue, actualValue)
	}
}

func TestUpdateVRAMDataIOEventStrategy_Handle_IncrementsVBy1_WhenIncrementBitClear(t *testing.T) {
	testCases := []struct {
		name       string
		controlVal uint8
	}{
		{name: "control is 0", controlVal: 0x00},
		{name: "control has other bits set but bit 2 is 0", controlVal: 0b11111011},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			bus := newMockBus()
			bus.WriteToMemory(domain.PPU_CONTROL, tc.controlVal)

			sut := NewUpdateVRAMDataIOEventStrategy(bus)
			state := domain.NewPPU()
			state.SetV(0x2000)

			sut.Handle(state)

			if state.GetV() != 0x2001 {
				t.Errorf("expected V to be incremented by 1 (0x2001), got 0x%04X", state.GetV())
			}
		})
	}
}

func TestUpdateVRAMDataIOEventStrategy_Handle_IncrementsVBy32_WhenIncrementBitSet(t *testing.T) {
	testCases := []struct {
		name       string
		controlVal uint8
	}{
		{name: "only bit 2 is set (0x04)", controlVal: 0b00000100},
		{name: "all bits set (0xFF)", controlVal: 0xFF},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			bus := newMockBus()
			bus.WriteToMemory(domain.PPU_CONTROL, tc.controlVal)

			sut := NewUpdateVRAMDataIOEventStrategy(bus)
			state := domain.NewPPU()
			state.SetV(0x2000)

			sut.Handle(state)

			if state.GetV() != 0x2020 {
				t.Errorf("expected V to be incremented by 32 (0x2020), got 0x%04X", state.GetV())
			}
		})
	}
}

func TestUpdateVRAMDataIOEventStrategy_ReadAndWriteVMemory(t *testing.T) {
	bus := newMockBus()
	sut := NewUpdateVRAMDataIOEventStrategy(bus)

	address := uint16(0x3000)
	value := uint8(0xBE)

	sut.writeToVMemory(address, value)

	readBack := sut.readVMemory(address)
	if readBack != value {
		t.Errorf("expected readVMemory to return 0x%02X, got 0x%02X", value, readBack)
	}
}

func TestUpdateVRAMDataIOEventStrategy_IsOffsetIncrementBy32(t *testing.T) {
	bus := newMockBus()
	sut := NewUpdateVRAMDataIOEventStrategy(bus)

	bus.WriteToMemory(domain.PPU_CONTROL, 0x00)
	if sut.isOffsetIncrementBy32() {
		t.Errorf("expected isOffsetIncrementBy32 to return false when bit 2 is 0")
	}

	bus.WriteToMemory(domain.PPU_CONTROL, 0x04)
	if !sut.isOffsetIncrementBy32() {
		t.Errorf("expected isOffsetIncrementBy32 to return true when bit 2 is 1")
	}
}
