package application

import "nes-emu/src/emulator/ppu/internal/domain"

type IOEventsStrategy interface {
	Handle(state *domain.PPU)
}

type IOEventsContext interface {
	HandleEvent(event *IOEvent, state *domain.PPU)
}

type IOEvent struct {
	Address uint16
	IsWrite bool
}
