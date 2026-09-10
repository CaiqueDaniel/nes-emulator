//go:build ignore

package emulator

import (
	"nes-emu/src/emulator/application"
	cpu "nes-emu/src/emulator/cpu/application"
	"nes-emu/src/emulator/delivery"
	ppu "nes-emu/src/emulator/ppu/application"
	memory "nes-emu/src/emulator/shared/persistance"
	emulator_shared_services "nes-emu/src/emulator/shared/services"
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
	screen := emulator_shared_services.NewShinyScreen(window, buffer)
	bus := emulator_shared_services.NewBus()
	cpu := cpu.NewCpu(bus)
	ppu := ppu.NewRenderGraphics(bus, ppu.NewPipeline(bus), screen)

	bus.AtatchWorkMemory(memory.NewMemory())
	bus.AtatchVideoMemory(memory.NewMemory())
	bus.AttachNMI(cpu)
	bus.AttachPictureProcessingUnit(ppu)

	e.startGameController = delivery.NewEmulatorController(application.NewStartGame(fs, bus, cpu))
}

func (e *emulatorModule) GetStartGameController() delivery.EmulatorController {
	return e.startGameController
}
