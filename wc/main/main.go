package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/skowe/concurency/wc"
)

// Main is only responsible for handling user input
// And starting/Handling a goroutine for each file passed, if any

func main() {

	// TODO: Implement flag functionalities
	// TODO: Implement Total counter if 2 or more files are passed

	wg := &sync.WaitGroup{}

	if len(os.Args) < 2 {
		printOut("", nil)
	} else {
		files := os.Args[1:]
		for _, filePath := range files {
			//TODO: add a mechanism for timing out a goroutine execution so output is in the same order as input
			wg.Add(1)
			go printOut(filePath, wg) // At this point goroutiunes are started at random
		}
		wg.Wait()
	}

}

// Starts a counter for the specific file.
// Prints out a formated output in order lines words bytes filename
func printOut(filePath string, wg *sync.WaitGroup) {
	//TODO: Use channels to send back counters to main function
	// in order to implement the total functionality
	defer func() {
		// wg is nil if source is os.Stdin
		if wg != nil {
			wg.Done()
		}
	}()
	lines, words, bytes := wc.ChanCount(filePath)
	format := "\t%d\t%d\t%d %s"
	message := fmt.Sprintf(format, lines, words, bytes, filePath)
	fmt.Println(message)
}
