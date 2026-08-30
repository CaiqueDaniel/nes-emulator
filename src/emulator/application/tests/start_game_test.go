package tests

import (
	"nes-emu/src/emulator/application"
	"nes-emu/src/emulator/components/bus"
	"nes-emu/src/emulator/components/cpu"
	"nes-emu/src/emulator/components/memory"
	shared_services "nes-emu/src/shared/services"
	"testing"
)

func createSut() (application.StartGame, application.Bus, application.Memory) {
	fs := shared_services.NewLocalFileSystem()
	mem := memory.NewMemory()
	videoMemory := memory.NewMemory()
	bus := bus.NewBus(mem)
	cpu := cpu.NewCpu(bus)

	return application.NewStartGame(fs, bus, cpu, videoMemory), bus, videoMemory
}

func TestItShouldValidateNesFile(t *testing.T) {
	sut, _, _ := createSut()
	input := application.StartGameInput{
		Path: "./../../../../test/resources/invalid.nes",
	}

	err := sut.Execute(input)

	if err == nil {
		t.Errorf("should not start invalid game: %v", err)
	}
}

func TestItShouldLoadPGRROMIntoMemory(t *testing.T) {
	sut, bus, _ := createSut()

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

func TestItShouldLoadCHRROMIntoMemory(t *testing.T) {
	sut, _, vRam := createSut()

	input := application.StartGameInput{
		Path: "./../../../../test/resources/Zelda.NES",
	}

	err := sut.Execute(input)

	if err == nil {
		t.Errorf("should load game: %v", err)
	}

	checkSum := 0

	for i := uint16(0x0000); i < uint16(0x2000); i++ {
		checkSum += int(vRam.Read(i))
	}

	if checkSum == 0 {
		t.Errorf("should have found a byte: %d", checkSum)
	}
}
