package emulator

import (
	"nes-emu/src/emulator/application"
	"nes-emu/src/emulator/components/bus"
	"nes-emu/src/emulator/components/cpu"
	"nes-emu/src/emulator/components/memory"
	"nes-emu/src/emulator/components/rp2C02"
	screen_driver "nes-emu/src/emulator/components/screen"
	"nes-emu/src/emulator/delivery"
	shared_services "nes-emu/src/shared/services"

	"golang.org/x/exp/shiny/screen"
)

type EmulatorModule interface {
	GetStartGameController() delivery.EmulatorController
}

type emulatorModule struct {
	startGameController delivery.EmulatorController
}

func NewEmulatorModule(window *screen.Window, buffer *screen.Buffer) *emulatorModule {
	module := &emulatorModule{}
	module.init(window, buffer)
	return module
}

func (e *emulatorModule) init(window *screen.Window, buffer *screen.Buffer) {
	fs := shared_services.NewLocalFileSystem()
	screen := screen_driver.NewShinyScreen(window, buffer)
	bus := bus.NewBus()
	cpu := cpu.NewCpu(bus)
	ppu := rp2C02.NewRp2C02(bus, rp2C02.NewPipeline(bus), screen)

	bus.AtatchWorkMemory(memory.NewMemory())
	bus.AtatchVideoMemory(memory.NewMemory())
	bus.AttachNMI(cpu)
	bus.AttachPictureProcessingUnit(ppu)

	e.startGameController = delivery.NewEmulatorController(application.NewStartGame(fs, bus, cpu))
}

func (e *emulatorModule) GetStartGameController() delivery.EmulatorController {
	return e.startGameController
}
