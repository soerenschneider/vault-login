package main

import (
	"context"
	"errors"
	"os"
	"strings"
)

type FileWriter struct {
	fileName string
}

func NewFileWriter(cfg Config) (*FileWriter, error) {
	if strings.TrimSpace(cfg.OutputSecretName) == "" {
		return nil, errors.New("empty file name provided")
	}

	return &FileWriter{
		fileName: cfg.OutputSecretName,
	}, nil
}

func (w *FileWriter) Write(_ context.Context, data []byte) error {
	return os.WriteFile(w.fileName, data, 0640)
}
