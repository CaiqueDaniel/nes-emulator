package application

type PPU interface {
	Render(request *PPUIOEvent)
	HandleIOEvents(address uint16, isWrite bool)
}

type PPUIOEvent struct {
	Address uint16
	IsWrite bool
}
