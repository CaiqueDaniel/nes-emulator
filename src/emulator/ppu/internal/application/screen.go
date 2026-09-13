package application

type Screen interface {
	ShowImage(buffer *[][]uint32)
}
