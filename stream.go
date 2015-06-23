package iostream

import (
	"io"

	"github.com/juju/errgo"
)

var (
	maskAny = errgo.MaskFunc(errgo.Any)
)

// Stream continously reads the data from r and writes them w using io.Copy.
// The copy operation can by sending any bool to cancel.
// Any returned error of `io.Copy` is ignored, if the request is already canceled.
//
// github.com/juju/errgo.Mask() is used internally.
func Stream(w io.Writer, r io.ReadCloser, cancel <-chan bool) error {
	errChan := make(chan error)

	// Execute the io.Copy asynchronously so we can wait for cancel events
	go func() {
		_, err := io.Copy(w, r)
		select {
		case errChan <- maskAny(err):
		default:
		}
		close(errChan)
	}()

	// Wait for the client to close the connection
	select {
	case err := <-errChan:
		return maskAny(err)
	case <-cancel:
		// Client closed the request
		return maskAny(r.Close())
	}
}
