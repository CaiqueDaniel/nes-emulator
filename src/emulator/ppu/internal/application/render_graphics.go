package application

import (
	"fmt"
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
	buffer          [][]uint32
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
		buffer:          make([][]uint32, domain.MAX_FRAME_SCANLINE+1),
		state:           domain.NewPPU(),
	}

	return p
}

func (p *RenderGraphics) Execute(input *RenderGraphicsInput) {
	p.checkIfInitilized()

	if input != nil {
		if input.Address == 0x2006 && input.IsWrite {
			fmt.Printf("Address: %X, V_Addr: %X, Value: %X\n", input.Address, p.state.GetV(), p.bus.ReadFromMemory(input.Address))
		}

		if input.Address == 0x2007 && input.IsWrite {
			fmt.Printf("Address: %X, V_Addr: %X, Value: %X\n", input.Address, p.state.GetV(), p.bus.ReadFromMemory(input.Address))
		}
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

	//fmt.Printf("Pixel Color: %x\n", p.state.GetT())

	p.buffer[p.state.GetCurrentScanline()] = append(p.buffer[p.state.GetCurrentScanline()], pixelColor)
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
