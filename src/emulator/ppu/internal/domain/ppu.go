package domain

const (
	MAX_DOTS_PER_SCANLINE    = 341
	MAX_VISIBLE_DOTS         = 256
	MAX_SCANLINES            = 262
	MAX_VISIBLE_SCANLINES    = 240
	INITIAL_V_BLANK_SCANLINE = 241
	LAST_V_BLANK_SCANLINE    = 261
)

var ColorPallet = [64]uint32{
	0x7C7C7C, 0x0000FC, 0x0000BC, 0x4428BC, 0x940084, 0xA80020, 0xA81000, 0x881400,
	0x503000, 0x007800, 0x006800, 0x005800, 0x004058, 0x000000, 0x000000, 0x000000,
	0xBCBCBC, 0x0070EC, 0x3C40FC, 0x7C00FA, 0xA800B4, 0xC43800, 0xE04000, 0xE86814,
	0x9C9400, 0x44B400, 0x34C400, 0x00D840, 0x00BCBC, 0x000000, 0x000000, 0x000000,
	0x3CBCFC, 0x0078F8, 0x0058F8, 0x6844FC, 0xD800CC, 0xE40058, 0xF83800, 0xE45C10,
	0xAC7C00, 0x00B800, 0x00A800, 0x00A844, 0x008888, 0x000000, 0x000000, 0x000000,
	0xF8F8F8, 0xA4E4FC, 0xB8B8F8, 0xD8B8F8, 0xF8B8F8, 0xF8A4C0, 0xF0D0B0, 0xFCE0A8,
	0xE8D878, 0xD8F878, 0xB8F8B8, 0xB8F8D8, 0x00FCFC, 0xD8D8D8, 0x000000, 0x000000,
}

type PPU struct {
	scanline                  uint16
	dot                       uint16
	enableBackgroundRendering bool
	enableSpriteRendering     bool
	frameBuffer               *PPUFrameBuffer
	v                         uint16
	t                         uint16
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

func (p *PPU) GetV() uint16 {
	return p.v
}

func (p *PPU) GetNameTableAddress() uint16 {
	const nametable_base_address = 0x2000
	const nametable_v_register_mask = 0x0FFF
	return nametable_base_address | (p.v & nametable_v_register_mask)
}
