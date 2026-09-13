package delivery

import "nes-emu/src/emulator/ppu/internal/application"

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

func (p *PPUController) Render(request *RenderRequest) {
	p.checkIfInitilized()
	p.renderGraphics.Execute(&application.RenderGraphicsInput{
		Address: request.Address,
		IsWrite: request.IsWrite,
	})
}

func (p *PPUController) checkIfInitilized() {
	if !p.initialized {
		panic("PPUController was not initilized")
	}
}
