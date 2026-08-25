package tests

import (
	"bytes"
	"nes-emu/src/emulator/domain"
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
