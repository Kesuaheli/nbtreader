package main

import (
	"encoding/json"
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

	out, err := json.MarshalIndent(nbt, "", "	")
	if err != nil {
		fmt.Println("Error while marshalling json:")
		exitUsage(err)
	}

	out2, err := nbtreader.MarshalNJSON(nbt)
	if err != nil {
		fmt.Println("Error while marshalling njson:")
		exitUsage(err)
	}

	err = os.WriteFile("../files/out.json", out, 0664)
	if err != nil {
		fmt.Println("Error while writing json file:")
		exitUsage(err)
	}
	err = os.WriteFile("../files/out2.json", out2, 0664)
	if err != nil {
		fmt.Println("Error while writing njson file:")
		exitUsage(err)
	}

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
