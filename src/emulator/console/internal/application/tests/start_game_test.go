package tests

import (
	"nes-emu/src/emulator/console/internal/application"
	shared_application "nes-emu/src/emulator/shared/application"
	"nes-emu/src/emulator/shared/persistance"
	shared_console_services "nes-emu/src/emulator/shared/services"
	shared_services "nes-emu/src/shared/services"
	"nes-emu/test/fixtures"
	"testing"
)

func createSut() (application.StartGame, shared_application.Bus, shared_application.Memory, *fixtures.MockCPU) {
	fs := shared_services.NewLocalFileSystem()
	mem := persistance.NewMemory()
	videoMemory := persistance.NewMemory()
	bus := shared_console_services.NewBus()
	cpu := &fixtures.MockCPU{}

	bus.AtatchWorkMemory(mem)
	bus.AtatchVideoMemory(videoMemory)

	return application.NewStartGame(fs, bus, cpu), bus, videoMemory, cpu
}

func TestItShouldValidateNesFile(t *testing.T) {
	sut, _, _, _ := createSut()
	input := application.StartGameInput{
		Path: "./../../../../test/resources/invalid.nes",
	}

	err := sut.Execute(input)

	if err == nil {
		t.Errorf("should not start invalid game: %v", err)
	}
}

func TestItShouldLoadPGRROMIntoMemory(t *testing.T) {
	sut, bus, _, _ := createSut()

	input := application.StartGameInput{
		Path: "./../../../../../../test/resources/PALTEST.NES",
	}

	err := sut.Execute(input)

	if err != nil {
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
	sut, _, vRam, _ := createSut()

	input := application.StartGameInput{
		Path: "./../../../../../../test/resources/Zelda.NES",
	}

	err := sut.Execute(input)

	if err != nil {
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

func TestItShouldRunTheProgram(t *testing.T) {
	sut, _, _, cpu := createSut()

	input := application.StartGameInput{
		Path: "./../../../../../../test/resources/Zelda.NES",
	}

	err := sut.Execute(input)

	if err != nil {
		t.Errorf("should load game: %v", err)
	}

	if cpu.RunProgramCalled == 0 {
		t.Errorf("should run the program")
	}
}
