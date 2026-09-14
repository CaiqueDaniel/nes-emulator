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
	c.checkIfInitialized()
	c.cpu.RunProgram()
}

func (c *CPUController) Reset() {
	c.checkIfInitialized()
	c.cpu.Reset()
}

func (c *CPUController) SetNMI() {
	c.checkIfInitialized()
	c.cpu.SetNMI()
}

func (c *CPUController) checkIfInitialized() {
	if !c.initialize {
		panic("CPU not initialized")
	}
}
