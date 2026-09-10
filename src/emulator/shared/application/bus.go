package application

type Bus interface {
	Tick()
	AtatchWorkMemory(Memory)
	AtatchVideoMemory(Memory)
	ReadFromMemory(address uint16) uint8
	WriteToMemory(address uint16, value uint8)
	ReadFromVideoMemory(address uint16) uint8
	WriteToVideoMemory(address uint16, value uint8)
}

type MNIBus interface {
	Bus
	CallNMIHandler()
}

type Tickable interface {
	Tick()
}
