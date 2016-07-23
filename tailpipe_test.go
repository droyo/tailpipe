package tailpipe

import (
	"io"
	"io/ioutil"
	"os"
	"testing"
)

func tmpfile(t *testing.T, prefix string) *os.File {
	f, err := ioutil.TempFile("", prefix)
	if err != nil {
		t.Fatal(err)
	}
	return f
}

func write(t *testing.T, w io.Writer, data string) {
	if _, err := io.WriteString(w, data); err != nil {
		t.Fatal(err)
	}
}

func compare(t *testing.T, r io.Reader, want string) {
	buf := make([]byte, 3000)
	n, err := r.Read(buf)
	if err != nil {
		t.Fatal(err)
	}
	got := string(buf[:n])
	if got != want {
		t.Errorf("Read(r) = %q, wanted %q", got, want)
	}
}

func TestFile(t *testing.T) {
	f := tmpfile(t, "tailpipe-go-test")
	defer f.Close()

	follow, err := Open(f.Name())
	if err != nil {
		t.Fatal(err)
	}

	write(t, f, "hello, ")
	compare(t, follow, "hello, ")

	write(t, f, "world!")
	compare(t, follow, "world!")
}
