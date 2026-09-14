package strategies

import (
	"nes-emu/src/emulator/ppu/internal/application"
	"nes-emu/src/emulator/ppu/internal/domain"
	"testing"
)

func TestNewPPUIOEventContext(t *testing.T) {
	bus := newMockBus()
	sut := NewPPUIOEventContext(bus)

	if sut == nil {
		t.Fatal("expected NewPPUIOEventContext to return a non-nil instance")
	}

	if !sut.initialized {
		t.Errorf("expected initialized to be true")
	}

	if sut.bus != bus {
		t.Errorf("expected bus to be properly set")
	}
}

func TestPPUIOEventContext_PanicWhenNotInitialized(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic when PPUIOEventContext is not initialized")
		}
		if r != "PPUIOEventContext not initialized" {
			t.Errorf("expected panic 'PPUIOEventContext not initialized', got '%v'", r)
		}
	}()

	sut := &PPUIOEventContext{}
	event := &application.IOEvent{
		Address: domain.PPU_ADDRESS,
		IsWrite: true,
	}
	state := domain.NewPPU()
	sut.HandleEvent(event, state)
}

func TestPPUIOEventContext_HandleEvent_Write_PPUAddress(t *testing.T) {
	bus := newMockBus()
	bus.WriteToMemory(domain.PPU_ADDRESS, 0x24)

	sut := NewPPUIOEventContext(bus)
	state := domain.NewPPU()

	event := &application.IOEvent{
		Address: domain.PPU_ADDRESS,
		IsWrite: true,
	}

	sut.HandleEvent(event, state)

	if !state.IsWSet() {
		t.Errorf("expected W to be set after write to PPU_ADDRESS")
	}
	if state.GetT() != 0x2400 {
		t.Errorf("expected T to be 0x2400, got 0x%04X", state.GetT())
	}
}

func TestPPUIOEventContext_HandleEvent_Write_PPUData(t *testing.T) {
	bus := newMockBus()
	bus.WriteToMemory(domain.PPU_DATA, 0x55)
	bus.WriteToMemory(domain.PPU_CONTROL, 0x00) // increment by 1

	sut := NewPPUIOEventContext(bus)
	state := domain.NewPPU()
	state.SetV(0x2000)

	event := &application.IOEvent{
		Address: domain.PPU_DATA,
		IsWrite: true,
	}

	sut.HandleEvent(event, state)

	if bus.ReadFromVideoMemory(0x2000) != 0x55 {
		t.Errorf("expected video memory at 0x2000 to be 0x55, got 0x%02X", bus.ReadFromVideoMemory(0x2000))
	}
	if state.GetV() != 0x2001 {
		t.Errorf("expected V to be incremented to 0x2001, got 0x%04X", state.GetV())
	}
}

func TestPPUIOEventContext_HandleEvent_Write_UnhandledAddress(t *testing.T) {
	bus := newMockBus()
	sut := NewPPUIOEventContext(bus)
	state := domain.NewPPU()

	event := &application.IOEvent{
		Address: domain.PPU_CONTROL, // Not handled by write strategy
		IsWrite: true,
	}

	// Should not panic, and state should remain unchanged
	sut.HandleEvent(event, state)

	if state.IsWSet() {
		t.Errorf("expected W to remain unset")
	}
	if state.GetT() != 0 {
		t.Errorf("expected T to remain 0")
	}
	if state.GetV() != 0 {
		t.Errorf("expected V to remain 0")
	}
}

func TestPPUIOEventContext_HandleEvent_Read_PPUStatus(t *testing.T) {
	bus := newMockBus()
	bus.WriteToMemory(domain.PPU_STATUS, 0x80) // VBlank flag set

	sut := NewPPUIOEventContext(bus)
	state := domain.NewPPU()
	state.ToggleW() // Set W flag

	event := &application.IOEvent{
		Address: domain.PPU_STATUS,
		IsWrite: false,
	}

	sut.HandleEvent(event, state)

	if bus.ReadFromMemory(domain.PPU_STATUS) != 0x00 {
		t.Errorf("expected PPU_STATUS vblank bit to be cleared (0x00), got 0x%02X", bus.ReadFromMemory(domain.PPU_STATUS))
	}
	if state.IsWSet() {
		t.Errorf("expected W flag to be cleared after reading PPU_STATUS")
	}
}

func TestPPUIOEventContext_HandleEvent_Read_UnhandledAddress(t *testing.T) {
	bus := newMockBus()
	sut := NewPPUIOEventContext(bus)
	state := domain.NewPPU()
	state.ToggleW()

	event := &application.IOEvent{
		Address: domain.PPU_ADDRESS, // Not handled by read strategy
		IsWrite: false,
	}

	// Should not panic, and state should remain unchanged
	sut.HandleEvent(event, state)

	if !state.IsWSet() {
		t.Errorf("expected W to remain set")
	}
}

func TestPPUIOEventContext_GetStrategyForWriteSignal(t *testing.T) {
	bus := newMockBus()
	sut := NewPPUIOEventContext(bus)
	state := domain.NewPPU()

	stratAddr := sut.getStrategyForWriteSignal(domain.PPU_ADDRESS, state)
	if _, ok := stratAddr.(*UpdateVRAMAddressIOEventStrategy); !ok {
		t.Errorf("expected *UpdateVRAMAddressIOEventStrategy for PPU_ADDRESS, got %T", stratAddr)
	}

	stratData := sut.getStrategyForWriteSignal(domain.PPU_DATA, state)
	if _, ok := stratData.(*UpdateVRAMDataIOEventStrategy); !ok {
		t.Errorf("expected *UpdateVRAMDataIOEventStrategy for PPU_DATA, got %T", stratData)
	}

	stratUnknown := sut.getStrategyForWriteSignal(domain.PPU_STATUS, state)
	if stratUnknown != nil {
		t.Errorf("expected nil strategy for PPU_STATUS on write, got %v", stratUnknown)
	}
}

func TestPPUIOEventContext_GetStrategyForReadSignal(t *testing.T) {
	bus := newMockBus()
	sut := NewPPUIOEventContext(bus)
	state := domain.NewPPU()

	stratStatus := sut.getStrategyForReadSignal(domain.PPU_STATUS, state)
	if _, ok := stratStatus.(*ResetWAndVBlankFlagIOEventStrategy); !ok {
		t.Errorf("expected *ResetWAndVBlankFlagIOEventStrategy for PPU_STATUS, got %T", stratStatus)
	}

	stratUnknown := sut.getStrategyForReadSignal(domain.PPU_ADDRESS, state)
	if stratUnknown != nil {
		t.Errorf("expected nil strategy for PPU_ADDRESS on read, got %v", stratUnknown)
	}
}
