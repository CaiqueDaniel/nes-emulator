package application

type CPU interface {
	RunProgram()
	SetNMI()
	Reset()
}
