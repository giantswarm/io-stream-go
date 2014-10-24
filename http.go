package iostream

import (
	"github.com/juju/errgo"

	"io"
	"net/http"
)

// WriteFlusher is a port of Docker's implementation used in their API.
// (https://github.com/docker/docker/blob/9ae3134dc9f0652ef48ec1fd445f42d8fe26de35/utils/utils.go#L269)
// It combines io.Writer and http.Flusher to enable streaming of constant data
// flows via http connections.
type WriteFlusher struct {
	W       io.Writer
	Flusher http.Flusher
}

// Write flushes the data immediately after every write operation
func (wf *WriteFlusher) Write(b []byte) (n int, err error) {
	n, err = wf.W.Write(b)
	wf.Flusher.Flush()
	return n, err
}

func (wf *WriteFlusher) Flush() {
	wf.Flusher.Flush()
}
