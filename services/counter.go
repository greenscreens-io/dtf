package services

import (
	"fmt"
	"strings"

	"github.com/dustin/go-humanize"
)

// HTTPCounter counts the number of bytes written to it. It implements to the io.Writer interface
// and we can pass this into io.TeeReader() which will report progress on each write cycle.
type HTTPCounter struct {
	Total      uint64
	transfered uint64
}

// Write writer to count written bytes
func (wc *HTTPCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.transfered += uint64(n)
	wc.PrintProgress()
	return n, nil
}

// Read writer to count written bytes
func (wc *HTTPCounter) Read(p []byte) (int, error) {
	n := len(p)
	wc.transfered += uint64(n)
	wc.PrintProgress()
	return n, nil
}

// PrintProgress to output current size
func (wc *HTTPCounter) PrintProgress() {
	// Clear the line by using a character return to go back to the start and remove
	// the remaining characters by filling it with spaces
	fmt.Printf("\r%s", strings.Repeat(" ", 40))

	// Return again and print current status of download
	// We use the humanize package to print the bytes in a meaningful way (e.g. 10 MB)
	if wc.Total == 0 {
		fmt.Printf("\rProgress... %s complete", humanize.Bytes(wc.transfered))
	}

	if wc.Total > 0 {
		fmt.Printf("\rProgress... %s complete out of %s", humanize.Bytes(wc.transfered), humanize.Bytes(wc.Total))
	}
}

// Finish counter by printing new line
func (wc *HTTPCounter) Finish() {
	fmt.Print("\n")
}
