package strategies

import (
	"nes-emu/src/emulator/ppu/internal/domain"
	"testing"
)

func TestNewResetWAndVBlankFlagIOEventStrategy(t *testing.T) {
	bus := newMockBus()
	sut := NewResetWAndVBlankFlagIOEventStrategy(bus)

	if sut == nil {
		t.Fatal("expected NewResetWAndVBlankFlagIOEventStrategy to return a non-nil instance")
	}

	if !sut.initialized {
		t.Errorf("expected initialized to be true")
	}

	if sut.bus != bus {
		t.Errorf("expected bus to be properly set")
	}
}

func TestResetWAndVBlankFlagIOEventStrategy_PanicWhenNotInitialized(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic when ResetWAndVBlankFlagIOEventStrategy is not initialized")
		}
	}()

	sut := &ResetWAndVBlankFlagIOEventStrategy{}
	state := domain.NewPPU()
	sut.Handle(state)
}

func TestResetWAndVBlankFlagIOEventStrategy_Handle_ClearsVBlankFlagAndPreservesOtherBits(t *testing.T) {
	testCases := []struct {
		name          string
		initialStatus uint8
		expectedAfter uint8
	}{
		{
			name:          "clears bit 7 when only bit 7 is set",
			initialStatus: 0x80,
			expectedAfter: 0x00,
		},
		{
			name:          "clears bit 7 and preserves bits 6 and 5",
			initialStatus: 0xE0,
			expectedAfter: 0x60,
		},
		{
			name:          "clears bit 7 and preserves arbitrary lower bits",
			initialStatus: 0xA5,
			expectedAfter: 0x25,
		},
		{
			name:          "leaves status unchanged if bit 7 was already clear",
			initialStatus: 0x7F,
			expectedAfter: 0x7F,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			bus := newMockBus()
			bus.WriteToMemory(domain.PPU_STATUS, tc.initialStatus)

			sut := NewResetWAndVBlankFlagIOEventStrategy(bus)
			state := domain.NewPPU()

			sut.Handle(state)

			actual := bus.ReadFromMemory(domain.PPU_STATUS)
			if actual != tc.expectedAfter {
				t.Errorf("expected PPU_STATUS to be 0x%02X, got 0x%02X", tc.expectedAfter, actual)
			}
		})
	}
}

func TestResetWAndVBlankFlagIOEventStrategy_Handle_ClearsWFlag(t *testing.T) {
	bus := newMockBus()
	sut := NewResetWAndVBlankFlagIOEventStrategy(bus)
	state := domain.NewPPU()

	state.ToggleW()
	if !state.IsWSet() {
		t.Fatal("expected W to be set before strategy execution")
	}

	sut.Handle(state)

	if state.IsWSet() {
		t.Errorf("expected W to be cleared after Handle, but it was still set")
	}

	// Calling handle again when W is already cleared should keep W cleared
	sut.Handle(state)
	if state.IsWSet() {
		t.Errorf("expected W to remain cleared")
	}
}
