package domain

import (
	"bytes"
	"errors"
	"math"
)

type ROM struct {
	raw     []byte
	version int8
	prgSize uint
	chrSize uint
}

func NewROM(raw []byte) (*ROM, error) {
	rom := &ROM{
		raw: raw,
	}

	if !rom.isValidFile() {
		return nil, errors.New("file is not valid")
	}

	rom.verifyAndSetFileVersion()
	rom.setPRGSize()
	rom.setCHRSize()

	return rom, nil
}

func (r *ROM) GetRaw() []byte {
	return r.raw
}

func (r *ROM) HaveTrainer() bool {
	return r.raw[6]&0b100 != 0
}

func (r *ROM) GetPRGROMAddressRange() (uint, uint) {
	const initial_start_index = 0x10
	const trainer_size = 512

	if r.HaveTrainer() {
		startAddress := initial_start_index + trainer_size
		return uint(startAddress), uint(startAddress) + r.prgSize
	}

	return initial_start_index, initial_start_index + r.prgSize
}

func (r *ROM) GetPRGROM() []byte {
	initial_start_index, end_index := r.GetPRGROMAddressRange()
	return r.raw[initial_start_index:end_index]
}

func (r *ROM) GetCHRROMAddressRange() (uint, uint) {
	_, endPgrAddress := r.GetPRGROMAddressRange()
	return endPgrAddress + 1, endPgrAddress + 1 + r.chrSize
}

func (r *ROM) GetCHRROM() []byte {
	initial_start_index, end_index := r.GetCHRROMAddressRange()
	return r.raw[initial_start_index:end_index]
}

func (r *ROM) GetVersion() int8 {
	return r.version
}

func (r *ROM) isValidFile() bool {
	if len(r.raw) < 4 {
		return false
	}
	header := (r.raw)[0:4]
	return bytes.Equal(header[0:3], []byte{'N', 'E', 'S'}) && header[3] == 0x1A
}

func (r *ROM) verifyAndSetFileVersion() {
	if len(r.raw) < 8 {
		r.version = 1
		return
	}
	header := (r.raw)[0:8]
	if r.isValidFile() && (header[7]&0x0C) == 0x08 {
		r.version = 2
		return
	}

	r.version = 1
}

func (r *ROM) setPRGSize() {
	if r.version == 1 {
		r.prgSize = uint(r.raw[4]) * 16 * 1024
		return
	}

	highByte := uint(r.raw[9] & 0xF)
	lowByte := uint(r.raw[4])

	if highByte == 0xF {
		multiplier := float64((lowByte&0x3)*2 + 1)
		expoent := float64(lowByte >> 2)
		r.prgSize = uint(math.Pow(2, expoent)) * uint(multiplier)
		return
	}

	highByte <<= 8
	r.prgSize = (highByte | lowByte) * 1024 * 16
}

func (r *ROM) setCHRSize() {
	r.chrSize = uint(r.raw[5]) * 8 * 1024
}

func (r *ROM) GetPRGSize() uint {
	return r.prgSize
}

func (r *ROM) GetCHRSize() uint {
	return r.chrSize
}
