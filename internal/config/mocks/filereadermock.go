package mocks

import (
	"io/fs"
)

type FileReadingMock struct {
	Invocations []string

	Results struct {
		ReadFile struct {
			Success []byte
			Error   error
		}
		Stat struct {
			Success fs.FileInfo
			Error   error
		}
	}
}

const (
	InvocationReadFile = "ReadFile"
	InvocationStat     = "Stat"
)

func (mock *FileReadingMock) Reset() {
	mock.Invocations = nil
}

func (mock *FileReadingMock) ReadFile(filename string) ([]byte, error) {
	mock.Invocations = append(mock.Invocations, InvocationReadFile)
	return mock.Results.ReadFile.Success, mock.Results.ReadFile.Error
}

func (mock *FileReadingMock) Stat(filename string) (fs.FileInfo, error) {
	mock.Invocations = append(mock.Invocations, InvocationStat)
	return mock.Results.Stat.Success, mock.Results.Stat.Error
}
