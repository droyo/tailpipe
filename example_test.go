package tailpipe_test

import (
	"bufio"
	"bytes"
	"log"
	"os"

	"aqwari.net/io/tailpipe"
)

func ExampleOpen() {
	f, err := tailpipe.Open("/var/log/messages")
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if bytes.Contains(scanner.Bytes(), []byte("ntpd")) {
			if _, err := os.Stdout.Write(scanner.Bytes()); err != nil {
				break
			}
		}
	}
}
