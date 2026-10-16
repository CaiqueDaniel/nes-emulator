package domain

type PPURepository interface {
	Save(ppu PPU)
	Get() PPU
}
