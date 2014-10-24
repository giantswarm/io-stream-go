package iostream

import (
  "github.com/juju/errgo"
  
	"net/http"
	"io"
)

// HttpStream stems the given reader to the given ResponseWriter.
// The ResponseWriter must implement http.CloseNotifier.
// If the request is canceled, the reader will be closed.
func HttpStream(w http.ResponseWriter, r io.ReadCloser) error {
	// type assert to http.CloseNotifier to be able to handle disconnection of clients
	cn, ok := w.(http.CloseNotifier)
	if !ok {
		return errgo.Newf("response writer is not a close notifier")
	}
	closeChan := cn.CloseNotify()

  wf, err := NewWriteFlusher(w)
  if err != nil {
    return errgo.Mask(err)
  }

  return Stream(wf, r, closeChan)
}

// NewWriteFlusher ensures that the given io.Writer w also implements http.Flusher
// and returns a new WriteFlusher which flushes after every Write().
func NewWriteFlusher(w http.ResponseWriter) (*WriteFlusher, error) {
	// type assert to http.Flusher
	f, ok := w.(http.Flusher)
	if !ok {
		return nil, errgo.Newf("writer is not a flusher")
	}
	return &WriteFlusher{w: w, flusher: f}, nil
}

// WriteFlusher is a port of Docker's implementation used in their API.
// (https://github.com/docker/docker/blob/9ae3134dc9f0652ef48ec1fd445f42d8fe26de35/utils/utils.go#L269)
// It combines io.Writer and http.Flusher to enable streaming of constant data
// flows via http connections.
type WriteFlusher struct {
	w       io.Writer
	flusher http.Flusher
}

// Write flushes the data immediately after every write operation
func (wf *WriteFlusher) Write(b []byte) (n int, err error) {
	n, err = wf.w.Write(b)
	wf.flusher.Flush()
	return n, err
}

func (wf *WriteFlusher) Flush() {
	wf.flusher.Flush()
}
