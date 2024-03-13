package nbt

import (
	"fmt"
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

func popType(b []byte) (TagType, []byte, error) {
	if len(b) == 0 {
		return 0x00, b, fmt.Errorf("tried to get nbt type from empty data")
	}
	return TagType(b[0]), b[1:], nil
}

type NbtTag interface {
	String() string
	Type() TagType
	parse([]byte) (NbtTag, []byte, error)
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

type Short int16

func (t Short) String() string {
	return fmt.Sprintf("%ds", t)
}

func (t Short) Type() TagType {
	return Tag_Short
}

type Int int32

func (t Int) String() string {
	return fmt.Sprintf("%d", t)
}

func (t Int) Type() TagType {
	return Tag_Int
}

type Long int64

func (t Long) String() string {
	return fmt.Sprintf("%dl", t)
}

func (t Long) Type() TagType {
	return Tag_Long
}

type Float float32

func (t Float) String() string {
	return fmt.Sprintf("%gf", t)
}

func (t Float) Type() TagType {
	return Tag_Float
}

type Double float64

func (t Double) String() string {
	return fmt.Sprintf("%gd", t)
}

func (t Double) Type() TagType {
	return Tag_Double
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

type List []NbtTag

func (t List) String() string {
	var entriesString string
	for i, entry := range t {
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
