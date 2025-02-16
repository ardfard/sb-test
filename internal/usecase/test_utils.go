package usecase

import "io"

type mockReadCloser struct {
	io.Reader
}

func (m mockReadCloser) Close() error {
	return nil
}
