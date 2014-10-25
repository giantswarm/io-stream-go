io-stream-go
============

Helper to stream between Reader and Writer which may close at any point suddenly.

Docs at [godocs.org](http://godoc.org/github.com/giantswarm/io-stream-go).

## Use Case

For one of our services we needed to stream logging data from a backend (available via HTTP) to the frontend (requested also via HTTP).
We initially used simply `io.Copy` but ran into a few problems:

* If the backend does not provide any more data (but the stream should be keept open), `io.Copy` blocks on `r.Read()` even though the client request was closed. Thus the request to the backend ain't closed for a long time (maybe never).
* `http.ResponseWriter` uses a buffer internally so not all written data is immediately sent to the client

For this we implemented our own streaming version of `io.Copy` with the following features:

 * Still copies data with `io.Copy` (so it still supports io.WriteTo and others)
 * Watches for closing of the clients request using `net/http.CloseNotifier` and closes the backend reader
 * Any copy error after closing is ignored
 * Data to the `net/http.ResponseWriter` is automatically flushed (WriteFlusher is copied from Docker)

