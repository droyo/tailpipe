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
	scanner := bufio.NewScanner(tail)
	for scanner.Scan() {
		if bytes.Contains(scanner.Bytes(), []byte("ntpd")) {
			if _, err := os.Stdout.Write(scanner.Bytes()); err != nil {
				break
			}
		}
	}
}
