package application

type PPU interface {
	Render(request *PPUIOEvent)
}

type PPUIOEvent struct {
	Address uint16
	IsWrite bool
}
