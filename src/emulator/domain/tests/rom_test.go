package tests

import (
	"bytes"
	"nes-emu/src/emulator/domain"
	shared_services "nes-emu/src/shared/services"
	"testing"
)

func TestNewROM_ValidFile_iNESv1(t *testing.T) {
	rawHeader := []byte{'N', 'E', 'S', 0x1A, 0x01, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

	rom, err := domain.NewROM(rawHeader)

	if err != nil {
		t.Fatalf("Expected no error for valid NES file, got: %v", err)
	}

	if rom == nil {
		t.Fatal("Expected ROM instance, got nil")
	}

	if !bytes.Equal(rom.GetRaw(), rawHeader) {
		t.Errorf("Expected raw bytes to match input header")
	}

	if rom.GetVersion() != 1 {
		t.Errorf("Expected version 1 for standard iNES header, got %d", rom.GetVersion())
	}
}

func TestNewROM_ValidFile_NES20(t *testing.T) {
	// NES 2.0 header flag: (header[7] & 0x0C) == 0x08
	rawHeader := []byte{'N', 'E', 'S', 0x1A, 0x01, 0x01, 0x00, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

	rom, err := domain.NewROM(rawHeader)

	if err != nil {
		t.Fatalf("Expected no error for NES 2.0 file, got: %v", err)
	}

	if rom.GetVersion() != 2 {
		t.Errorf("Expected version 2 for NES 2.0 header, got %d", rom.GetVersion())
	}
}

func TestNewROM_InvalidMagicHeader(t *testing.T) {
	rawHeader := []byte{'F', 'A', 'K', 'E', 0x1A, 0x01, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

	rom, err := domain.NewROM(rawHeader)

	if err == nil {
		t.Error("Expected error for invalid magic header, got nil")
	}

	if rom != nil {
		t.Errorf("Expected nil ROM when error occurs, got %v", rom)
	}

	if err.Error() != "file is not valid" {
		t.Errorf("Expected error message 'file is not valid', got '%s'", err.Error())
	}
}

func TestNewROM_InvalidHeaderTerminator(t *testing.T) {
	// Valid "NES" letters but 4th byte is 0x00 instead of 0x1A
	rawHeader := []byte{'N', 'E', 'S', 0x00, 0x01, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

	_, err := domain.NewROM(rawHeader)

	if err == nil {
		t.Error("Expected error when 4th byte is not 0x1A, got nil")
	}
}

func TestNewROM_ShortRawBuffer(t *testing.T) {
	shortBuffers := [][]byte{
		{},
		{'N'},
		{'N', 'E'},
		{'N', 'E', 'S'},
	}

	for _, buf := range shortBuffers {
		rom, err := domain.NewROM(buf)

		if err == nil {
			t.Errorf("Expected error for short buffer of length %d, got nil", len(buf))
		}

		if rom != nil {
			t.Errorf("Expected nil ROM for short buffer of length %d, got %v", len(buf), rom)
		}
	}
}

func TestNewROM_ShouldGetNumberOfPRGBlocksOnNES1File(t *testing.T) {
	file := getTestFile()
	rom, _ := domain.NewROM(*file)

	if rom.GetVersion() != 1 {
		t.Errorf("Expected version 1 for standard iNES header, got %d", rom.GetVersion())
	}

	if rom.GetPRGSize() != 32768 {
		t.Errorf("Expected 32768 bytes for standard iNES header, got %d", rom.GetPRGSize())
	}
}

func TestNewROM_ShouldGetNumberOfPRGBlocksOnNES20File(t *testing.T) {
	file := []byte{'N', 'E', 'S', 0x1A, 0x01, 0x01, 0x00, 0x08, 0x00, 0xE, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	rom, _ := domain.NewROM(file)

	//byte 4=1; byte 9=0xE; ((0xE<<8)+1) * 1024 * 16
	const expected_size = 58736640

	if rom.GetVersion() != 2 {
		t.Errorf("Expected version 2 for standard iNES header, got %d", rom.GetVersion())
	}

	if rom.GetPRGSize() != expected_size {
		t.Errorf("Expected %d bytes for standard iNES header, got %d", expected_size, rom.GetPRGSize())
	}
}

func TestNewROM_ShouldGetPRGROMWithoutTrainerData(t *testing.T) {
	file := getTestFile()
	rom, _ := domain.NewROM(*file)

	const expected_start_address = 0x10
	const expected_end_address = 0x10 + 32768

	startAddress, endAddress := rom.GetPRGROMAddressRange()
	prgBytes := rom.GetPRGROM()

	if rom.HaveTrainer() {
		t.Errorf("rom should not have a trainer")
	}

	if startAddress != expected_start_address {
		t.Errorf("start address should be %d", expected_start_address)
	}

	if endAddress != expected_end_address {
		t.Errorf("end address should be %d", expected_end_address)
	}

	if uint(len(prgBytes)) != rom.GetPRGSize() {
		t.Errorf("prg bytes should be %d. got %d", rom.GetPRGSize(), len(prgBytes))
	}
}

func TestNewROM_ShouldGetPRGROMWithTrainerData(t *testing.T) {
	file := getTestWithTrainerFile()
	rom, _ := domain.NewROM(*file)

	const expected_start_address = 0x10 + 512
	const expected_end_address = expected_start_address + 32768

	startAddress, endAddress := rom.GetPRGROMAddressRange()
	prgBytes := rom.GetPRGROM()

	if !rom.HaveTrainer() {
		t.Errorf("rom should have a trainer")
	}

	if startAddress != expected_start_address {
		t.Errorf("start address should be %d", expected_start_address)
	}

	if endAddress != expected_end_address {
		t.Errorf("end address should be %d", expected_end_address)
	}

	if uint(len(prgBytes)) != rom.GetPRGSize() {
		t.Errorf("prg bytes should be %d. got %d", rom.GetPRGSize(), len(prgBytes))
	}
}

func getTestFile() *[]byte {
	fs := shared_services.NewLocalFileSystem()
	file, err := fs.ReadFile("./../../../../test/resources/PALTEST.NES")

	if err != nil {
		panic("file not loaded")
	}

	return &file
}

func getTestWithTrainerFile() *[]byte {
	fs := shared_services.NewLocalFileSystem()
	file, err := fs.ReadFile("./../../../../test/resources/PALTEST_With_Trainer.NES")

	if err != nil {
		panic("file not loaded")
	}

	return &file
}
