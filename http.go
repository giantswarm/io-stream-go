package iostream

import (
	"github.com/juju/errgo"

	"io"
	"net/http"
)

func HttpStream(res http.ResponseWriter, w io.Writer, r io.ReadCloser) error {
	// type assert to http.Flusher to be able to stream journald's output
	f, ok := res.(http.Flusher)
	if !ok {
		return errgo.Newf("response writer is not a flusher")
	}
	wf := &iostream.WriteFlusher{W: w, Flusher: f}

	// type assert to http.CloseNotifier to be able to handle disconnection of clients
	cn, ok := res.(http.CloseNotifier)
	if !ok {
		return errgo.Newf("response writer is not a close notifier")
	}
	closeChan := cn.CloseNotify()

	return iostream.Stream(wf, r, closeChan)
}

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
