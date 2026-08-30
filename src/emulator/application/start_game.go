package application

import (
	"nes-emu/src/emulator/domain"
	shared_application "nes-emu/src/shared/application"
)

type StartGame interface {
	Execute(input StartGameInput) error
}

type startGame struct {
	fs   shared_application.FileSystem
	bus  Bus
	cpu  CPU
	vRam Memory
}

func NewStartGame(fs shared_application.FileSystem, bus Bus, cpu CPU, vRam Memory) StartGame {
	return &startGame{
		fs:   fs,
		cpu:  cpu,
		bus:  bus,
		vRam: vRam,
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
	f.loadCHRROMDataIntoMemory(rom)
	f.cpu.RunProgram()

	return nil
}

func (f *startGame) loadPRGROMDataIntoMemory(rom *domain.ROM) {
	const initial_pgr_rom_address uint16 = 0x8000

	for index, pgrByte := range rom.GetPRGROM() {
		f.bus.WriteToMemory(initial_pgr_rom_address+uint16(index), pgrByte)
	}
}

func (f *startGame) loadCHRROMDataIntoMemory(rom *domain.ROM) {
	for index, chrByte := range rom.GetCHRROM() {
		f.vRam.Write(uint16(index), chrByte)
	}
}

type StartGameInput struct {
	Path string
}
