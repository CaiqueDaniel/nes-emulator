package application

import (
	"fmt"
	"nes-emu/src/emulator/domain"
	shared_application "nes-emu/src/shared/application"
)

type StartGame interface {
	Execute(input StartGameInput) error
}

type startGame struct {
	fs  shared_application.FileSystem
	bus Bus
	cpu CPU
}

func NewStartGame(fs shared_application.FileSystem, bus Bus, cpu CPU) StartGame {
	return &startGame{
		fs:  fs,
		cpu: cpu,
		bus: bus,
	}
}

func (f *startGame) Execute(input StartGameInput) error {
	file, err := f.fs.ReadFile(input.Path)

	if err != nil {
		return err
	}

	rom, err := domain.NewROM(file)

	if err != nil {
		return err
	}

	f.loadPRGROMDataIntoMemory(rom)

	return fmt.Errorf("invalid file: %s", input.Path)
}

func (f *startGame) loadPRGROMDataIntoMemory(rom *domain.ROM) {
	const initial_pgr_rom_address uint16 = 0x8000

	for index, pgrByte := range rom.GetPRGROM() {
		f.bus.WriteToMemory(initial_pgr_rom_address+uint16(index), pgrByte)
	}
}

type StartGameInput struct {
	Path string
}
