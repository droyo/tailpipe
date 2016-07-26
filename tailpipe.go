// Package tailpipe allows for reading normal files indefinitely.
// With the tailpipe package, code that uses the standard library's
// io.Reader interface can be transparently adapted to receive
// future updates to normal files.
package tailpipe

import (
	"errors"
	"io"
	"os"
	"time"
)

// The Follow function allows for the creation of a File with
// an underlying stream that may not implement all interfaces
// which a File implements. Such Files will return ErrNotSupported
// when this is the case.
var ErrNotSupported = errors.New("Operation not supported by underlying stream")

// A File represents an open normal file. A File is effectively of
// infinite length; all reads to the file will block until data are available,
// even if EOF on the underlying file is reached.
//
// The tailpipe package will attempt to detect when a file has been
// rotated. Programs that wish to be notified when such a rotation
// occurs should receive from the Rotated channel.
type File struct {
	r       io.Reader
	Rotated chan struct{}
}

// Read reads up to len(p) bytes into p. If end-of-file is reached,
// Read will block until new data are available. Read returns the
// number of bytes read and any errors other than io.EOF.
//
// If the underlying stream is an *os.File, Read will attempt to
// detect if it has been replaced, such as during log rotation. If
// so, Read will re-open the file at the original path provided to
// Open.
func (f *File) Read(p []byte) (n int, err error) {
	for {
		n, err = f.r.Read(p)
		if n == 0 && err == io.EOF {
			time.Sleep(time.Millisecond * 100)
			if file, ok := f.r.(*os.File); ok {
				if file, ok = newFile(file); ok {
					select {
					case f.Rotated <- struct{}{}:
					default:
					}
					f.r = file
				}
			}
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
	if err != nil {
		return nil, err
	}
	return Follow(f), err
}

// Follow converts an existing io.Reader to a
// tailpipe.File. This can be useful when opening
// files with special permissions. In general, the
// behavior of a tailpipe.File is only suitable for
// normal files; using Follow on an io.Pipe, net.Conn,
// or other non-file stream will yield undesirable
// results.
func Follow(r io.Reader) *File {
	// BUG(droyo): a File's Rotated channel should not
	// be relied upon to provide an accurate count of
	// file rotations; because the tailpipe will only perform
	// a non-blocking send on the Rotated channel,
	// a goroutine may miss a new notification while it
	// is responding to a previous notification. This is addressed
	// by buffering the channel, but can still be a problem with
	// files that are rotaetd very frequently.
	return &File{r: r, Rotated: make(chan struct{}, 1)}
}

// Name returns the name of the underlying file,
// if available. If the underlying stream does not
// have a name, Name returns an empty string.
func (f *File) Name() string {
	if v, ok := f.r.(interface {
		Name() string
	}); ok {
		return v.Name()
	}
	return ""
}

// Seek calls Seek on the underlying stream. If the
// underlying stream does not provide a Seek method,
// ErrNotSupported is returned.
func (f *File) Seek(offset int64, whence int) (int64, error) {
	if v, ok := f.r.(interface {
		Seek(int64, int) (int64, error)
	}); ok {
		return v.Seek(offset, whence)
	}
	return 0, ErrNotSupported
}

// Close closes the underlying stream
func (f *File) Close() error {
	if v, ok := f.r.(io.ReadCloser); ok {
		return v.Close()
	}
	return ErrNotSupported
}
