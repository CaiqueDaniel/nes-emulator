package application

type Memory interface {
	Read(address uint16) uint8
	Write(address uint16, value uint8)
}
