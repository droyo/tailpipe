// Package tailpipe allows for reading normal files indefinitely.
// With the tailpipe package, code that uses the standard library's
// io.Reader interface can be transparently adapted to read receive
// future updates to normal files. Rather than returning EOF when
// the end of file is reached, the Read routine in the tailpipe package
// waits for future updates to the file. This can be useful, for instance,
// when watching log files for updates.
package tailpipe

import (
	"io"
	"os"
	"time"
)

// A File represents an open normal file. A File is effectively of
// infinite length; all reads to the file will block until data are available,
// even if EOF on the underlying file is reached.
type File struct {
	r io.Reader
}

// Read reads up to len(p) bytes into p. If end-of-file is reached,
// Read will block until new data are available. Read returns the
// number of bytes read and any errors other than io.EOF.
func (f *File) Read(p []byte) (n int, err error) {
	for {
		n, err = f.r.Read(p)
		if n == 0 && err == io.EOF {
			time.Sleep(time.Millisecond * 100)
		} else {
			break
		}
	}
	if err == io.EOF {
		return n, nil
	}
	return n, err
}

// Open opens the given file for reading.
func Open(path string) (*File, error) {
	f, err := os.Open(path)
	return &File{r: f}, err
}
