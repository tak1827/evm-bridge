package log

import (
	"io"
)

type Writer struct {
	Out io.Writer
}

func (w *Writer) Write(p []byte) (n int, err error) {
	return w.Out.Write(p)
}

func (w *Writer) SetWriter(writer io.Writer) {
	w.Out = writer
}
