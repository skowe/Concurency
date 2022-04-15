package wc

import (
	"io"
	"log"
	"os"
)

// Open file in specified path return a pointer to its object
// Raise fatal error if file can not be opened
func getFile(path string) *os.File {
	f, err := os.Open(path)
	if err != nil {
		log.Fatalln(err)
	}
	return f
}

// Send all bytes read from file in path to channel c
// use STDIN if file is not specified
// Raise fatal error if any reading fails
// Return the number of bytes read from file/STDIN
func LoadFile(path string, c chan []byte) int {
	var buf []byte
	var err error
	if path == "" {
		buf, err = io.ReadAll(os.Stdin)
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		f := getFile(path)
		buf, err = io.ReadAll(f)
		if err != nil {
			log.Fatalln(err)
		}
	}
	c <- buf
	close(c)
	return len(buf)
}
