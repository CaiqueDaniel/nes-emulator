package emulator

import (
	"nes-emu/src/emulator/console"
	console_delivery "nes-emu/src/emulator/console/delivery"
	"nes-emu/src/emulator/cpu"
	"nes-emu/src/emulator/ppu"
	shared_persistence "nes-emu/src/emulator/shared/persistance"
	shared_services "nes-emu/src/emulator/shared/services"

	"golang.org/x/exp/shiny/screen"
)

type EmulatorModule interface {
	GetStartGameController() *console_delivery.ConsoleController
}

type emulatorModule struct {
	consoleController *console_delivery.ConsoleController
}

func NewEmulatorModule(window *screen.Window, buffer *screen.Buffer) *emulatorModule {
	module := &emulatorModule{}
	module.init(window, buffer)
	return module
}

func (e *emulatorModule) init(window *screen.Window, buffer *screen.Buffer) {
	bus := shared_services.NewBus()
	cpu := cpu.NewCPUModule(bus).GetController()
	ppu := ppu.NewPPUModule(bus, window, buffer).GetController()

	bus.AtatchWorkMemory(shared_persistence.NewMemory())
	bus.AtatchVideoMemory(shared_persistence.NewMemory())
	bus.AttachNMI(cpu)
	bus.AttachPictureProcessingUnit(ppu)

	e.consoleController = console.NewConsoleModule(bus, cpu).GetController()
}

func (e *emulatorModule) GetStartGameController() *console_delivery.ConsoleController {
	return e.consoleController
}
