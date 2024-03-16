package nbtreader

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
)

// NBT is a go struct representation of a Minecraft NBT object.
type NBT struct {
	Data []byte

	buf []byte

	rootName String
	root     NbtTag
}

// New creates a new NBT object. The given data will be completely parsed, including decompression
// (if compressed).
//
// The resulting NBT object can be used to change or get single nbt values and compose it again.
func New(data []byte) (nbt *NBT, err error) {
	nbt = &NBT{
		Data: data,
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
	nbt.buf = nbt.Data
	err := nbt.decompress()
	if err != nil {
		return fmt.Errorf("nbt: %v", err)
	}

	var rootType TagType
	rootType, nbt.buf, err = popType(nbt.buf)
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

	nbt.rootName, nbt.buf, err = popString(nbt.buf)
	if err != nil {
		return fmt.Errorf("nbt: %v", err)
	}

	nbt.root, nbt.buf, err = nbt.root.parse(nbt.buf)
	if err != nil {
		return err
	}
	if len(nbt.buf) > 0 {
		return fmt.Errorf("nbt: has %d bytes of rest data after parsing:\n% 02x", len(nbt.buf), nbt.buf)
	}
	return nil
}

// Compose takes the NBT object and composes it to the nbt binary format. It also will be compressed
// using gzip. The data can be accessed using nbt.Data. Compose conveniently returns this Data
// already.
func (nbt *NBT) Compose() []byte {
	nbt.buf = []byte{}
	nbt.buf = pushByte(nbt.buf, nbt.root.Type())
	nbt.buf = pushString(nbt.buf, nbt.rootName)

	nbt.buf = append(nbt.buf, nbt.root.compose()...)
	nbt.compress()

	nbt.Data = nbt.buf
	return nbt.Data
}

type compression = byte

const (
	NONE compression = iota
	GZIP
	ZIP
	TAR
)

func (nbt *NBT) decompress() error {
	buf := bytes.NewBuffer(nbt.buf)

	switch c := nbt.getCompressionType(); c {
	case NONE:
		return nil
	case GZIP:
		gz, err := gzip.NewReader(buf)
		defer gz.Close()
		if err != nil {
			return err
		}
		nbt.buf, err = io.ReadAll(gz)
		return err
	case ZIP:
		return fmt.Errorf("file has ZIP compression: ZIP is not supportet yet")
	case TAR:
		return fmt.Errorf("file has TAR compression: TAR is not supportet yet")
	default:
		return fmt.Errorf("file has unsupported compression: %2x", c)
	}
}

func (nbt *NBT) compress() error {
	buf := bytes.NewBuffer([]byte{})
	gz := gzip.NewWriter(buf)
	_, err := gz.Write(nbt.buf)
	gz.Close()
	nbt.buf = buf.Bytes()
	return err
}

func (nbt NBT) getCompressionType() compression {
	switch {
	case bytes.Equal(nbt.buf[:3], []byte{0x1f, 0x8b, 0x08}):
		return GZIP
	case bytes.Equal(nbt.buf[:4], []byte{0x50, 0x4b, 0x03, 0x04}):
		return ZIP
	case bytes.Equal(nbt.buf[:4], []byte{0x75, 0x73, 0x74, 0x61}):
		return TAR
	default:
		return NONE
	}
}