package nbtreader

import (
	"encoding/binary"
	"fmt"
	"math"
)

func parseType(b []byte, tagType TagType) (NbtTag, []byte, error) {
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
		return tag, b, fmt.Errorf("unknown type %02x, next 10 bytes: % 02x", tagType, b[:10])
	}
	tag, b, err := tag.parse(b)
	return tag, b, err
}

func (t Byte) parse(b []byte) (NbtTag, []byte, error) {
	return popByte(b)
}

func (t Short) parse(b []byte) (NbtTag, []byte, error) {
	return popShort(b)
}

func (t Int) parse(b []byte) (NbtTag, []byte, error) {
	return popInt(b)
}

func (t Long) parse(b []byte) (NbtTag, []byte, error) {
	return popLong(b)
}

func (t Float) parse(b []byte) (NbtTag, []byte, error) {
	return popFloat(b)
}

func (t Double) parse(b []byte) (NbtTag, []byte, error) {
	return popDouble(b)
}

func (t ByteArray) parse(b []byte) (NbtTag, []byte, error) {
	itemCap, b, err := popInt(b)
	if err != nil {
		return t, b, err
	}
	if int(itemCap) > len(b) {
		return t, b, fmt.Errorf("tried to parse %d bytes but only %d left:\n%v\n%s", itemCap, len(b), b, b)
	}
	t = make([]Byte, 0, itemCap)
	for i := 0; i < int(itemCap); i++ {
		var item Byte
		item, b, err = popByte(b)
		if err != nil {
			return t, b, err
		}
		t = append(t, item)
	}
	return t, b, nil
}

func (t String) parse(b []byte) (NbtTag, []byte, error) {
	return popString(b)
}

func (t List) parse(b []byte) (NbtTag, []byte, error) {
	i, b, err := popByte(b)
	if err != nil {
		return t, b, err
	}
	tagType := TagType(i)
	itemCap, b, err := popInt(b)
	if err != nil {
		return t, b, err
	}
	if tagType == Tag_End && itemCap > 0 {
		return t, b, fmt.Errorf("list cannot be of type TAG_END")
	}

	t = make([]NbtTag, 0, itemCap)
	for i := 0; i < int(itemCap); i++ {
		var entry NbtTag
		entry, b, err = parseType(b, tagType)
		if err != nil {
			return t, b, err
		}
		t = append(t, entry)
	}
	return t, b, nil
}

func (t Compound) parse(b []byte) (NbtTag, []byte, error) {
	for {
		var i Byte
		var err error
		i, b, err = popByte(b)
		if err != nil {
			return t, b, err
		}
		tagType := TagType(i)

		if tagType == Tag_End {
			return t, b, nil
		}

		var key String
		var child NbtTag
		key, b, err = popString(b)
		child, b, err = parseType(b, tagType)
		if err != nil {
			return t, b, err
		}
		t[key] = child
	}
}

func (t IntArray) parse(b []byte) (NbtTag, []byte, error) {
	itemCap, b, err := popInt(b)
	if err != nil {
		return t, b, err
	}
	t = make([]Int, 0, itemCap)
	for i := 0; i < int(itemCap); i++ {
		var item Int
		item, b, err = popInt(b)
		if err != nil {
			return t, b, err
		}
		t = append(t, item)
	}
	return t, b, nil
}

func (t LongArray) parse(b []byte) (NbtTag, []byte, error) {
	itemCap, b, err := popInt(b)
	if err != nil {
		return t, b, err
	}
	t = make([]Long, 0, itemCap)
	for i := 0; i < int(itemCap); i++ {
		var item Long
		item, b, err = popLong(b)
		if err != nil {
			return t, b, err
		}
		t = append(t, item)
	}
	return t, b, nil
}

func popByte(b []byte) (Byte, []byte, error) {
	return Byte(b[0]), b[1:], nil
}

func popShort(b []byte) (Short, []byte, error) {
	return Short(binary.BigEndian.Uint16(b[:2])), b[2:], nil
}

func popInt(b []byte) (Int, []byte, error) {
	return Int(binary.BigEndian.Uint32(b[:4])), b[4:], nil
}

func popLong(b []byte) (Long, []byte, error) {
	return Long(binary.BigEndian.Uint64(b[:8])), b[8:], nil
}

func popFloat(b []byte) (Float, []byte, error) {
	i, b, err := popInt(b)
	f := Float(math.Float32frombits(uint32(i)))
	return f, b, err
}

func popDouble(b []byte) (Double, []byte, error) {
	l, b, err := popLong(b)
	d := Double(math.Float64frombits(uint64(l)))
	return d, b, err
}

func popString(b []byte) (String, []byte, error) {
	lenName, b, err := popShort(b)
	if err != nil {
		return "", b, err
	}

	if len(b) < int(lenName) {
		return "", b, fmt.Errorf("tried to get %d byte string from %d bytes of data", lenName, len(b))
	}

	return String(b[:lenName]), b[lenName:], nil
}
