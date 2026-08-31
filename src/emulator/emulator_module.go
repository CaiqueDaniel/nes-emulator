package emulator

import (
	"nes-emu/src/emulator/application"
	"nes-emu/src/emulator/components/bus"
	"nes-emu/src/emulator/components/cpu"
	"nes-emu/src/emulator/components/memory"
	"nes-emu/src/emulator/components/ppu"
	"nes-emu/src/emulator/delivery"
	shared_services "nes-emu/src/shared/services"
)

type EmulatorModule interface {
	GetStartGameController() delivery.EmulatorController
}

type emulatorModule struct {
	startGameController delivery.EmulatorController
}

func NewEmulatorModule() *emulatorModule {
	module := &emulatorModule{}
	module.init()
	return module
}

func (e *emulatorModule) init() {
	fs := shared_services.NewLocalFileSystem()
	bus := bus.NewBus()
	cpu := cpu.NewCpu(bus)
	ppu := ppu.NewPPU(bus, ppu.NewPipeline(bus))

	bus.AtatchWorkMemory(memory.NewMemory())
	bus.AtatchVideoMemory(memory.NewMemory())
	bus.AttachNMI(cpu)
	bus.AttachPictureProcessingUnit(ppu)

	e.startGameController = delivery.NewEmulatorController(application.NewStartGame(fs, bus, cpu))
}
