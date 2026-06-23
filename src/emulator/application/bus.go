package application

type Bus interface {
	Tick(cycles int)
	AttachTickable(tickable Tickable)
}

type Tickable interface {
	Tick()
}
