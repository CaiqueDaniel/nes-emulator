package memory

type Memory struct {
	data []byte
}

func NewMemory(size uint16) *Memory {
	return &Memory{
		data: make([]byte, uint(size)+1),
	}
}

func (m *Memory) Write(address uint16, value byte) {
	m.data[address] = value
}

func (m *Memory) Read(address uint16) byte {
	return m.data[address]
}
