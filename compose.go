package nbtreader

import (
	"encoding/binary"
	"io"
	"math"
)

func (t Byte) compose(w io.Writer) error {
	return pushByte(w, t)
}
func (t Short) compose(w io.Writer) error {
	return pushShort(w, t)
}
func (t Int) compose(w io.Writer) error {
	return pushInt(w, t)
}
func (t Long) compose(w io.Writer) error {
	return pushLong(w, t)
}
func (t Float) compose(w io.Writer) error {
	return pushFloat(w, t)
}
func (t Double) compose(w io.Writer) error {
	return pushDouble(w, t)
}
func (t ByteArray) compose(w io.Writer) error {
	if err := pushInt(w, len(t)); err != nil {
		return err
	}
	for _, b := range t {
		if err := pushByte(w, b); err != nil {
			return err
		}
	}
	return nil
}
func (t String) compose(w io.Writer) error {
	return pushString(w, t)
}
func (t List) compose(w io.Writer) error {
	itemCap := len(t.Elements)
	if itemCap == 0 {
		if err := pushByte(w, Tag_End); err != nil {
			return err
		}
		return pushInt(w, 0)
	}

	if err := pushByte(w, t.TagType); err != nil {
		return err
	}
	if err := pushInt(w, itemCap); err != nil {
		return err
	}
	for _, entry := range t.Elements {
		if err := entry.compose(w); err != nil {
			return err
		}
	}
	return nil
}
func (t Compound) compose(w io.Writer) error {
	for _, tag := range t.getOrdered() {
		if err := pushByte(w, tag.Value.Type()); err != nil {
			return err
		}
		if err := pushString(w, tag.Key); err != nil {
			return nil
		}
		if err := tag.Value.compose(w); err != nil {
			return err
		}
	}
	return pushByte(w, Tag_End)
}
func (t IntArray) compose(w io.Writer) error {
	if err := pushInt(w, len(t)); err != nil {
		return err
	}
	for _, i := range t {
		if err := pushInt(w, i); err != nil {
			return err
		}
	}
	return nil
}
func (t LongArray) compose(w io.Writer) error {
	if err := pushInt(w, len(t)); err != nil {
		return err
	}
	for _, l := range t {
		if err := pushLong(w, l); err != nil {
			return err
		}
	}
	return nil
}

func pushByte[B Byte | int8 | TagType](w io.Writer, b B) error {
	_, err := w.Write([]byte{byte(b)})
	return err
}

func pushShort[S Short | int16 | uint8](w io.Writer, s S) error {
	var buf [2]byte
	binary.BigEndian.AppendUint16(buf[:], uint16(s))
	_, err := w.Write(buf[:])
	return err
}

func pushInt[I Int | int32 | uint16 | int](w io.Writer, i I) error {
	var buf [4]byte
	binary.BigEndian.AppendUint32(buf[:], uint32(i))
	_, err := w.Write(buf[:])
	return err
}

func pushLong[L Long | int64 | uint32 | int](w io.Writer, l L) error {
	var buf [8]byte
	binary.BigEndian.AppendUint64(buf[:], uint64(l))
	_, err := w.Write(buf[:])
	return err
}

func pushFloat[F Float | float32](w io.Writer, f F) error {
	return pushInt(w, Int(math.Float32bits(float32(f))))
}

func pushDouble[D Double | float64](w io.Writer, d D) error {
	return pushLong(w, Long(math.Float64bits(float64(d))))
}

func pushString(w io.Writer, s String) error {
	if err := pushShort(w, s.Len()); err != nil {
		return err
	}
	_, err := w.Write([]byte(s))
	return err
}
