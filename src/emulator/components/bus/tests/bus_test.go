package tests

import (
	"nes-emu/src/emulator/components/bus"
	"nes-emu/src/emulator/components/memory"
	"nes-emu/test/fixtures"
	"testing"
)

type mockPPU struct {
	renderCalled int
}

func (m *mockPPU) Render() {
	m.renderCalled++
}

func TestNewBus(t *testing.T) {
	b := bus.NewBus()
	if b == nil {
		t.Fatalf("Expected NewBus to return a non-nil object")
	}

	if b.GetTickCount() != 0 {
		t.Errorf("Expected initial tickCount to be 0, got %d", b.GetTickCount())
	}
}

func TestNewBusWithWorkMemory(t *testing.T) {
	mem := memory.NewMemory()
	b := bus.NewBusWithWorkMemory(mem)
	if b == nil {
		t.Fatalf("Expected NewBusWithWorkMemory to return a non-nil object")
	}

	b.WriteToMemory(0x1234, 0x42)
	if b.ReadFromMemory(0x1234) != 0x42 {
		t.Errorf("Expected work memory to be initialized and accessible")
	}
}

func TestAttachWorkMemory(t *testing.T) {
	b := bus.NewBus()
	mem := memory.NewMemory()
	b.AtatchWorkMemory(mem)

	b.WriteToMemory(0x0100, 0xAB)
	result := b.ReadFromMemory(0x0100)
	if result != 0xAB {
		t.Errorf("Expected to read 0xAB from attached work memory, got 0x%X", result)
	}
}

func TestAttachVideoMemory(t *testing.T) {
	b := bus.NewBus()
	vmem := memory.NewMemory()
	b.AtatchVideoMemory(vmem)

	b.WriteToVideoMemory(0x2000, 0xCD)
	result := b.ReadFromVideoMemory(0x2000)
	if result != 0xCD {
		t.Errorf("Expected to read 0xCD from attached video memory, got 0x%X", result)
	}
}

func TestAttachPictureProcessingUnitAndTick(t *testing.T) {
	b := bus.NewBus()
	ppu := &mockPPU{}
	b.AttachPictureProcessingUnit(ppu)

	b.Tick()

	if ppu.renderCalled != 3 {
		t.Errorf("Expected ppu.Render to be called 3 times after 1 tick, got %d", ppu.renderCalled)
	}

	if b.GetTickCount() != 1 {
		t.Errorf("Expected tickCount to be 1, got %d", b.GetTickCount())
	}

	b.Tick()

	if ppu.renderCalled != 6 {
		t.Errorf("Expected ppu.Render to be called 6 times after 2 ticks, got %d", ppu.renderCalled)
	}

	if b.GetTickCount() != 2 {
		t.Errorf("Expected tickCount to be 2, got %d", b.GetTickCount())
	}
}

func TestResetTickCount(t *testing.T) {
	b := bus.NewBus()
	ppu := &mockPPU{}
	b.AttachPictureProcessingUnit(ppu)

	b.Tick()
	b.Tick()
	if b.GetTickCount() != 2 {
		t.Errorf("Expected tickCount to be 2, got %d", b.GetTickCount())
	}

	b.ResetTickCount()
	if b.GetTickCount() != 0 {
		t.Errorf("Expected tickCount to be reset to 0, got %d", b.GetTickCount())
	}
}

func TestAttachNMIAndCallNMIHandler(t *testing.T) {
	b := bus.NewBus()
	cpu := &fixtures.MockCPU{}

	b.AttachNMI(cpu)
	b.CallNMIHandler()

	if cpu.NmiCalled != 1 {
		t.Errorf("Expected cpu.SetNMI to be called 1 time, got %d", cpu.NmiCalled)
	}

	b.CallNMIHandler()

	if cpu.NmiCalled != 2 {
		t.Errorf("Expected cpu.SetNMI to be called 2 times, got %d", cpu.NmiCalled)
	}
}

func TestReadWorkMemory_PanicsWhenUnattached(t *testing.T) {
	b := bus.NewBus()

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected ReadFromMemory to panic when work memory is nil")
		}
	}()
	b.ReadFromMemory(0x0000)
}

func TestWriteWorkMemory_PanicsWhenUnattached(t *testing.T) {
	b := bus.NewBus()

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected WriteToMemory to panic when work memory is nil")
		}
	}()
	b.WriteToMemory(0x0000, 0xFF)
}

func TestReadVideoMemory_PanicsWhenUnattached(t *testing.T) {
	b := bus.NewBus()

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected ReadFromVideoMemory to panic when video memory is nil")
		}
	}()
	b.ReadFromVideoMemory(0x0000)
}

func TestWriteVideoMemory_PanicsWhenUnattached(t *testing.T) {
	b := bus.NewBus()

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected WriteToVideoMemory to panic when video memory is nil")
		}
	}()
	b.WriteToVideoMemory(0x0000, 0xFF)
}

func TestReadAndWriteToMemory_MemoryMirroing_RAM(t *testing.T) {
	b := bus.NewBus()
	b.AtatchWorkMemory(memory.NewMemory())

	b.WriteToMemory(0x0000, 0x12)
	b.WriteToMemory(0x7FF, 0x34)

	if b.ReadFromMemory(0x0000) != 0x12 || b.ReadFromMemory(0x800) != 0x12 || b.ReadFromMemory(0x1000) != 0x12 || b.ReadFromMemory(0x1800) != 0x12 {
		t.Errorf("Expected to read 0x12 from memory")
	}

	if b.ReadFromMemory(0x7FF) != 0x34 && b.ReadFromMemory(0x1FF) != 0x34 && b.ReadFromMemory(0x9FF) != 0x34 && b.ReadFromMemory(0x1FFF) != 0x34 {
		t.Errorf("Expected to read 0x34 from memory")
	}

	b.WriteToMemory(0x800, 0xAB)

	if b.ReadFromMemory(0x0000) != 0xAB {
		t.Errorf("Expected to read 0xAB from memory at address 0x0000, got 0x%X", b.ReadFromMemory(0x0000))
	}
}

func TestReadAndWriteToMemory_MemoryMirroing_PPULatches(t *testing.T) {
	b := bus.NewBus()
	b.AtatchWorkMemory(memory.NewMemory())

	b.WriteToMemory(0x2000, 0x12)
	b.WriteToMemory(0x2007, 0x34)

	if b.ReadFromMemory(0x2000) != 0x12 {
		t.Errorf("Expected to read 0x12 from memory")
	}

	if b.ReadFromMemory(0x2008) != 0x12 {
		t.Errorf("Expected to read 0x12 from memory")
	}

	if b.ReadFromMemory(0x2010) != 0x12 {
		t.Errorf("Expected to read 0x12 from memory")
	}

	if b.ReadFromMemory(0x2018) != 0x12 {
		t.Errorf("Expected to read 0x12 from memory")
	}

	if b.ReadFromMemory(0x2007) != 0x34 && b.ReadFromMemory(0x3FFF) != 0x34 {
		t.Errorf("Expected to read 0x34 from memory")
	}

	b.WriteToMemory(0x2008, 0xAB)

	if b.ReadFromMemory(0x2000) != 0xAB {
		t.Errorf("Expected to read 0xAB from memory at address 0x2000, got 0x%X", b.ReadFromMemory(0x2000))
	}
}
