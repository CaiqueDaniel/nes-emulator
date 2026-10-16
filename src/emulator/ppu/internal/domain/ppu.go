package domain

type PPU struct {
	scanline                  uint16
	dot                       uint16
	enableBackgroundRendering bool
	enableSpriteRendering     bool
}

func NewPPU() *PPU {
	return &PPU{
		scanline:                  0,
		dot:                       0,
		enableBackgroundRendering: false,
		enableSpriteRendering:     false,
	}
}

func (p *PPU) GetScanline() uint16 {
	return p.scanline
}

func (p *PPU) GetDot() uint16 {
	return p.dot
}

func (p *PPU) IsBackgroundRenderingEnabled() bool {
	return p.enableBackgroundRendering
}

func (p *PPU) IsSpriteRenderingEnabled() bool {
	return p.enableSpriteRendering
}

func (p *PPU) IsRenderingEnabled() bool {
	return p.enableBackgroundRendering || p.enableSpriteRendering
}
