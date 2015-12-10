package iso

func NewWriter() Writer {
	return FakeWriter{}
}
