package delivery

import "nes-emu/src/emulator/cpu/internal/application"

type CPUController struct {
	initialize bool
	cpu        *application.CPU
}

func NewCPUController(cpu *application.CPU) *CPUController {
	return &CPUController{
		initialize: true,
		cpu:        cpu,
	}
}

func (c *CPUController) RunProgram() {
	c.checkIfInitialize()
	c.cpu.RunProgram()
}

func (c *CPUController) SetNMI() {
	c.checkIfInitialize()
	c.cpu.SetNMI()
}

func (c *CPUController) checkIfInitialize() {
	if !c.initialize {
		panic("CPU not initialized")
	}
}
