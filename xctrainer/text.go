package xctrainer

import (
	"bufio"
	"log"
	"strings"
)

func Lines(r *bufio.Reader) (c chan string) {
	c = make(chan string)
	go func() {
		defer close(c)
		for {
			line, err := r.ReadString('\n')
			if err != nil {
				if err.Error() != "EOF" {
					log.Fatal(err)
				}
				return
			}
			line = strings.TrimRight(line, "\r\n")
			if line == "" {
				continue
			}
			c <- line
		}
	}()
	return c
}
