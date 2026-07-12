package application

type Bus interface {
	Tick(cycles int)
	AttachTickable(tickable Tickable)
}

type MNIBus interface {
	CallNMIHandler()
}

type Tickable interface {
	Tick()
}
