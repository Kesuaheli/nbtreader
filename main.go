package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"os"

	"nbtreader/mca"
	"nbtreader/nbt"
	"path/filepath"
)

type Compression int

const (
	NONE Compression = iota
	GZIP
	ZIP
	TAR
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

	dataRaw, err := decompress(file)
	if err != nil {
		fmt.Println("Error while reading file:")
		exitUsage(err)
	}

	if match := mca.GetRegex().FindStringSubmatch(os.Args[1]); match != nil {
		region, err := mca.ParseRegex(match)
		if err != nil {
			fmt.Println("Error parsing regex:")
			exitUsage(err)
		}

		err = region.ParseData(dataRaw)
		if err != nil {
			fmt.Println("Error while parsing region:")
			exitUsage(err)
		}
		fmt.Println(region)
	} else {
		nbtRoot, restData, err := nbt.NewParser(dataRaw)
		if err != nil {
			fmt.Println("Error while parsing data:")
			exitUsage(err)
		}

		fmt.Println("{", nbtRoot, "}")
		if len(restData) > 0 {
			fmt.Println("WARNING: rest data:", restData)
		}
	}
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

func decompress(f *os.File) ([]byte, error) {
	b, err := io.ReadAll(f)
	if err != nil {
		return []byte{}, err
	}

	buf := bytes.NewBuffer(b)

	switch getCompressionType(b) {
	case NONE:
		return b, nil
	case GZIP:
		gz, err := gzip.NewReader(buf)
		if err != nil {
			return []byte{}, err
		}
		defer gz.Close()
		return io.ReadAll(gz)
	case ZIP:
		return []byte{}, fmt.Errorf("file has ZIP compression: ZIP is not supportet yet")
	case TAR:
		return []byte{}, fmt.Errorf("file has TAR compression: TAR is not supportet yet")
	default:
		return []byte{}, fmt.Errorf("file has unsupported compression")
	}
}

func getCompressionType(b []byte) Compression {
	switch {
	case bytes.Equal(b[:3], []byte{0x1f, 0x8b, 0x08}):
		return GZIP
	case bytes.Equal(b[:4], []byte{0x50, 0x4b, 0x03, 0x04}):
		return ZIP
	case bytes.Equal(b[:4], []byte{0x75, 0x73, 0x74, 0x61}):
		return TAR
	default:
		return NONE
	}
}
