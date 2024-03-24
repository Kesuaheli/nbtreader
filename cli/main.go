package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Kesuaheli/nbtreader"
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

	nbt, err := nbtreader.New(file)
	if err != nil {
		fmt.Println("Error while reading file:")
		exitUsage(err)
	}

	fmt.Println("SNBT:\n", nbt)

	// TODO: update Compose to io.Writer interface
	/* Outdated code
	data = nbt.Compose()

	if err = os.MkdirAll("files", 0644); err != nil {
		fmt.Println("Error while crating output dir:")
		exitUsage(err)
	}
	if err = os.WriteFile("files/output.dat", data, 0644); err != nil {
		fmt.Println("Error while writing output file:")
		exitUsage(err)
	}
	fmt.Println("Wrote file to files/output.dat")
	*/
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
