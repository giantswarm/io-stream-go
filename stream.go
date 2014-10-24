package iostream

import (
	"io"

	"github.com/juju/errgo"
)

// Stream continously reads the data from r and writes them w using io.Copy.
// The copy operation can by sending any bool to cancel.
func Stream(w io.Writer, r io.ReadCloser, cancel <-chan bool) error {
	errChan := make(chan error)

	// Execute the io.Copy asynchronously so we can wait for cancel events
	go func() {
		_, err := io.Copy(w, r)
		select {
		case errChan <- errgo.Mask(err):
		default:
		}
		close(errChan)
	}()

	// Wait for the client to close the connection
	select {
	case err := <-errChan:
		return errgo.Mask(err)
	case <-cancel:
		// Client closed the request
		return errgo.Mask(r.Close())
	}
}
