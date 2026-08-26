package tests

import (
	"nes-emu/src/emulator/application"
	"nes-emu/src/emulator/components/bus"
	"nes-emu/src/emulator/components/cpu"
	"nes-emu/src/emulator/components/memory"
	shared_services "nes-emu/src/shared/services"
	"testing"
)

func createSut() (application.StartGame, application.Bus) {
	fs := shared_services.NewLocalFileSystem()
	memory := memory.NewMemory()
	bus := bus.NewBus(memory)
	cpu := cpu.NewCpu(bus)

	return application.NewStartGame(fs, bus, cpu), bus
}

func TestItShouldValidateNesFile(t *testing.T) {
	sut, _ := createSut()
	input := application.StartGameInput{
		Path: "./../../../../test/resources/invalid.nes",
	}

	err := sut.Execute(input)

	if err == nil {
		t.Errorf("should not start invalid game: %v", err)
	}
}

func TestItShouldLoadPGRROMIntoMemory(t *testing.T) {
	sut, bus := createSut()

	input := application.StartGameInput{
		Path: "./../../../../test/resources/PALTEST.NES",
	}

	err := sut.Execute(input)

	if err == nil {
		t.Errorf("should load game: %v", err)
	}

	checkSum := 0

	for i := uint16(0x8000); i < uint16(0xFFFF); i++ {
		checkSum += int(bus.ReadFromMemory(i))
	}

	if checkSum == 0 {
		t.Errorf("should have found a byte: %d", checkSum)
	}
}
