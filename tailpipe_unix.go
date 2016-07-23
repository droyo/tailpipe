// +build android darwin dragonfly freebsd linux nacl netbsd openbsd solaris

package tailpipe

import (
	"os"
	"syscall"
)

func newFile(oldfile *os.File) (*os.File, bool) {
	var isOld bool

	oldstat, err := oldfile.Stat()
	if err != nil {
		isOld = true
	}

	newfile, err := os.Open(oldfile.Name())
	if err != nil {
		// NOTE(droyo) time will tell whether this is the right thing
		// to do. The file could be gone for good, or we could just
		// be in-between rotations.
		return nil, false
	}
	if isOld {
		return newfile, true
	}

	newstat, err := newfile.Stat()
	if err != nil {
		return nil, false
	}

	old := oldstat.Sys().(*syscall.Stat_t)
	new := newstat.Sys().(*syscall.Stat_t)

	if old.Dev == new.Dev && old.Ino == new.Ino {
		return nil, false
	}
	return newfile, true
}
