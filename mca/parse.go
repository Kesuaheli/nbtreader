package mca

import (
	"encoding/binary"
	"fmt"
	"time"
)

type Chunck struct {
	// Position of the Chunk in the region
	X, Z int
	// Timestamp of the last modification in this Chunk
	LastModified time.Time

	// Type of compression of 'Data':
	//  00 zlib
	//  01 gzip
	DataCompressionType byte
	// The compressed chunk data (compressed with 'DataCompressionType')
	Data []byte
}

func Parse(data []byte) (region []Chunck, err error) {
	// region files begin with two 4KiB tables
	if len(data) < 2*0x1000 {
		return region, fmt.Errorf("given file is too short: needs to be at least 8192 bytes but only %d where given", len(data))
	}

	region = make([]Chunck, 0, 1024)

	for z := 0; z <= 31; z++ {
		for x := 0; x <= 31; x++ {
			chunk, err := parseChunk(x, z, data)
			if err != nil {
				return region, fmt.Errorf("failed to parse chunk at [%d,%d]: %v", x, z, err)
			}
			region = append(region, chunk)
		}
	}

	return region, nil
}

func parseChunk(x, z int, data []byte) (chunk Chunck, err error) {
	chunk.X = x
	chunk.Z = z
	var (
		// chunk index in tables, calculated from x and z chunk
		// coordinates:
		//  tableIndex  = (x + z*32) * 4
		tableIndex int = (x + z*32) * 4
		// 4 byte location data of this chunk. Contains 3 bytes (big
		// endian) of offset (in 4KiB sectors) in the region data and
		// 1 byte of length (also in 4KiB) of this chunk in the region
		// data.
		locationTable []byte = data[tableIndex : tableIndex+4]
		// 4 byte big endian timestamp of last modification in this
		// chunk.
		timestampTable []byte = data[tableIndex+4096 : tableIndex+4096+4]
		// offset (in 4KiB sections) of this chunk in the region
		// file. Multiply with 4096 for offset in bytes.
		offset uint32
		// 4 byte big endian timestamp of lastmodification in this
		// chunk.
		timestamp uint32
		// length in bytes used by this chunk.
		//
		// This may not represent all bytes occupied by this chunk in
		// the whole region file. To get the occupied byte length,
		// round this number up to the next full multiple of 4096.
		chunkDataSize uint32
		// offset in bytes of this chunk in the region file.
		totalOffset int
	)

	offsetBytes := append([]byte{0x00}, locationTable[:3]...)
	offset = binary.BigEndian.Uint32(offsetBytes)
	totalOffset = int(offset) * 4096

	timestamp = binary.BigEndian.Uint32(timestampTable)
	chunk.LastModified = time.Unix(int64(timestamp), 0)

	chunkDataSize = binary.BigEndian.Uint32(data[totalOffset : totalOffset+4])
	if len(data) < totalOffset+5+int(chunkDataSize) {
		err = fmt.Errorf("tried to read data from chunk [%d,%d] from byte %d to byte %d, but only %d bytes where given",
			x, z,
			totalOffset,
			totalOffset+5+int(chunkDataSize),
			len(data),
		)
		return chunk, err
	}

	chunk.DataCompressionType = data[totalOffset+4]
	chunk.Data = data[totalOffset+5 : totalOffset+5+int(chunkDataSize)]

	return chunk, nil
}
func (c Chunck) String() string {
	var compression string
	switch c.DataCompressionType {
	case 0x01:
		compression = "gzip"
	case 0x02:
		compression = "zlib"
	default:
		compression = fmt.Sprintf("UNKNOWN COMPRESSION %2x", c.DataCompressionType)
	}

	return fmt.Sprintf("Chunk %2d, %2d\n\tModified: %s,\n\n\tData (%s) [%d bytes]\n", c.X, c.Z, c.LastModified, compression, len(c.Data))
}
