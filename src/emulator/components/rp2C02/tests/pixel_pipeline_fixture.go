package tests

type PixelPipelineFixture struct {
	lowPatternShiftRegister    uint16
	highPatternShiftRegister   uint16
	lowAttributeShiftRegister  uint16
	highAttributeShiftRegister uint16
}

func (p *PixelPipelineFixture) StepUpPipeline(currentDot uint, vValue uint16, fineY uint16) bool {
	if currentDot%8 == 0 {
		p.fillPatternShiftRegister(0b00000011, 0b00000010)
		p.fillAttrShiftRegister(0b00000001, 0b00000001)
		return true
	}

	return false
}

func (p *PixelPipelineFixture) RenderPixel(fineX byte) uint32 {
	p.shiftRegisters()
	return 0x00000000
}

func (p *PixelPipelineFixture) GetShiftRegisters() (uint16, uint16, uint16, uint16) {
	return p.lowPatternShiftRegister, p.highPatternShiftRegister, p.lowAttributeShiftRegister, p.highAttributeShiftRegister
}

func (p *PixelPipelineFixture) shiftRegisters() {
	p.lowAttributeShiftRegister = p.lowAttributeShiftRegister << 1
	p.highAttributeShiftRegister = p.highAttributeShiftRegister << 1
	p.lowPatternShiftRegister = p.lowPatternShiftRegister << 1
	p.highPatternShiftRegister = p.highPatternShiftRegister << 1
}

func (p *PixelPipelineFixture) fillPatternShiftRegister(highByte byte, lowByte byte) {
	p.fillShiftRegister(highByte, &p.highPatternShiftRegister)
	p.fillShiftRegister(lowByte, &p.lowPatternShiftRegister)
}

func (p *PixelPipelineFixture) fillAttrShiftRegister(highByte byte, lowByte byte) {
	p.fillShiftRegister(highByte, &p.highAttributeShiftRegister)
	p.fillShiftRegister(lowByte, &p.lowAttributeShiftRegister)
}

func (p *PixelPipelineFixture) fillShiftRegister(value byte, register *uint16) {
	*register = *register | uint16(value)
}
