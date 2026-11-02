package domain

type BGLayer struct {
	tileIndex          byte
	tileAttributeIndex byte

	lowPatternLatch  byte
	highPatternLatch byte

	lowPatternShiftRegister    uint16
	highPatternShiftRegister   uint16
	lowAttributeShiftRegister  uint16
	highAttributeShiftRegister uint16
}

func NewBGLayer() *BGLayer {
	return &BGLayer{}
}

func (b *BGLayer) LoadShiftRegistersFromLatches() {
	//TODO: implement
}

func (b *BGLayer) ShiftRegisters() {
	b.lowPatternShiftRegister <<= 1
	b.highPatternShiftRegister <<= 1
	b.lowAttributeShiftRegister <<= 1
	b.highAttributeShiftRegister <<= 1
}

// Getters
func (b *BGLayer) GetTileIndex() byte {
	return b.tileIndex
}

func (b *BGLayer) GetTileAttributeIndex() byte {
	return b.tileAttributeIndex
}

func (b *BGLayer) GetLowPatternLatch() byte {
	return b.lowPatternLatch
}

func (b *BGLayer) GetHighPatternLatch() byte {
	return b.highPatternLatch
}

func (b *BGLayer) GetLowPatternShiftRegister() uint16 {
	return b.lowPatternShiftRegister
}

func (b *BGLayer) GetHighPatternShiftRegister() uint16 {
	return b.highPatternShiftRegister
}

func (b *BGLayer) GetLowAttributeShiftRegister() uint16 {
	return b.lowAttributeShiftRegister
}

func (b *BGLayer) GetHighAttributeShiftRegister() uint16 {
	return b.highAttributeShiftRegister
}

// Setters
func (b *BGLayer) SetTileIndex(tileIndex byte) {
	b.tileIndex = tileIndex
}

func (b *BGLayer) SetTileAttributeIndex(tileAttributeIndex byte) {
	b.tileAttributeIndex = tileAttributeIndex
}

func (b *BGLayer) SetLowPatternLatch(lowPatternLatch byte) {
	b.lowPatternLatch = lowPatternLatch
}

func (b *BGLayer) SetHighPatternLatch(highPatternLatch byte) {
	b.highPatternLatch = highPatternLatch
}
