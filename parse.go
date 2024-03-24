package nbtreader

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
)

func parseType(r io.Reader, tagType TagType) (NbtTag, error) {
	var tag NbtTag
	switch tagType {
	case Tag_Byte:
		tag = Byte(0)
	case Tag_Short:
		tag = Short(0)
	case Tag_Int:
		tag = Int(0)
	case Tag_Long:
		tag = Long(0)
	case Tag_Float:
		tag = Float(0)
	case Tag_Double:
		tag = Double(0)
	case Tag_Byte_Array:
		tag = ByteArray{}
	case Tag_String:
		tag = String("")
	case Tag_List:
		tag = List{}
	case Tag_Compound:
		tag = Compound{}
	case Tag_Int_Array:
		tag = IntArray{}
	case Tag_Long_Array:
		tag = LongArray{}
	default:
		return nil, fmt.Errorf("unknown type %02x", tagType)
	}
	return tag.parse(r)
}

func (t Byte) parse(r io.Reader) (NbtTag, error) {
	return popByte(r)
}

func (t Short) parse(r io.Reader) (NbtTag, error) {
	return popShort(r)
}

func (t Int) parse(r io.Reader) (NbtTag, error) {
	return popInt(r)
}

func (t Long) parse(r io.Reader) (NbtTag, error) {
	return popLong(r)
}

func (t Float) parse(r io.Reader) (NbtTag, error) {
	return popFloat(r)
}

func (t Double) parse(r io.Reader) (NbtTag, error) {
	return popDouble(r)
}

func (t ByteArray) parse(r io.Reader) (NbtTag, error) {
	itemCap, err := popInt(r)
	if err != nil {
		return t, err
	}

	t = make([]Byte, itemCap)
	for i, item := range t {
		item, err = popByte(r)
		if err != nil {
			return t, err
		}
		t[i] = item
	}
	return t, nil
}

func (t String) parse(r io.Reader) (NbtTag, error) {
	return popString(r)
}

func (t List) parse(r io.Reader) (NbtTag, error) {
	i, err := popByte(r)
	if err != nil {
		return t, err
	}
	tagType := TagType(i)
	itemCap, err := popInt(r)
	if err != nil {
		return t, err
	}
	if tagType == Tag_End && itemCap > 0 {
		return t, fmt.Errorf("list cannot be of type TAG_END")
	}

	t = make([]NbtTag, itemCap)
	for i := 0; i < int(itemCap); i++ {
		var entry NbtTag
		entry, err = parseType(r, tagType)
		if err != nil {
			return t, err
		}
		t[i] = entry
	}
	return t, nil
}

func (t Compound) parse(r io.Reader) (NbtTag, error) {
	for {
		var i Byte
		var err error
		i, err = popByte(r)
		if err != nil {
			return t, err
		}
		tagType := TagType(i)

		if tagType == Tag_End {
			return t, nil
		}

		var key String
		var child NbtTag
		key, err = popString(r)
		child, err = parseType(r, tagType)
		if err != nil {
			return t, err
		}
		t[key] = child
	}
}

func (t IntArray) parse(r io.Reader) (NbtTag, error) {
	itemCap, err := popInt(r)
	if err != nil {
		return t, err
	}

	t = make([]Int, itemCap)
	for i, item := range t {
		item, err = popInt(r)
		if err != nil {
			return t, err
		}
		t[i] = item
	}
	return t, nil
}

func (t LongArray) parse(r io.Reader) (NbtTag, error) {
	itemCap, err := popInt(r)
	if err != nil {
		return t, err
	}

	t = make([]Long, itemCap)
	for i, item := range t {
		item, err = popLong(r)
		if err != nil {
			return t, err
		}
		t[i] = item
	}
	return t, nil
}

func popByte(r io.Reader) (Byte, error) {
	var buf [1]byte
	_, err := r.Read(buf[:])
	return Byte(buf[0]), err
}

func popShort(r io.Reader) (Short, error) {
	var buf [2]byte
	_, err := r.Read(buf[:])
	return Short(binary.BigEndian.Uint16(buf[:])), err
}

func popInt(r io.Reader) (Int, error) {
	var buf [4]byte
	_, err := r.Read(buf[:])
	return Int(binary.BigEndian.Uint32(buf[:])), err
}

func popLong(r io.Reader) (Long, error) {
	var buf [8]byte
	_, err := r.Read(buf[:])
	return Long(binary.BigEndian.Uint64(buf[:])), err
}

func popFloat(r io.Reader) (Float, error) {
	i, err := popInt(r)
	f := Float(math.Float32frombits(uint32(i)))
	return f, err
}

func popDouble(r io.Reader) (Double, error) {
	l, err := popLong(r)
	d := Double(math.Float64frombits(uint64(l)))
	return d, err
}

func popString(r io.Reader) (String, error) {
	lenName, err := popShort(r)
	if err != nil {
		return "", err
	}

	p := make([]byte, lenName)
	_, err = r.Read(p)
	return String(p), err
}
