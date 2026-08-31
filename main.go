package main

import (
	"nes-emu/src/emulator"
)

func main() {
	em := emulator.NewEmulatorModule()
	em.GetStartGameController().StartGame("./../../../test/resources/PALTEST.NES")
}
