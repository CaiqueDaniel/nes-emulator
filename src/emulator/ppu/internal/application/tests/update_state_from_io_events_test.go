package tests

import (
	"nes-emu/src/emulator/ppu/internal/application"
	"nes-emu/src/emulator/ppu/internal/persistence"
	"nes-emu/src/emulator/ppu/internal/strategies"
	shared_application "nes-emu/src/emulator/shared/application"
	shared_persistence "nes-emu/src/emulator/shared/persistance"
	shared_services "nes-emu/src/emulator/shared/services"
	"testing"
)

func TestShouldBeAbleToSetVRegister(t *testing.T) {
	updateVRegister, repository, bus := createSut()

	bus.WriteToMemory(0x2006, 0x3F)
	updateVRegister.Execute(0x2006, true)

	bus.WriteToMemory(0x2006, 0)
	updateVRegister.Execute(0x2006, true)

	state := repository.GetState()

	if state.GetV() != 0x3F00 {
		t.Errorf("V register is not working %X", state.GetV())
	}
}

func createSut() (*application.UpdateStateFromIOEvents, *persistence.PPUMemoryRepository, shared_application.Bus) {
	bus := shared_services.NewBus()
	bus.AtatchWorkMemory(shared_persistence.NewMemory())
	bus.AtatchVideoMemory(shared_persistence.NewMemory())

	repository := persistence.NewPPUMemoryRepository()
	ioEventsContext := strategies.NewPPUIOEventContext(bus)

	return application.NewUpdateStateFromIOEvents(repository, ioEventsContext), repository, bus
}
