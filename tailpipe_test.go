package tailpipe

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func tmpfile(t *testing.T, prefix string) (*os.File, func()) {
	f, err := ioutil.TempFile("", prefix)
	if err != nil {
		t.Fatal(err)
	}
	return f, func() { f.Close(); os.Remove(f.Name()) }
}

func write(t *testing.T, w io.Writer, data string) {
	if _, err := io.WriteString(w, data); err != nil {
		t.Fatal(err)
	}
}

func compare(t *testing.T, r io.Reader, want string) {
	buf := make([]byte, 3000)
	done := make(chan struct{})

	var (
		n   int
		err error
	)
	go func() {
		n, err = r.Read(buf)
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(time.Second * 5):
		t.Fatal("read timeout")
	}
	if err != nil {
		t.Fatal(err)
	}
	got := string(buf[:n])
	if got != want {
		t.Errorf("Read(r) = %q, wanted %q", got, want)
	}
}

func TestFile(t *testing.T) {
	f, teardown := tmpfile(t, "tailpipe-go-test")
	defer teardown()

	follow, err := Open(f.Name())
	if err != nil {
		t.Fatal(err)
	}
	defer follow.Close()

	if follow.Name() != f.Name() {
		t.Error("Name() got %q, want %q", follow.Name(), f.Name())
	}
	write(t, f, "hello, ")
	compare(t, follow, "hello, ")

	write(t, f, "world!")
	compare(t, follow, "world!")

	go compare(t, follow, "there!")
	time.Sleep(time.Millisecond / 2)
	write(t, f, "there!")

	// Test double close
	follow.Close()
	follow.Close()

	if _, ok := <-follow.Rotated; ok {
		t.Errorf("closing %s did not close Rotated channel", f.Name())
	}
}

func TestSeeker(t *testing.T) {
	r := bytes.NewReader([]byte("hello, world!"))
	follow := Follow(r)
	defer follow.Close()

	if n, err := follow.Seek(-6, 2); err != nil {
		t.Error(err)
	} else if n != 7 {
		t.Errorf("Seek(-6, 2) got %d, wanted 7", n)
	}

	compare(t, follow, "world!")
}

func TestReader(t *testing.T) {
	var buf bytes.Buffer
	follow := Follow(&buf)
	defer follow.Close()

	if follow.Name() != "" {
		t.Errorf("Name() got %q, want \"\"", follow.Name())
	}
	if _, err := follow.Seek(0, 2); err == nil {
		t.Error("seek on unseekable *bytes.Buffer did not return error")
	}
	go compare(t, follow, "hello, ")
	time.Sleep(time.Millisecond / 2)
	write(t, &buf, "hello, ")
}
