package mocks

type MockCloser struct {
	CloseCount int
	Err        error
}

func (m *MockCloser) Close() error {
	m.CloseCount++
	return m.Err
}