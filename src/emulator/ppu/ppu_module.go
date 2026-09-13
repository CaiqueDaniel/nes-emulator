package ppu

import (
	"nes-emu/src/emulator/ppu/delivery"
	"nes-emu/src/emulator/ppu/internal/application"
	"nes-emu/src/emulator/ppu/internal/services"
	"nes-emu/src/emulator/ppu/internal/strategies"
	shared_application "nes-emu/src/emulator/shared/application"

	"golang.org/x/exp/shiny/screen"
)

type PPUModule struct {
	initilized bool
	controller *delivery.PPUController
}

func NewPPUModule(bus shared_application.MNIBus, window *screen.Window, buffer *screen.Buffer) *PPUModule {
	pipeline := application.NewPipeline(bus)
	screen := services.NewShinyScreen(window, buffer)
	renderGraphics := application.NewRenderGraphics(bus, pipeline, screen, strategies.NewPPUIOEventContext(bus))
	controller := delivery.NewPPUController(renderGraphics)

	return &PPUModule{
		initilized: true,
		controller: controller,
	}
}

func (p *PPUModule) checkIfInitilized() {
	if !p.initilized {
		panic("PPUModule was not initilized")
	}
}
