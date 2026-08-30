package fixtures

type MockCPU struct {
	NmiCalled        int
	RunProgramCalled int
}

func (m *MockCPU) RunProgram() {
	m.RunProgramCalled++
}
func (m *MockCPU) Reset() {}
func (m *MockCPU) SetNMI() {
	m.NmiCalled++
}
