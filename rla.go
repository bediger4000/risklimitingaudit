package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {

	fin, closefn, err := openInputFile()
	if err != nil {
		log.Fatal(err)
	}
	defer closefn()

	scanAllines(fin, timestampParser)
}

// scanAllines calls a function (argument fn) on all lines
// of fin argument one at a time. Can print some error messages
// on os.Stderr.
func scanAllines(fin *os.File, fn func(string) error) {

	scanner := bufio.NewScanner(fin)
	/* For longer lines:
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)
	*/

	lineCounter := 0

	for scanner.Scan() {
		lineCounter++
		line := scanner.Text()
		if err := fn(line); err != nil {
			fmt.Fprintf(os.Stderr, "line %d: %v\n", lineCounter, err)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "problem line %d: %v", lineCounter, err)
	}
}

// openInputFile either open a file named by os.Args[1],
// and return an *os.File, or if that command line argument doesn't
// exist, return os.Stdin. Also return a closing function, which can
// always be called, even if openInputFile returns os.Stdin
func openInputFile() (*os.File, func(), error) {
	var closeFunc = func() {}
	fin := os.Stdin
	if len(os.Args) > 1 {
		var err error
		if fin, err = os.Open(os.Args[1]); err != nil {
			return nil, closeFunc, err
		}
		closeFunc = func() { fin.Close() }
	}
	return fin, closeFunc, nil
}

// example function called on each input line, with input line
// as formal argument.
func timestampParser(text string) error {
	timestamp, err := time.Parse("02/Jan/2006:15:04:05", text)
	if err != nil {
		return fmt.Errorf("%q unparseable\n", text)
	}
	fmt.Println(timestamp.Format(time.RFC3339))
	return nil
}
