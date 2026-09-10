package application

import (
	"nes-emu/src/emulator/application"
	"nes-emu/src/emulator/ppu/domain"
	shared "nes-emu/src/emulator/shared/application"
)

type RenderGraphics struct {
	bus      shared.MNIBus
	pipeline application.PixelPipeline
	screen   application.Screen
	state    *domain.PPU
	buffer   [][]uint32
}

func NewRenderGraphics(bus shared.MNIBus, pipeline application.PixelPipeline, screen application.Screen) *RenderGraphics {
	p := &RenderGraphics{
		pipeline: pipeline,
		bus:      bus,
		screen:   screen,
		buffer:   make([][]uint32, domain.MAX_FRAME_SCANLINE+1),
		state:    domain.NewPPU(),
	}

	return p
}

func (p *RenderGraphics) Execute() {
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

func (p *RenderGraphics) TriggerLatchWithWriteSignal(address uint16) {
	switch address {
	case domain.PPU_ADDRESS:
		p.updateVRAMAddress()

	case domain.PPU_DATA:
		p.updateVRAMData()
	}
}

func (p *RenderGraphics) TriggerLatchWithReadSignal(address uint16) {
	switch address {
	case domain.PPU_STATUS:
		p.resetWAndVBlankFlagOnRead()
	}
}

func (p *RenderGraphics) GetCurrentScanline() uint16 {
	return p.state.GetCurrentScanline()
}

func (p *RenderGraphics) GetCurrentScanlinePixel() uint8 {
	return p.state.GetCurrentScanlinePixel()
}

func (p *RenderGraphics) renderPixel() {
	p.buffer[p.state.GetCurrentScanline()] = append(p.buffer[p.state.GetCurrentScanline()], p.pipeline.RenderPixel(p.state.GetX()))
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

// TODO 1: separate this logic into independent component
func (p *RenderGraphics) updateVRAMAddress() {
	const high_byte_mask = 0x3F00

	if !p.state.IsWSet() {
		value := p.bus.ReadFromMemory(domain.PPU_ADDRESS)
		p.state.SetT((uint16(value) << 8) & high_byte_mask)
	} else {
		p.state.SetT(p.state.GetT() | uint16(p.bus.ReadFromMemory(domain.PPU_ADDRESS)))
		p.state.CopyTToV()
	}

	p.state.ToggleW()
}

//TODO 1

// TODO 2: separate this logic into independent component
func (p *RenderGraphics) updateVRAMData() {
	value := p.bus.ReadFromMemory(domain.PPU_DATA)
	p.writeToVMemory(p.state.GetV(), value)
	p.increaseVByOffset()
}

func (p *RenderGraphics) increaseVByOffset() {
	if p.isOffsetIncrementBy32() {
		p.state.IncreaseVByOffset(32)
	} else {
		p.state.IncreaseVByOffset(1)
	}
}

func (p *RenderGraphics) isOffsetIncrementBy32() bool {
	return p.bus.ReadFromMemory(domain.PPU_CONTROL)&0b100 != 0
}

//TODO 2

//TODO 3separate this logic into independent component

func (p *RenderGraphics) resetWAndVBlankFlagOnRead() {
	status := p.bus.ReadFromMemory(domain.PPU_STATUS)
	p.bus.WriteToMemory(domain.PPU_STATUS, status&0x7F)
	p.state.ClearW()
}

//TODO 3

func (p *RenderGraphics) readVMemory(address uint16) uint8 {
	return p.bus.ReadFromVideoMemory(address)
}

func (p *RenderGraphics) writeToVMemory(address uint16, value byte) {
	p.bus.WriteToVideoMemory(address, value)
}
