package tailpipe

import (
	"os"
	"testing"
	"time"
)

func TestRotation(t *testing.T) {
	f, teardown := tmpfile(t, "tailpipe-go-test")
	defer teardown()

	follow, err := Open(f.Name())
	if err != nil {
		t.Fatal(err)
	}

	write(t, f, "hello, ")
	compare(t, follow, "hello, ")

	done := make(chan struct{})
	go func() {
		select {
		case <-follow.Rotated:
			t.Logf("got rotate for %s", f.Name())
		case <-time.After(time.Second):
			t.Error("did not receive rotation notification")
		}
		close(done)
	}()
	teardown()
	f, err = os.Create(f.Name())
	if err != nil {
		t.Fatal(err)
	}

	write(t, f, "world!")
	compare(t, follow, "world!")
	<-done
}

func TestDelayedRotation(t *testing.T) {
	f, teardown := tmpfile(t, "tailpipe-go-test")
	defer teardown()

	follow, err := Open(f.Name())
	if err != nil {
		t.Fatal(err)
	}

	write(t, f, "hello, ")
	compare(t, follow, "hello, ")

	teardown()
	go compare(t, follow, "world!")

	time.Sleep(time.Millisecond * 100)
	f, err = os.Create(f.Name())
	if err != nil {
		t.Fatal(err)
	}

	write(t, f, "world!")
}
