package domain

import (
	"bytes"
	"errors"
)

type ROM struct {
	raw     []byte
	version int8
	prgSize uint
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

	return rom, nil
}

func (r *ROM) GetRaw() []byte {
	return r.raw
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

	highByte := uint(r.raw[9]&0xF) << 8
	lowByte := uint(r.raw[4])
	r.prgSize = (highByte | lowByte) * 1024 * 16
}

func (r *ROM) GetPRGSize() uint {
	return r.prgSize
}
