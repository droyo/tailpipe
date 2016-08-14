package tailpipe_test

import (
	"bufio"
	"bytes"
	"log"
	"os"

	"aqwari.net/io/tailpipe"
)

func ExampleOpen() {
	tail, err := tailpipe.Open("/var/log/messages")
	if err != nil {
		log.Fatal(err)
	}
	defer tail.Close()
	go func() {
		for range tail.Rotated {
			log.Printf("file %s rotated; following new file", tail.Name())
		}
	}()
	scanner := bufio.NewScanner(tail)
	for scanner.Scan() {
		if bytes.Contains(scanner.Bytes(), []byte("ntpd")) {
			if _, err := os.Stdout.Write(scanner.Bytes()); err != nil {
				break
			}
		}
	}
}
