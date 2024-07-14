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
	inputType    *string
	output       *string
	outputType   *string
	uncompressed *bool
)

func init() {
	inputType = flag.String("inType", fileTypeNBT, "The filetype of input file.")
	output = flag.String("out", "", "The file to write the output to. If ommitted, output is written to stdout.")
	outputType = flag.String("outType", fileTypeSNBT, "The filetype of output file.")
	uncompressed = flag.Bool("uncompressed", false, "If the output NBT data should be raw. Otherwise using GZip compression.")
}

func main() {
	flag.Parse()

	*inputType = strings.ToLower(*inputType)
	*outputType = strings.ToLower(*outputType)

	if *inputType != fileTypeNBT {
		exitUsage(fmt.Errorf("unknown or unsupported input type '%s'", *inputType))
	}

	var (
		inFile  *os.File
		outFile *os.File
		err     error
	)

	if *output == "" {
		outFile = os.Stdout
	} else {
		outFile, err = os.Create(*output)
		if err != nil {
			exitUsage(err)
		}
		defer outFile.Close()
	}
	if flag.Arg(0) == "" {
		inFile = os.Stdin
	} else if flag.Arg(0) == *output {
		exitUsage(fmt.Errorf("flag '-out': Writing the output to the same file as reading from is not supportet.\nConsider using a temporarily file and rename is afterwards."))
	} else {
		inFile, err = os.Open(flag.Arg(0))
		if err != nil {
			exitUsage(err)
		}
		defer inFile.Close()
	}

	nbt, err := nbtreader.New(inFile, outFile)
	if err != nil {
		fmt.Println("Error while reading file:")
		exitUsage(err)
	}

	var out []byte
	switch *outputType {
	case fileTypeJSON:
		out, err = json.MarshalIndent(nbt, "", "	")
		out = append(out, '\n')
	case fileTypeNBT:
		err = nbt.NBT(!*uncompressed)
		if err != nil {
			exitUsage(err)
		}
		return
	case fileTypeNJSON:
		out, err = nbtreader.MarshalNJSON(nbt)
		out = append(out, '\n')
	case fileTypeSNBT:
		out = []byte(fmt.Sprint(nbt))
		out = append(out, '\n')
	default:
		exitUsage(fmt.Errorf("unknown or unsupported output type '%s'", *outputType))
	}

	if err != nil {
		fmt.Printf("Error while marshalling to output '%s':\n", *outputType)
		exitUsage(err)
	}

	_, err = outFile.Write(out)
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
