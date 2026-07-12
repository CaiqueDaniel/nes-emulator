package application

type Bus interface {
	Tick(cycles int)
	AttachTickable(tickable Tickable)
	ReadFromMemory(address uint16) uint8
	WriteToMemory(address uint16, value uint8)
}

type MNIBus interface {
	Bus
	CallNMIHandler()
}

type Tickable interface {
	Tick()
}
