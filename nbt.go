package nbtreader

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
)

// NBT is a go struct representation of a Minecraft NBT object.
type NBT struct {
	rw *bufio.ReadWriter
	w  io.Writer

	rootName String
	root     NbtTag
}

// New creates a new NBT object. The given data will be completely parsed, including decompression
// (if compressed).
//
// The resulting NBT object can be used to change or get single nbt values and compose it again.
func New(r io.Reader, w io.Writer) (nbt *NBT, err error) {
	nbt = &NBT{
		w: w,
		rw: bufio.NewReadWriter(
			bufio.NewReader(r),
			bufio.NewWriter(w),
		),
	}

	err = nbt.parse()
	return nbt, err
}

// String implements the fmt.Stringer interface. The given NBT object will be converted to a SNBT
// string, including linebreaks.
func (nbt NBT) String() string {
	return fmt.Sprint(nbt.root)
}

func (nbt *NBT) parse() error {
	err := nbt.decompress()
	if err != nil {
		return fmt.Errorf("nbt: %v", err)
	}

	var rootType TagType
	rootType, err = popType(nbt.rw)
	if err != nil {
		return err
	}

	switch rootType {
	case Tag_Compound:
		nbt.root = Compound{}
	case Tag_List:
		nbt.root = List{}
	default:
		return fmt.Errorf("nbt: found invalid root tag: %s", rootType)
	}

	nbt.rootName, err = popString(nbt.rw)
	if err != nil {
		return fmt.Errorf("nbt: %v", err)
	}
	nbt.root, err = nbt.root.parse(nbt.rw)
	if err != nil {
		return err
	}

	// TODO: check for rest data in reader
	/* Outdated code
	if len(nbt.buf) > 0 {
		return fmt.Errorf("nbt: has %d bytes of rest data after parsing:\n% 02x", len(nbt.buf), nbt.buf)
	}
	*/
	return nil
}

// TODO: update Compose to io.Writer interface
// Compose takes the NBT object and composes it to the nbt binary format. It also will be compressed
// using gzip. The data can be accessed using nbt.Data. Compose conveniently returns this Data
// already.
func (nbt *NBT) NBT(compressed bool) error {
	if compressed {
		gzipWriter := gzip.NewWriter(nbt.rw)
		nbt.rw.Writer = bufio.NewWriter(gzipWriter)
	}

	if err := pushByte(nbt.rw, nbt.root.Type()); err != nil {
		return err
	}
	if err := pushString(nbt.rw, nbt.rootName); err != nil {
		return err
	}

	if err := nbt.root.compose(nbt.rw); err != nil {
		return err
	}
	return nbt.rw.Flush()
}

type compression = byte

const (
	NONE compression = iota
	GZIP
	ZIP
	TAR
)

func (nbt *NBT) decompress() error {
	switch c := nbt.getCompressionType(); c {
	case NONE:
		return nil
	case GZIP:
		gzipReader, err := gzip.NewReader(nbt.rw.Reader)
		if err != nil {
			return err
		}
		nbt.rw.Reader = bufio.NewReader(gzipReader)
		return nil
	case ZIP:
		return fmt.Errorf("file has ZIP compression: ZIP is not supportet yet")
	case TAR:
		return fmt.Errorf("file has TAR compression: TAR is not supportet yet")
	default:
		return fmt.Errorf("file has unsupported compression: %2x", c)
	}
}

func (nbt *NBT) MarshalJSON() ([]byte, error) {
	return json.Marshal(nbt.root)
}

func (nbt *NBT) MarshalNJSON() ([]byte, error) {
	return MarshalNJSON(nbt.root)
}

// TODO: update compress to use io.Writer interface
/* Outdated code
func (nbt *NBT) compress() error {
	buf := bytes.NewBuffer([]byte{})
	gz := gzip.NewWriter(buf)
	_, err := gz.Write(nbt.buf)
	gz.Close()
	nbt.buf = buf.Bytes()
	return err
}
*/

func (nbt NBT) getCompressionType() compression {
	buf, err := nbt.rw.Peek(4)
	if err != nil {
		panic(err)
	}

	switch {
	case bytes.Equal(buf[:3], []byte{0x1f, 0x8b, 0x08}):
		return GZIP
	case bytes.Equal(buf[:4], []byte{0x50, 0x4b, 0x03, 0x04}):
		return ZIP
	case bytes.Equal(buf[:4], []byte{0x75, 0x73, 0x74, 0x61}):
		return TAR
	default:
		return NONE
	}
}
