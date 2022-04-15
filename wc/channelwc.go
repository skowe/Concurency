package wc

import (
	"bytes"
	"sync"
)

// Counts the number of lines, words and bytes in the file specified by path
// Does not return an error but file loading operations can cause a fatal error
// that crash the program
func ChanCount(path string) (lines, words, byts int) {
	c := make(chan []byte, 1)

	byts = LoadFile(path, c)

	lines, words = exec(c)
	return
}

// Goroutine starter for counting lines and words
func exec(c <-chan []byte) (lines, words int) {
	lineFlag := make(chan bool)
	wordFlag := make(chan bool)
	doneFlag := make(chan bool)

	defer func() {
		close(lineFlag)
		close(wordFlag)
		close(doneFlag)
	}()

	bs := <-c
	wg := &sync.WaitGroup{}
	go counter(&lines, &words, lineFlag, wordFlag, doneFlag)
	wg.Add(2)
	go countLines(bs, lineFlag, wg)
	go countWords(bs, wordFlag, wg)
	wg.Wait()
	doneFlag <- true

	return
}

// Waits for input on any of the passed channels
// reading from lineFlag increases the nuber of lines
// reading from wordFlag increases the number of words
// reading doneFlag stops the count
func counter(lines, words *int, lineFlag, wordFlag, doneFlag <-chan bool) {
	for {
		select {
		case <-lineFlag:
			*lines++
		case <-wordFlag:
			*words++
		case <-doneFlag:
			return
		}
	}
}

// Implements loginc for counting words
// When a condition for a word is met send a bool flag on the channel to increment word count
func countWords(b []byte, wordFlag chan<- bool, wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
	}()
	var prev = []byte{' '}
	for _, x := range b {
		switch x {
		case ' ', '\t', '\n':
			if !bytes.Contains([]byte{' ', '\n', '\t'}, prev) {
				wordFlag <- true
			}
			prev[0] = x
		default:
			prev[0] = x
		}
	}
}

// Implements loginc for counting lines
// When a condition for a line is met send a bool flag on the channel to increment line count
func countLines(b []byte, lineFlag chan<- bool, wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
	}()

	for _, x := range b {
		if x == '\n' {
			lineFlag <- true
		}
	}
}
