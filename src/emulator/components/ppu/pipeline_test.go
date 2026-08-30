package ppu

import (
	"nes-emu/src/emulator/application"
	"testing"
)

type mockMemory struct {
	data map[uint16]byte
}

func newMockMemory() *mockMemory {
	return &mockMemory{data: make(map[uint16]byte)}
}

func (m *mockMemory) Read(address uint16) uint8 {
	if v, ok := m.data[address]; ok {
		return v
	}
	return 0
}

func (m *mockMemory) Write(address uint16, value uint8) {
	m.data[address] = value
}

type mockBus struct {
	memory map[uint16]byte
}

func newMockBus() *mockBus {
	return &mockBus{memory: make(map[uint16]byte)}
}

func (b *mockBus) Tick()                                        {}
func (b *mockBus) AttachTickable(tickable application.Tickable) {}
func (b *mockBus) ReadFromMemory(address uint16) uint8 {
	if v, ok := b.memory[address]; ok {
		return v
	}
	return 0
}
func (b *mockBus) WriteToMemory(address uint16, value uint8) { b.memory[address] = value }
func (b *mockBus) CallNMIHandler()                           {}
func (b *mockBus) ReadFromVideoMemory(address uint16) uint8 {
	if v, ok := b.memory[address]; ok {
		return v
	}
	return 0
}
func (b *mockBus) WriteToVideoMemory(address uint16, value uint8) {
	b.memory[address] = value
}

func (b *mockBus) AtatchWorkMemory(memory application.Memory)  {}
func (b *mockBus) AtatchVideoMemory(memory application.Memory) {}

// Test that when the pipeline is advanced to the tile index fetch dot it reads
// the correct tile index from nametable memory.
func TestStepUpPipelineSetsTileIndex(t *testing.T) {
	mem := newMockMemory()
	bus := newMockBus()
	p := NewPipeline(mem, bus).(*pipeline)

	vValue := uint16(0x0005)
	expectedTile := byte(0xAA)
	addr := base_nametable_address | (vValue & 0x0FFF)
	mem.Write(addr, expectedTile)

	// dot 2 is fetch_tile_index_dot
	p.StepUpPipeline(2, vValue, 0)

	if p.tileIndex != expectedTile {
		t.Fatalf("tileIndex expected 0x%X, got 0x%X", expectedTile, p.tileIndex)
	}
}

func TestStepUpPipelineExtractPaletteBits(t *testing.T) {
	mem := newMockMemory()
	bus := newMockBus()

	p := NewPipeline(mem, bus).(*pipeline)

	vValue := uint16(0x0005)
	// Attribute table for nametable 0 is at base_nametable_address + 0x03C0 when vValue small
	attrAddr := uint16(base_nametable_address | 0x03C0)
	// choose an attribute byte with known low two bits
	attrByte := byte(0b10)
	mem.Write(attrAddr, attrByte)

	p.StepUpPipeline(4, vValue, 0)

	// per extractPaletteBits for this vValue shift == 0 so palette is attrByte & 0x03
	expectedLow := byte(0)
	expectedHigh := byte(255) // when bit set code returns 255

	if p.lowPallette != expectedLow {
		t.Fatalf("low palette expected %d, got %d", expectedLow, p.lowPallette)
	}
	if p.highPallette != expectedHigh {
		t.Fatalf("high palette expected %d, got %d", expectedHigh, p.highPallette)
	}
}

func TestStepUpPipelineReturnsPipelineResult(t *testing.T) {
	mem := newMockMemory()
	bus := newMockBus()

	p := NewPipeline(mem, bus)

	vValue := uint16(0x0005)
	// set tile index at nametable
	tile := byte(0x10)
	mem.Write(base_nametable_address|(vValue&0x0FFF), tile)

	// set attribute table (palette) at default location
	mem.Write(base_nametable_address|0x03C0, 0b01)

	// set control register so pattern table index bit is 1
	bus.memory[ppu_control] = 0x10

	fineY := uint16(3)

	// compute address used by the current implementation (note: code sets high_byte_offset unconditionally)
	patternIndex := uint16(1) << 12
	highByteOffset := uint16(0b1000)
	address := patternIndex | (uint16(tile) << 4) | highByteOffset | fineY
	lowVal := byte(0x12)
	// because of an implementation detail both low and high read the same address
	mem.Write(address, lowVal)

	// advance pipeline in order: 2 -> 4 -> 6 -> 0
	p.StepUpPipeline(2, vValue, fineY)
	p.StepUpPipeline(4, vValue, fineY)
	p.StepUpPipeline(6, vValue, fineY)
	res := p.StepUpPipeline(0, vValue, fineY)

	if res == nil {
		t.Fatal("expected PipelineResult, got nil")
	}

	if res.LowPatternByte != lowVal {
		t.Fatalf("LowPatternByte expected 0x%X, got 0x%X", lowVal, res.LowPatternByte)
	}
	if res.HighPatternByte != lowVal {
		t.Fatalf("HighPatternByte expected 0x%X, got 0x%X (implementation reads same address)", lowVal, res.HighPatternByte)
	}
}
