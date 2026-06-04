package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	result := make(chan string)
	go func() {
		var currentLine string
		readBuffer := make([]byte, 8)
		for {
			n, err := f.Read(readBuffer)
			if err != nil {
				if err == io.EOF {
					break
				}
				continue
			}

			if n == 0 {
				break
			}
			idx := bytes.Index(readBuffer[:n], []byte("\n"))
			if idx != -1 {
				currentLine += string(readBuffer[:idx])
				result <- currentLine
				currentLine = ""
				currentLine += string(readBuffer[idx+1 : n])
			} else {
				currentLine += string(readBuffer[:n])
			}

		}
		result <- currentLine
		close(result)

	}()

	return result

}

func main() {
	file, err := os.Open("message.txt")
	if err != nil {
		log.Fatal(err)
	}
	resultLines := getLinesChannel(file)
	for line := range resultLines {
		fmt.Printf("read: %s\n", line)
	}

}
