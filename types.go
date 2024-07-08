package nbtreader

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"strings"
)

type TagType byte

const (
	Tag_End TagType = iota
	Tag_Byte
	Tag_Short
	Tag_Int
	Tag_Long
	Tag_Float
	Tag_Double
	Tag_Byte_Array
	Tag_String
	Tag_List
	Tag_Compound
	Tag_Int_Array
	Tag_Long_Array
)

func (t TagType) String() string {
	switch t {
	case Tag_End:
		return "End of Compound"
	case Tag_Byte:
		return "Byte (int8)"
	case Tag_Short:
		return "Short (int16)"
	case Tag_Int:
		return "Int (int32)"
	case Tag_Long:
		return "Long (int64)"
	case Tag_Float:
		return "Float (float32)"
	case Tag_Double:
		return "Double (float64)"
	case Tag_Byte_Array:
		return "ByteArray ([]int8)"
	case Tag_String:
		return "String"
	case Tag_List:
		return "List"
	case Tag_Compound:
		return "Start of Compound"
	case Tag_Int_Array:
		return "IntArray ([]int32)"
	case Tag_Long_Array:
		return "LongArray ([]int64)"
	default:
		return fmt.Sprintf("*Unknown Type %d*", t)
	}
}

func (t TagType) Annotation() TypeAnnotation {
	switch t {
	case Tag_Byte:
		return ByteAnnotation
	case Tag_Short:
		return ShortAnnotation
	case Tag_Int:
		return IntAnnotation
	case Tag_Long:
		return LongAnnotation
	case Tag_Float:
		return FloatAnnotation
	case Tag_Double:
		return DoubleAnnotation
	case Tag_Byte_Array:
		return ByteArrayAnnotation
	case Tag_String:
		return StringAnnotation
	case Tag_List:
		return NoAnnotation
	case Tag_Compound:
		return CompoundAnnotation
	case Tag_Int_Array:
		return IntArrayAnnotation
	case Tag_Long_Array:
		return LongArrayAnnotation
	default:
		// also for Tag_End and Tag_List
		return NoAnnotation
	}
}

func popType(r io.Reader) (TagType, error) {
	var ttype [1]byte
	if _, err := r.Read(ttype[:]); err != nil {
		return 0x00, fmt.Errorf("pop type: %v", err)
	}

	t := TagType(ttype[0])
	return t, nil
}

type NbtTag interface {
	String() string
	Type() TagType
	parse(io.Reader) (NbtTag, error)
	compose() []byte
}

var indention int = 0

func indent() string {
	return "\n" + strings.Repeat("  ", indention)
}
func indentIncr() {
	indention = indention + 1
}
func indentDecr() {
	indention = indention - 1
}

type Byte byte

func (t Byte) String() string {
	return fmt.Sprintf("%db", t)
}

func (t Byte) Type() TagType {
	return Tag_Byte
}

func (t Byte) MarshalJSON() ([]byte, error) {
	return numberToJSON(t)
}

type Short int16

func (t Short) String() string {
	return fmt.Sprintf("%ds", t)
}

func (t Short) Type() TagType {
	return Tag_Short
}

func (t Short) MarshalJSON() ([]byte, error) {
	return numberToJSON(t)
}

type Int int32

func (t Int) String() string {
	return fmt.Sprintf("%d", t)
}

func (t Int) Type() TagType {
	return Tag_Int
}

func (t Int) MarshalJSON() ([]byte, error) {
	return numberToJSON(t)
}

type Long int64

func (t Long) String() string {
	return fmt.Sprintf("%dl", t)
}

func (t Long) Type() TagType {
	return Tag_Long
}

func (t Long) MarshalJSON() ([]byte, error) {
	return numberToJSON(t)
}

type Float float32

func (t Float) String() string {
	return fmt.Sprintf("%gf", t)
}

func (t Float) Type() TagType {
	return Tag_Float
}

func (t Float) MarshalJSON() ([]byte, error) {
	return numberToJSON(t)
}

type Double float64

func (t Double) String() string {
	return fmt.Sprintf("%gd", t)
}

func (t Double) Type() TagType {
	return Tag_Double
}

func (t Double) MarshalJSON() ([]byte, error) {
	return numberToJSON(t)
}

