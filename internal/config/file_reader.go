package config

import (
	"io/fs"
	"os"
)

type FileReading interface {
	ReadFile(filename string) ([]byte, error)
	Stat(filename string) (fs.FileInfo, error)
}

// Simple facade around the os file reading functions
type FileReader struct{}

func (FileReader) ReadFile(filename string) ([]byte, error) {
	return os.ReadFile(filename)
}

func (FileReader) Stat(filename string) (fs.FileInfo, error) {
	return os.Stat(filename)
}
