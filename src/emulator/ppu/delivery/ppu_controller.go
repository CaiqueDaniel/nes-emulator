package delivery

import (
	application "nes-emu/src/emulator/ppu/internal/application"
	shared_application "nes-emu/src/emulator/shared/application"
)

type PPUController struct {
	initialized             bool
	renderGraphics          *application.RenderGraphics
	updateStateFromIOEvents *application.UpdateStateFromIOEvents
}

type RenderRequest application.RenderGraphicsInput

func NewPPUController(renderGraphics *application.RenderGraphics, updateStateFromIOEvents *application.UpdateStateFromIOEvents) *PPUController {
	return &PPUController{
		initialized:             true,
		renderGraphics:          renderGraphics,
		updateStateFromIOEvents: updateStateFromIOEvents,
	}
}

func (p *PPUController) Render(request *shared_application.PPUIOEvent) {
	p.checkIfInitilized()

	var input *application.RenderGraphicsInput

	if request != nil {
		input = &application.RenderGraphicsInput{
			Address: request.Address,
			IsWrite: request.IsWrite,
		}
	}

	p.renderGraphics.Execute(input)
}

func (p *PPUController) HandleIOEvents(address uint16, isWrite bool) {
	p.updateStateFromIOEvents.Execute(address, isWrite)
}

func (p *PPUController) GetV() uint16 {
	return p.renderGraphics.GetV()
}

func (p *PPUController) checkIfInitilized() {
	if !p.initialized {
		panic("PPUController was not initilized")
	}
}
