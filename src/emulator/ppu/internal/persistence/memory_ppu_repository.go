package persistence

import "nes-emu/src/emulator/ppu/internal/domain"

type PPUMemoryRepository struct {
	state *domain.PPU
}

func NewPPUMemoryRepository() *PPUMemoryRepository {
	return &PPUMemoryRepository{
		state: domain.NewPPU(),
	}
}

func (p *PPUMemoryRepository) Save(state domain.PPU) {
	p.state = &state
}

func (p *PPUMemoryRepository) GetState() domain.PPU {
	return *p.state
}
