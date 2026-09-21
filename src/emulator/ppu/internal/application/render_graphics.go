package application

import (
	"nes-emu/src/emulator/ppu/internal/domain"
	shared "nes-emu/src/emulator/shared/application"
)

type RenderGraphics struct {
	initialized     bool
	bus             shared.MNIBus
	pipeline        PixelPipeline
	screen          Screen
	ioEventsContext IOEventsContext
	state           *domain.PPU
	buffer          [domain.MAX_FRAME_SCANLINE + 1][domain.MAX_PIXEL_PER_SCANLINE + 1]uint32
}

type RenderGraphicsInput struct {
	Address uint16
	IsWrite bool
}

func NewRenderGraphics(bus shared.MNIBus, pipeline PixelPipeline, screen Screen, ioEventsContext IOEventsContext) *RenderGraphics {
	p := &RenderGraphics{
		initialized:     true,
		pipeline:        pipeline,
		bus:             bus,
		screen:          screen,
		ioEventsContext: ioEventsContext,
		buffer:          [domain.MAX_FRAME_SCANLINE + 1][domain.MAX_PIXEL_PER_SCANLINE + 1]uint32{},
		state:           domain.NewPPU(),
	}

	return p
}

func (p *RenderGraphics) Execute(input *RenderGraphicsInput) {
	p.checkIfInitilized()

	if input != nil {
		p.ioEventsContext.HandleEvent(&IOEvent{Address: input.Address, IsWrite: input.IsWrite}, p.state)
	}

	if p.state.ShouldRenderScanlinePixel() {
		p.renderPixel()
		p.state.AdvanceToNextScanlinePixel()
	}

	p.state.AdvanceToNextScanline()
	p.fetchGraphics()
	p.updateStatusRegister()

	if p.state.IsVBlankStarted() {
		p.screen.ShowImage(&p.buffer)
	}

	if p.checkIfNMIShouldBeCalled() {
		p.bus.CallNMIHandler()
	}

	p.state.IncreaseDots()
}

func (p *RenderGraphics) GetCurrentScanline() uint16 {
	return p.state.GetCurrentScanline()
}

func (p *RenderGraphics) GetCurrentScanlinePixel() uint8 {
	return p.state.GetCurrentScanlinePixel()
}

func (p *RenderGraphics) renderPixel() {
	pixelColor := p.pipeline.RenderPixel(p.state.GetX())

	p.buffer[p.state.GetCurrentScanline()][p.state.GetCurrentScanlinePixel()] = pixelColor
}

func (p *RenderGraphics) fetchGraphics() {
	p.pipeline.StepUpPipeline(uint(p.state.GetCurrentDots()), p.state.GetV(), p.state.GetFineY())
}

func (p *RenderGraphics) updateStatusRegister() {
	if p.state.IsOnPreRender() {
		value := p.bus.ReadFromMemory(domain.PPU_STATUS) & 0b00011111
		p.bus.WriteToMemory(domain.PPU_STATUS, value)
	}

	if p.state.IsVBlankStarted() {
		value := p.bus.ReadFromMemory(domain.PPU_STATUS) ^ 0b10000000
		p.bus.WriteToMemory(domain.PPU_STATUS, value)
	}
}

func (p *RenderGraphics) checkIfNMIShouldBeCalled() bool {
	return p.isNMIFlagEnabled() && p.state.IsVBlankStarted()
}

func (p *RenderGraphics) isNMIFlagEnabled() bool {
	const nmiFlagMask = 0b10000000
	return p.bus.ReadFromMemory(domain.PPU_CONTROL)&nmiFlagMask != 0
}

func (p *RenderGraphics) checkIfInitilized() {
	if !p.initialized {
		panic("RenderGraphics was not initialized")
	}
}

func (p *RenderGraphics) readVMemory(address uint16) uint8 {
	return p.bus.ReadFromVideoMemory(address)
}

func (p *RenderGraphics) writeToVMemory(address uint16, value byte) {
	p.bus.WriteToVideoMemory(address, value)
}
