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
	frameBuffer               *PPUFrameBuffer
}

type PPUFrameBuffer [MAX_VISIBLE_SCANLINES][MAX_VISIBLE_DOTS]uint32

func NewPPU() *PPU {
	return &PPU{
		scanline:                  0,
		dot:                       0,
		enableBackgroundRendering: false,
		enableSpriteRendering:     false,
		frameBuffer:               &PPUFrameBuffer{},
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

func (p *PPU) IsVBlankStarted() bool {
	return p.scanline+1 == INITIAL_V_BLANK_SCANLINE && p.dot == 0
}

func (p *PPU) ClearFrameBuffer() {
	*p.frameBuffer = PPUFrameBuffer{}
}

func (p *PPU) GetFrameBuffer() *PPUFrameBuffer {
	return p.frameBuffer
}
