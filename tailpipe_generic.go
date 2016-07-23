//+build !android,!darwin,!dragonfly,!freebsd,!linux,!nacl,!netbsd,!openbsd,!solaris

package tailpipe

import "os"

// We provide a stub implementation of newFile
// so that tailpipe can be built on non-unix operating
// systems without file rotation detection. Contributions
// for newFile implementation on other operating systems
// are welcome.
func newFile(oldfile *os.File) (*os.File, bool) {
	return nil, false
}
