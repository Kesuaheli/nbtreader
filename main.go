package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) == 0 {
		os.Exit(1)
	} else if len(os.Args) == 1 {
		fmt.Println("Please specify a file:")
		exitUsage(nil)
	}
	file, err := os.Open(os.Args[1])
	if err != nil {
		exitUsage(err)
	}
	defer file.Close()

	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		exitUsage(err)
	}
	defer gzipReader.Close()

	nbtData, err := io.ReadAll(gzipReader)
	if err != nil {
		exitUsage(err)
	}

	fmt.Println(nbtData)

}

// exitUsage prints the error, if any, and the command usage and then
// calls os.Exit(1) to exit the program
func exitUsage(err error) {
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%s <file_to_read>\n", filepath.Base(os.Args[0]))

	os.Exit(1)
}
