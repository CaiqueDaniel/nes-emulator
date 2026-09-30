package application

import "nes-emu/src/emulator/ppu/internal/domain"

type UpdateStateFromIOEvents struct {
	repository      domain.PPURepository
	ioEventsContext IOEventsContext
}

func NewUpdateStateFromIOEvents(repository domain.PPURepository, ioEventsContext IOEventsContext) *UpdateStateFromIOEvents {
	return &UpdateStateFromIOEvents{repository: repository, ioEventsContext: ioEventsContext}
}

func (u *UpdateStateFromIOEvents) Execute(workMemoryAddress uint16, isWrite bool) {
	state := u.repository.GetState()

	u.ioEventsContext.HandleEvent(&IOEvent{Address: workMemoryAddress, IsWrite: isWrite}, &state)
	u.repository.Save(state)
}
