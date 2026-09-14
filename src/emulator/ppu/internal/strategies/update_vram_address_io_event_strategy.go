package strategies

import (
	"nes-emu/src/emulator/ppu/internal/domain"
	shared "nes-emu/src/emulator/shared/application"
)

type UpdateVRAMAddressIOEventStrategy struct {
	initilized bool
	bus        shared.MNIBus
}

func NewUpdateVRAMAddressIOEventStrategy(bus shared.MNIBus) *UpdateVRAMAddressIOEventStrategy {
	return &UpdateVRAMAddressIOEventStrategy{
		initilized: true,
		bus:        bus,
	}
}

func (s *UpdateVRAMAddressIOEventStrategy) Handle(state *domain.PPU) {
	s.checkIfInitilized()

	const high_byte_mask = 0x3F00

	if !state.IsWSet() {
		value := s.bus.ReadFromMemory(domain.PPU_ADDRESS)
		state.SetT((uint16(value) << 8) & high_byte_mask)
	} else {
		state.SetT(state.GetT() | uint16(s.bus.ReadFromMemory(domain.PPU_ADDRESS)))
		state.CopyTToV()
	}

	state.ToggleW()
}

func (s *UpdateVRAMAddressIOEventStrategy) checkIfInitilized() {
	if !s.initilized {
		panic("UpdateVRAMAddressIOEventStrategy not initialized")
	}
}
