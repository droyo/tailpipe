package tailpipe

import (
	"os"
	"testing"
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

	teardown()
	f, err = os.Create(f.Name())
	if err != nil {
		t.Fatal(err)
	}

	write(t, f, "world!")
	compare(t, follow, "world!")
}
