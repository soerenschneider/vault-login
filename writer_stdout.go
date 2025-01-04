package main

import (
	"context"
	"fmt"
)

type StdoutWriter struct {
}

func NewStdoutWriter() *StdoutWriter {
	return &StdoutWriter{}
}

func (w *StdoutWriter) Write(_ context.Context, data []byte) error {
	fmt.Println(data)
	return nil
}
