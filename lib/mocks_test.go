package lib

type MockReader struct {
	seek func(int64, int) (int64, error)
	read func(p []byte) (int, error)
}

func (r *MockReader) Seek(offset int64, whence int) (int64, error) {
	return r.seek(offset, whence)
}

func (r *MockReader) Read(p []byte) (int, error) {
	return r.read(p)
}
