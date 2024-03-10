package nbt

import (
	"encoding/binary"
	"fmt"
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
type EndTag struct {
}
type Byte byte
type Short int16
type Int int32
type Long int64
type Float float32
type Double float64
type ByteArray []Byte
type String string
type List []NbtTag
type Compound map[String]NbtTag
type IntArray []Int
type LongArray []Long

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

func (t EndTag) String() string {
	return fmt.Sprint("END_TAG")
}

func (t EndTag) Type() TagType {
	return Tag_End
}

func (t Byte) String() string {
	return fmt.Sprintf("%db", t)
}

func (t Byte) Type() TagType {
	return Tag_Byte
}

func (t Short) String() string {
	return fmt.Sprintf("%ds", t)
}

func (t Short) Type() TagType {
	return Tag_Short
}

func (t Int) String() string {
	return fmt.Sprintf("%d", t)
}

func (t Int) Type() TagType {
	return Tag_Int
}

func (t Long) String() string {
	return fmt.Sprintf("%dl", t)
}

func (t Long) Type() TagType {
	return Tag_Long
}

func (t Float) String() string {
	return fmt.Sprintf("%gf", t)
}

func (t Float) Type() TagType {
	return Tag_Float
}

func (t Double) String() string {
	return fmt.Sprintf("%gd", t)
}

func (t Double) Type() TagType {
	return Tag_Double
}

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

func (t String) String() string {
	return "\"" + string(t) + "\""
}

func (t String) Type() TagType {
	return Tag_String
}

func (t String) Len() Short {
	return Short(len(t))
}

func (t List) String() string {
	var entriesString string
	for i, entry := range t {
		if i == 0 {
			indentIncr()
			entriesString = entry.String()
		} else {
			entriesString = entriesString + ", " + entry.String()
		}
		if i == len(t)-1 {
			indentDecr()
			// entriesString = entriesString
		}
	}
	return fmt.Sprintf("[%s]", entriesString)
}

func (t List) Type() TagType {
	return Tag_List
}

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

// parsing

func parseType(b []byte, tagType TagType) (NbtTag, []byte, error) {
	var tag NbtTag
	switch tagType {
	case Tag_End:
		tag = EndTag{}
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
		tag = EndTag{}
	}
	tag, b, err := tag.parse(b)
	return tag, b, err
}

func (t EndTag) parse(b []byte) (NbtTag, []byte, error) {
	fmt.Printf("parse tag end? %d bytes left: next 6: %v (%s)\n", len(b), b[:6], b[:6])
	return t, b, nil
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

func pushByte[B Byte | byte | int8 | TagType](data []byte, b B) []byte {
	return append(data, byte(b))
}

func popShort(b []byte) (Short, []byte, error) {
	return Short(binary.BigEndian.Uint16(b[:2])), b[2:], nil
}

func pushShort[S Short | int16 | uint16](data []byte, s S) []byte {
	return binary.BigEndian.AppendUint16(data, uint16(s))
}

func popInt(b []byte) (Int, []byte, error) {
	return Int(binary.BigEndian.Uint32(b[:4])), b[4:], nil
}

func pushInt[I Int | int32 | uint32 | int](data []byte, i I) []byte {
	return binary.BigEndian.AppendUint32(data, uint32(i))
}

func popLong(b []byte) (Long, []byte, error) {
	return Long(binary.BigEndian.Uint64(b[:8])), b[8:], nil
}

func pushLong[L Long | int64 | uint64 | int](data []byte, l L) []byte {
	return binary.BigEndian.AppendUint64(data, uint64(l))
}

func popFloat(b []byte) (Float, []byte, error) {
	i, b, err := popInt(b)
	f := Float(math.Float32frombits(uint32(i)))
	return f, b, err
}

func pushFloat[F Float | float32](data []byte, f F) []byte {
	return pushInt(data, Int(math.Float32bits(float32(f))))
}

func popDouble(b []byte) (Double, []byte, error) {
	l, b, err := popLong(b)
	d := Double(math.Float64frombits(uint64(l)))
	return d, b, err
}

func pushDouble[D Double | float64](data []byte, d D) []byte {
	return pushLong(data, Long(math.Float64bits(float64(d))))
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

func pushString(b []byte, s String) []byte {
	b = pushShort(b, s.Len())
	return append(b, []byte(s)...)
}
