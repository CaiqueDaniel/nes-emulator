package console

import (
	"nes-emu/src/emulator/console/delivery"
	"nes-emu/src/emulator/console/internal/application"
	emu_application "nes-emu/src/emulator/shared/application"
	shared_services "nes-emu/src/shared/services"
)

type ConsoleModule struct {
	initialized bool
	controller  *delivery.ConsoleController
}

func NewConsoleModule(bus emu_application.Bus, cpu emu_application.CPU) *ConsoleModule {
	fs := shared_services.NewLocalFileSystem()
	startGameUseCase := application.NewStartGame(fs, bus, cpu)
	controller := delivery.NewConsoleController(startGameUseCase)

	return &ConsoleModule{
		initialized: true,
		controller:  controller,
	}
}

func (c *ConsoleModule) checkIfInitilized() {
	if !c.initialized {
		panic("ConsoleModule was not initilized")
	}
}
