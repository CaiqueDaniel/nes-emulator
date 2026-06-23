package bus

import (
	"nes-emu/src/emulator/application"
)

type bus struct {
	tickables []application.Tickable
}

func NewBus() application.Bus {
	return &bus{
		tickables: make([]application.Tickable, 0),
	}
}

func (b *bus) AttachTickable(tickable application.Tickable) {
	b.tickables = append(b.tickables, tickable)
}

func (b *bus) Tick(cycles int) {
	for _, tickable := range b.tickables {
		tickable.Tick()
	}
}
