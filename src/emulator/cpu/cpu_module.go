package cpu

import (
	"nes-emu/src/emulator/cpu/delivery"
	"nes-emu/src/emulator/cpu/internal/application"
	shared_application "nes-emu/src/emulator/shared/application"
)

type CPUModule struct {
	initialize bool
	controller *delivery.CPUController
}

func NewCPUModule(bus shared_application.Bus) *CPUModule {
	cpu := application.NewCpu(bus)
	controller := delivery.NewCPUController(cpu)

	return &CPUModule{
		controller: controller,
	}
}

func (c *CPUModule) checkIfInitialize() {
	if !c.initialize {
		panic("CPU not initialized")
	}
}
