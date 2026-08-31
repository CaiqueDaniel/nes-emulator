package delivery

import "nes-emu/src/emulator/application"

type EmulatorController interface {
	StartGame(string) error
}

type emulatorController struct {
	startGameUseCase application.StartGame
}

func NewEmulatorController(startGameUseCase application.StartGame) *emulatorController {
	return &emulatorController{
		startGameUseCase: startGameUseCase,
	}
}

func (e *emulatorController) StartGame(path string) error {
	return e.startGameUseCase.Execute(application.StartGameInput{
		Path: path,
	})
}
