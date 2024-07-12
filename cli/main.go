package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/Kesuaheli/nbtreader"
)

type FileType string

const (
	fileTypeJSON  = "json"
	fileTypeNBT   = "nbt"
	fileTypeNJSON = "njson"
	fileTypeSNBT  = "snbt"
)

var (
	input      *string
	inputType  *string
	output     *string
	outputType *string
)

func init() {
	input = flag.String("file", "", "The input file to read. If ommitted, file is read from stdin.")
	inputType = flag.String("inType", fileTypeNBT, "The filetype of input file.")
	output = flag.String("out", "", "The file to write the output to. If ommitted, output is written to stdout.")
	outputType = flag.String("outType", fileTypeSNBT, "The filetype of output file.")
}

func main() {
	flag.Parse()

	*inputType = strings.ToLower(*inputType)
	*outputType = strings.ToLower(*outputType)

	if *inputType != fileTypeNBT {
		exitUsage(fmt.Errorf("unknown or unsupported input type '%s'", *inputType))
	}
	if *outputType != fileTypeJSON && *outputType != fileTypeNJSON && *outputType != fileTypeSNBT {
		exitUsage(fmt.Errorf("unknown or unsupported output type '%s'", *inputType))
	}

	var file *os.File
	var err error
	if *input == "" {
		file = os.Stdin
	} else {
		file, err = os.Open(*input)
		if err != nil {
			exitUsage(err)
		}
	}

	nbt, err := nbtreader.New(file)
	if err != nil {
		fmt.Println("Error while reading file:")
		exitUsage(err)
	}

	var out []byte
	switch *outputType {
	case fileTypeJSON:
		out, err = json.MarshalIndent(nbt, "", "	")
	case fileTypeNJSON:
		out, err = nbtreader.MarshalNJSON(nbt)
	case fileTypeSNBT:
		out = []byte(fmt.Sprint(nbt))
		out = append(out, '\n')
	default:
		panic(fmt.Sprintf("unhandled output type '%s'", *outputType))
	}

	if err != nil {
		fmt.Printf("Error while marshalling to output '%s':\n", *outputType)
		exitUsage(err)
	}

	if *output == "" {
		fmt.Print(string(out))
		return
	}

	err = os.WriteFile(*output, out, 0664)
	if err != nil {
		fmt.Println("Error while writing output file:")
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
	flag.Usage()

	os.Exit(1)
}
