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
	lowAttributeIndex, highAttributeIndex := b.getTransformedAttributeIndexBitToByte(b.tileAttributeIndex)

	b.fillShiftRegister(b.lowPatternLatch, &b.lowPatternShiftRegister)
	b.fillShiftRegister(b.highPatternLatch, &b.highPatternShiftRegister)
	b.fillShiftRegister(lowAttributeIndex, &b.lowAttributeShiftRegister)
	b.fillShiftRegister(highAttributeIndex, &b.highAttributeShiftRegister)
}

func (b *BGLayer) ShiftRegisters() {
	b.lowPatternShiftRegister <<= 1
	b.highPatternShiftRegister <<= 1
	b.lowAttributeShiftRegister <<= 1
	b.highAttributeShiftRegister <<= 1
}

func (b *BGLayer) getTransformedAttributeIndexBitToByte(value byte) (byte, byte) {
	low := value & 0b1
	high := (value & 0b10) >> 1

	var convertedLow, convertedHigh byte
	if low == 1 {
		convertedLow = 0xFF
	}
	if high == 1 {
		convertedHigh = 0xFF
	}

	return convertedLow, convertedHigh
}

func (b *BGLayer) fillShiftRegister(value byte, register *uint16) {
	*register = *register | uint16(value)
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
