package delivery

import (
	application "nes-emu/src/emulator/ppu/internal/application"
	shared_application "nes-emu/src/emulator/shared/application"
)

type PPUController struct {
	initialized    bool
	renderGraphics *application.RenderGraphics
}

type RenderRequest application.RenderGraphicsInput

func NewPPUController(renderGraphics *application.RenderGraphics) *PPUController {
	return &PPUController{
		initialized:    true,
		renderGraphics: renderGraphics,
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

func (p *PPUController) checkIfInitilized() {
	if !p.initialized {
		panic("PPUController was not initilized")
	}
}
