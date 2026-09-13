package delivery

import "nes-emu/src/emulator/console/internal/application"

type ConsoleController struct {
	initialized bool
	startGame   application.StartGame
}

func NewConsoleController(startGame application.StartGame) *ConsoleController {
	return &ConsoleController{
		initialized: true,
		startGame:   startGame,
	}
}

func (c *ConsoleController) StartGame(gamePath string) error {
	c.checkIfInitilized()

	return c.startGame.Execute(application.StartGameInput{Path: gamePath})
}

func (c *ConsoleController) checkIfInitilized() {
	if !c.initialized {
		panic("ConsoleController was not initilized")
	}
}
