The tailpipe package provides an implementation of `tail -F` for Go. In
effect, it produces files that are infinite length, and re-opens the
file when it is changed, such as during log rotation. It can be used
for programs that need to watch log files indefinitely.

There are several `tail` packages for Go available. However, most
packages do not preserve the io.Reader interface, instead opting for
a line-oriented approach involving `chan string`. By preserving the
`io.Reader` interface, this package composes better with the standard
library. In addition, the `tailpipe` package will detect file rotations and
re-open the underlying file, if possible. File rotation detection is
currently only provided on unix operating systems.
