package domain

type PPURepository interface {
	Save(state PPU)
	GetState() PPU
}
