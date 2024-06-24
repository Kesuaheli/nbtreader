package nbtreader

import (
	"encoding/binary"
	"math"
)

func (t Byte) compose() []byte {
	return pushByte([]byte{}, t)
}
func (t Short) compose() []byte {
	return pushShort([]byte{}, t)
}
func (t Int) compose() []byte {
	return pushInt([]byte{}, t)
}
func (t Long) compose() []byte {
	return pushLong([]byte{}, t)
}
func (t Float) compose() []byte {
	return pushFloat([]byte{}, t)
}
func (t Double) compose() []byte {
	return pushDouble([]byte{}, t)
}
func (t ByteArray) compose() []byte {
	buf := pushInt([]byte{}, len(t))
	for _, b := range t {
		buf = pushByte(buf, b)
	}
	return buf
}
func (t String) compose() []byte {
	return pushString([]byte{}, t)
}
func (t List) compose() []byte {
	itemCap := len(t.Elements)
	if itemCap == 0 {
		buf := pushByte([]byte{}, Tag_End)
		return pushInt(buf, 0)
	}

	buf := pushByte([]byte{}, t.TagType)
	buf = pushInt(buf, itemCap)
	for _, entry := range t.Elements {
		buf = append(buf, entry.compose()...)
	}
	return buf
}
func (t Compound) compose() []byte {
	var data []byte
	for name, tag := range t {
		data = pushByte(data, tag.Type())
		data = pushString(data, name)
		data = append(data, tag.compose()...)
	}
	data = pushByte(data, Tag_End)
	return data
}
func (t IntArray) compose() []byte {
	buf := pushInt([]byte{}, len(t))
	for _, i := range t {
		buf = pushInt(buf, i)
	}
	return buf
}
func (t LongArray) compose() []byte {
	buf := pushInt([]byte{}, len(t))
	for _, l := range t {
		buf = pushLong(buf, l)
	}
	return buf
}

func pushByte[B Byte | byte | int8 | TagType](data []byte, b B) []byte {
	return append(data, byte(b))
}

func pushShort[S Short | int16 | uint16](data []byte, s S) []byte {
	return binary.BigEndian.AppendUint16(data, uint16(s))
}

func pushInt[I Int | int32 | uint32 | int](data []byte, i I) []byte {
	return binary.BigEndian.AppendUint32(data, uint32(i))
}

func pushLong[L Long | int64 | uint64 | int](data []byte, l L) []byte {
	return binary.BigEndian.AppendUint64(data, uint64(l))
}

func pushFloat[F Float | float32](data []byte, f F) []byte {
	return pushInt(data, Int(math.Float32bits(float32(f))))
}

func pushDouble[D Double | float64](data []byte, d D) []byte {
	return pushLong(data, Long(math.Float64bits(float64(d))))
}

func pushString(b []byte, s String) []byte {
	b = pushShort(b, s.Len())
	return append(b, []byte(s)...)
}
