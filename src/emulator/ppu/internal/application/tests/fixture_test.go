package tests

import "nes-emu/src/emulator/memory"

func writeChar1ToPatternTable(memory *memory.Memory) {
	patternTable := []byte{
		0x10, // Linha 0: ...#....
		0x30, // Linha 1: ..##....
		0x10, // Linha 2: ...#....
		0x10, // Linha 3: ...#....
		0x10, // Linha 4: ...#....
		0x10, // Linha 5: ...#....
		0x70, // Linha 6: .#####..
		0x00, // Linha 7: ........
		0x10, // (Plano Alto - High Byte) correspondente às mesmas linhas
		0x30,
		0x10,
		0x10,
		0x10,
		0x10,
		0x70,
		0x00,
	}

	for i, b := range patternTable {
		memory.Write(uint16(initial_pattern_address)+uint16(i), b)
	}
}
