package fixtures

type MockCPU struct {
	NmiCalled int
}

func (m *MockCPU) RunProgram() {}
func (m *MockCPU) Reset()      {}
func (m *MockCPU) SetNMI() {
	m.NmiCalled++
}
