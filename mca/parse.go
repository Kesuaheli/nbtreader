package mca

import (
	"fmt"
	"regexp"
	"strconv"
)

type MCARegion struct {
	// Position of the region in the world
	X, Z int
	// All 32 by 32 Chunks in the region
	Chunks []Chunk
}

var regex *regexp.Regexp = regexp.MustCompile(`r\.(?P<x>-?\d+)\.(?P<z>-?\d+)\.mca`)

func GetRegex() *regexp.Regexp {
	return regex
}

func ParseRegex(match []string) (region MCARegion, err error) {
	for n, name := range regex.SubexpNames() {
		if n == 0 {
			continue
		}

		i, err := strconv.Atoi(match[n])
		if err != nil {
			return MCARegion{}, err
		}

		switch name {
		case "x":
			region.X = i
		case "z":
			region.Z = i
		default:
			return MCARegion{}, fmt.Errorf("unknown group '%s' in mca regex", name)
		}
	}

	region.Chunks = make([]Chunk, 1024)

	return region, nil
}

func (r *MCARegion) ParseData(data []byte) error {
	// region files have to begin with two 4KiB tables
	if len(data) < 2*0x1000 {
		return fmt.Errorf("given file is too short: needs to be at least 8192 bytes but only %d where given", len(data))
	}

	for z := 0; z <= 31; z++ {
		for x := 0; x <= 31; x++ {
			chunk, err := r.parseChunk(x, z, data)
			if err != nil {
				return fmt.Errorf("failed to parse chunk at [%d,%d]: %v", x, z, err)
			}
			r.Chunks[x+z*32] = chunk
		}
	}

	return nil
}
