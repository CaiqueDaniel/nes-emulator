package fixtures

type NMIBusMock struct {
	NMICalled int
}

func (n *NMIBusMock) CallNMIHandler() {
	n.NMICalled++
}
