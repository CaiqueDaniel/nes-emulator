package tests

import (
	"nes-emu/src/emulator/components/bus"
	"nes-emu/src/emulator/components/memory"
	"nes-emu/test/fixtures"
	"testing"
)

type mockTickable struct {
	tickCalled int
}

func (m *mockTickable) Tick() {
	m.tickCalled++
}

func TestNewBus(t *testing.T) {
	mem := memory.NewMemory()
	b := bus.NewBus(mem)
	if b == nil {
		t.Errorf("Expected NewBus to return a non-nil object")
	}
}

func TestAttachTickableAndTick(t *testing.T) {
	mem := memory.NewMemory()
	b := bus.NewBus(mem)
	mock1 := &mockTickable{}
	mock2 := &mockTickable{}

	b.AttachTickable(mock1)
	b.AttachTickable(mock2)

	b.Tick(1)

	if mock1.tickCalled != 1 {
		t.Errorf("Expected mock1.Tick to be called 1 time, got %d", mock1.tickCalled)
	}

	if mock2.tickCalled != 1 {
		t.Errorf("Expected mock2.Tick to be called 1 time, got %d", mock2.tickCalled)
	}

	b.Tick(5)

	if mock1.tickCalled != 2 {
		t.Errorf("Expected mock1.Tick to be called 2 times, got %d", mock1.tickCalled)
	}

	if mock2.tickCalled != 2 {
		t.Errorf("Expected mock2.Tick to be called 2 times, got %d", mock2.tickCalled)
	}
}

func TestAttachNMIAndCallNMIHandler(t *testing.T) {
	mem := memory.NewMemory()
	b := bus.NewBus(mem)
	cpu := &fixtures.MockCPU{}

	b.AttachNMI(cpu)
	b.CallNMIHandler()

	if cpu.NmiCalled != 1 {
		t.Errorf("Expected cpu.CallNMI to be called 1 time, got %d", cpu.NmiCalled)
	}

	b.CallNMIHandler()

	if cpu.NmiCalled != 2 {
		t.Errorf("Expected cpu.CallNMI to be called 2 times, got %d", cpu.NmiCalled)
	}
}

func TestReadFromMemory(t *testing.T) {
	mem := memory.NewMemory()
	b := bus.NewBus(mem)

	mockTick := &mockTickable{}
	b.AttachTickable(mockTick)

	mem.Write(0x1234, 0x42)

	result := b.ReadFromMemory(0x1234)

	if result != 0x42 {
		t.Errorf("Expected to read 0x42 from memory, got %d", result)
	}
}

func TestWriteToMemory(t *testing.T) {
	mem := memory.NewMemory()
	b := bus.NewBus(mem)

	mockTick := &mockTickable{}
	b.AttachTickable(mockTick)

	b.WriteToMemory(0x1234, 0x42)

	result := mem.Read(0x1234)

	if result != 0x42 {
		t.Errorf("Expected to write 0x42 to memory, got %d", result)
	}
}
