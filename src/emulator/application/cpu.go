package application

type CPU interface {
	RunProgram()
	Reset()
	SetNMI()
}