type ByteArray []Byte

func (t ByteArray) String() string {
	var itemsString string
	for i, item := range t {
		if i == 0 {
			itemsString = fmt.Sprint(item)
		} else {
			itemsString = itemsString + ", " + fmt.Sprint(item)
		}
	}
	return fmt.Sprintf("[B; %s]", itemsString)
}

func (t ByteArray) Type() TagType {
	return Tag_Byte_Array
}

func (t ByteArray) MarshalJSON() ([]byte, error) {
	return arrayMarshalJSON([]Byte(t))
}

type String string

func (t String) String() string {
	return "\"" + string(t) + "\""
}

func (t String) Type() TagType {
	return Tag_String
}

func (t String) Len() Short {
	return Short(len(t))
}

func (t String) MarshalJSON() ([]byte, error) {
	return []byte(t.String()), nil
}

type List struct {
	TagType  TagType
	Elements []NbtTag
}

func (t List) String() string {
	var entriesString string
	for i, entry := range t.Elements {
		if i == 0 {
			entriesString = entry.String()
		} else {
			entriesString = entriesString + ", " + entry.String()
		}
	}
	return fmt.Sprintf("[%s]", entriesString)
}

func (t List) Type() TagType {
	return Tag_List
}

func (t List) MarshalJSON() ([]byte, error) {
	return arrayMarshalJSON(t.Elements)
}

type Compound map[String]NbtTag

func (t Compound) String() string {
	var childsString string
	i := 0
	for key, child := range t {
		if i == 0 {
			indentIncr()
		} else {
			childsString = childsString + ","
		}
		childsString += fmt.Sprintf("%s%s: %s", indent(), string(key), child.String())
		if i == len(t)-1 {
			indentDecr()
			childsString = childsString + indent()
		}
		i++
	}
	return fmt.Sprintf("{%s}", childsString)
}

func (t Compound) Type() TagType {
	return Tag_Compound
}

type IntArray []Int

func (t IntArray) String() string {
	var itemsString string
	for i, item := range t {
		if i == 0 {
			itemsString = fmt.Sprint(item)
		} else {
			itemsString = itemsString + ", " + fmt.Sprint(item)
		}
	}
	return fmt.Sprintf("[I; %s]", itemsString)
}

func (t IntArray) Type() TagType {
	return Tag_Int_Array
}

func (t IntArray) MarshalJSON() ([]byte, error) {
	return arrayMarshalJSON([]Int(t))
}

type LongArray []Long

func (t LongArray) String() string {
	var itemsString string
	for i, item := range t {
		if i == 0 {
			itemsString = fmt.Sprint(item)
		} else {
			itemsString = itemsString + ", " + fmt.Sprint(item)
		}
	}
	return fmt.Sprintf("[L; %s]", itemsString)
}

func (t LongArray) Type() TagType {
	return Tag_Long_Array
}

func (t LongArray) MarshalJSON() ([]byte, error) {
	return arrayMarshalJSON([]Long(t))
}

// numberToJSON is a helper function that formats any NBT number to a valid JSON number.
func numberToJSON[N Byte | Short | Int | Long | Float | Double](num N) ([]byte, error) {
	_, okF := any(num).(Float)
	_, okD := any(num).(Double)
	if !okF && !okD {
		return []byte(fmt.Sprintf("%d", int64(num))), nil
	}
	if math.Mod(float64(num), 1) == 0 {
		return []byte(fmt.Sprintf("%d.0", int64(num))), nil
	}
	return []byte(fmt.Sprintf("%f", float64(num))), nil
}

// arrayToJSON is a helper function that formats any NBT array or list to a valid JSON array.
func arrayMarshalJSON[S []E, E NbtTag](s S) ([]byte, error) {
	var (
		childsBytes bytes.Buffer
		err         error
	)
	childsBytes.WriteByte('[')
	for i, entry := range s {
		var childJSON []byte
		childJSON, err = json.Marshal(entry)
		childsBytes.Write(childJSON)
		if err != nil {
			break
		}
		if i < len(s)-1 {
			childsBytes.WriteByte(',')
		}
	}
	childsBytes.WriteByte(']')
	return childsBytes.Bytes(), err
}
