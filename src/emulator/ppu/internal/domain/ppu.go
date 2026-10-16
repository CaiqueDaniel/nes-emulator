package domain

const (
	MAX_DOTS_PER_SCANLINE    = 341
	MAX_VISIBLE_DOTS         = 256
	MAX_SCANLINES            = 262
	MAX_VISIBLE_SCANLINES    = 240
	INITIAL_V_BLANK_SCANLINE = 241
	LAST_V_BLANK_SCANLINE    = 261
)

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

func (p *PPU) MoveToNextDot() {
	p.dot = (p.dot + 1) % MAX_DOTS_PER_SCANLINE

	if p.dot == 0 {
		p.scanline = (p.scanline + 1) % MAX_SCANLINES
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
