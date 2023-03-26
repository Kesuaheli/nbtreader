package nbt

import (
	"bytes"
	"encoding/binary"
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
	StringVal() string
	parse([]byte, bool) (NbtTag, []byte, error)
}
type EndTag struct {
}
type Byte struct {
	Name  string
	Value int8
}
type Short struct {
	Name  string
	Value int16
}
type Int struct {
	Name  string
	Value int32
}
type Long struct {
	Name  string
	Value int64
}
type Float struct {
	Name  string
	Value float32
}
type Double struct {
	Name  string
	Value float64
}
type ByteArray struct {
	Name  string
	Items []int8
}
type String struct {
	Name  string
	Value string
}
type List struct {
	Name    string
	Type    TagType
	Entries []NbtTag
}
type Compound struct {
	Name     string
	Children []NbtTag
}
type IntArray struct {
	Name  string
	Items []int32
}
type LongArray struct {
	Name  string
	Items []int64
}

// named tags (name + value)

func (t EndTag) String() string {
	return fmt.Sprint("END_TAG")
}
func (t Byte) String() string {
	return fmt.Sprintf("%s: %s", t.Name, t.StringVal())
}
func (t Short) String() string {
	return fmt.Sprintf("%s: %s", t.Name, t.StringVal())
}
func (t Int) String() string {
	return fmt.Sprintf("%s: %s", t.Name, t.StringVal())
}
func (t Long) String() string {
	return fmt.Sprintf("%s: %s", t.Name, t.StringVal())
}
func (t Float) String() string {
	return fmt.Sprintf("%s: %s", t.Name, t.StringVal())
}
func (t Double) String() string {
	return fmt.Sprintf("%s: %s", t.Name, t.StringVal())
}
func (t ByteArray) String() string {
	return fmt.Sprintf("%s: %s", t.Name, t.StringVal())
}
func (t String) String() string {
	return fmt.Sprintf("%s: %s", t.Name, t.StringVal())
}
func (t List) String() string {
	return fmt.Sprintf("%s: %s", t.Name, t.StringVal())
}
func (t Compound) String() string {
	return fmt.Sprintf("%s: %s", t.Name, t.StringVal())
}
func (t IntArray) String() string {
	return fmt.Sprintf("%s: %s", t.Name, t.StringVal())
}
func (t LongArray) String() string {
	return fmt.Sprintf("%s: %s", t.Name, t.StringVal())
}

// unnamed tags (only values)
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

func (t EndTag) StringVal() string {
	return fmt.Sprint("END_TAG")
}
func (t Byte) StringVal() string {
	return fmt.Sprintf("%db", t.Value)
}
func (t Short) StringVal() string {
	return fmt.Sprintf("%ds", t.Value)
}
func (t Int) StringVal() string {
	return fmt.Sprintf("%d", t.Value)
}
func (t Long) StringVal() string {
	return fmt.Sprintf("%dl", t.Value)
}
func (t Float) StringVal() string {
	return fmt.Sprintf("%gf", t.Value)
}
func (t Double) StringVal() string {
	return fmt.Sprintf("%gd", t.Value)
}
func (t ByteArray) StringVal() string {
	var itemsString string
	for i, item := range t.Items {
		if i == 0 {
			itemsString = fmt.Sprint(item)
		} else {
			itemsString = itemsString + ", " + fmt.Sprint(item)
		}
	}
	return fmt.Sprintf("[B; %s]", itemsString)
}
func (t String) StringVal() string {
	return fmt.Sprintf("\"%s\"", t.Value)
}
func (t List) StringVal() string {
	var entriesString string
	for i, entry := range t.Entries {
		if i == 0 {
			indentIncr()
			entriesString = entry.StringVal()
		} else {
			entriesString = entriesString + ", " + entry.StringVal()
		}
		if i == len(t.Entries)-1 {
			indentDecr()
			// entriesString = entriesString
		}
	}
	return fmt.Sprintf("[%s]", entriesString)
}
func (t Compound) StringVal() string {
	var childsString string
	for i, child := range t.Children {
		if i == 0 {
			indentIncr()
			childsString = indent() + child.String()
		} else {
			childsString = childsString + "," + indent() + child.String()
		}
		if i == len(t.Children)-1 {
			indentDecr()
			childsString = childsString + indent()
		}
	}
	return fmt.Sprintf("{%s}", childsString)
}
func (t IntArray) StringVal() string {
	var itemsString string
	for i, item := range t.Items {
		if i == 0 {
			itemsString = fmt.Sprint(item)
		} else {
			itemsString = itemsString + ", " + fmt.Sprint(item)
		}
	}
	return fmt.Sprintf("[I; %s]", itemsString)
}
func (t LongArray) StringVal() string {
	var itemsString string
	for i, item := range t.Items {
		if i == 0 {
			itemsString = fmt.Sprint(item)
		} else {
			itemsString = itemsString + ", " + fmt.Sprint(item)
		}
	}
	return fmt.Sprintf("[L; %s]", itemsString)
}

// parsing

func parseType(b []byte, named bool, tagType TagType) (NbtTag, []byte, error) {
	var tag NbtTag
	switch tagType {
	case Tag_End:
		tag = EndTag{}
	case Tag_Byte:
		tag = Byte{}
	case Tag_Short:
		tag = Short{}
	case Tag_Int:
		tag = Int{}
	case Tag_Long:
		tag = Long{}
	case Tag_Float:
		tag = Float{}
	case Tag_Double:
		tag = Double{}
	case Tag_Byte_Array:
		tag = ByteArray{}
	case Tag_String:
		tag = String{}
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
	tag, b, err := tag.parse(b, named)
	return tag, b, err
}

func (t EndTag) parse(b []byte, named bool) (NbtTag, []byte, error) {
	return t, b, nil
}

func (t Byte) parse(b []byte, named bool) (NbtTag, []byte, error) {
	var err error
	if named {
		t.Name, b, err = popName(b)
		if err != nil {
			return t, b, err
		}
	}
	t.Value, b, err = popByte(b)
	return t, b, err
}

func (t Short) parse(b []byte, named bool) (NbtTag, []byte, error) {
	var err error
	if named {
		t.Name, b, err = popName(b)
		if err != nil {
			return t, b, err
		}
	}
	t.Value, b, err = popShort(b)
	return t, b, err
}

func (t Int) parse(b []byte, named bool) (NbtTag, []byte, error) {
	var err error
	if named {
		t.Name, b, err = popName(b)
		if err != nil {
			return t, b, err
		}
	}
	t.Value, b, err = popInt(b)
	return t, b, err
}

func (t Long) parse(b []byte, named bool) (NbtTag, []byte, error) {
	var err error
	if named {
		t.Name, b, err = popName(b)
		if err != nil {
			return t, b, err
		}
	}
	t.Value, b, err = popLong(b)
	return t, b, err
}

func (t Float) parse(b []byte, named bool) (NbtTag, []byte, error) {
	var err error
	if named {
		t.Name, b, err = popName(b)
		if err != nil {
			return t, b, err
		}
	}
	t.Value, b, err = popFloat(b)
	return t, b, err
}

func (t Double) parse(b []byte, named bool) (NbtTag, []byte, error) {
	var err error
	if named {
		t.Name, b, err = popName(b)
		if err != nil {
			return t, b, err
		}
	}
	t.Value, b, err = popDouble(b)
	return t, b, err
}

func (t ByteArray) parse(b []byte, named bool) (NbtTag, []byte, error) {
	var err error
	if named {
		t.Name, b, err = popName(b)
		if err != nil {
			return t, b, err
		}
	}
	itemCap, b, err := popInt(b)
	if err != nil {
		return t, b, err
	}
	t.Items = make([]int8, 0, itemCap)
	for i := 0; i < int(itemCap); i++ {
		var item int8
		item, b, err = popByte(b)
		if err != nil {
			return t, b, err
		}
		t.Items = append(t.Items, item)
	}
	return t, b, nil
}

func (t String) parse(b []byte, named bool) (NbtTag, []byte, error) {
	var err error
	if named {
		t.Name, b, err = popName(b)
		if err != nil {
			return t, b, err
		}
	}
	t.Value, b, err = popName(b)
	return t, b, err
}

func (t List) parse(b []byte, named bool) (NbtTag, []byte, error) {
	var err error
	if named {
		t.Name, b, err = popName(b)
		if err != nil {
			return t, b, err
		}
	}

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

	t.Entries = make([]NbtTag, 0, itemCap)
	for i := 0; i < int(itemCap); i++ {
		var entry NbtTag
		entry, b, err = parseType(b, false, tagType)
		if err != nil {
			return t, b, err
		}
		t.Entries = append(t.Entries, entry)
	}
	return t, b, nil
}

func (t Compound) parse(b []byte, named bool) (NbtTag, []byte, error) {
	var err error
	if named {
		t.Name, b, err = popName(b)
		if err != nil {
			return t, b, err
		}
	}

	t.Children = make([]NbtTag, 0)
	for {
		var i int8
		i, b, err = popByte(b)
		if err != nil {
			return t, b, err
		}
		tagType := TagType(i)

		if tagType == Tag_End {
			break
		}

		var child NbtTag
		child, b, err = parseType(b, true, tagType)
		if err != nil {
			return t, b, err
		}
		t.Children = append(t.Children, child)
	}
	return t, b, nil
}

func (t IntArray) parse(b []byte, named bool) (NbtTag, []byte, error) {
	var err error
	if named {
		t.Name, b, err = popName(b)
		if err != nil {
			return t, b, err
		}
	}
	itemCap, b, err := popInt(b)
	if err != nil {
		return t, b, err
	}
	t.Items = make([]int32, 0, itemCap)
	for i := 0; i < int(itemCap); i++ {
		var item int32
		item, b, err = popInt(b)
		if err != nil {
			return t, b, err
		}
		t.Items = append(t.Items, item)
	}
	return t, b, nil
}

func (t LongArray) parse(b []byte, named bool) (NbtTag, []byte, error) {
	var err error
	if named {
		t.Name, b, err = popName(b)
		if err != nil {
			return t, b, err
		}
	}
	itemCap, b, err := popInt(b)
	if err != nil {
		return t, b, err
	}
	t.Items = make([]int64, 0, itemCap)
	for i := 0; i < int(itemCap); i++ {
		var item int64
		item, b, err = popLong(b)
		if err != nil {
			return t, b, err
		}
		t.Items = append(t.Items, item)
	}
	return t, b, nil
}

func popByte(b []byte) (int8, []byte, error) {
	return int8(b[0]), b[1:], nil
}

func popShort(b []byte) (int16, []byte, error) {
	return int16(binary.BigEndian.Uint16(b[:2])), b[2:], nil
}

func popInt(b []byte) (int32, []byte, error) {
	return int32(binary.BigEndian.Uint32(b[:4])), b[4:], nil
}

func popLong(b []byte) (int64, []byte, error) {
	return int64(binary.BigEndian.Uint64(b[:8])), b[8:], nil
}

func popFloat(b []byte) (float32, []byte, error) {
	if len(b) < 4 {
		return 0, b, fmt.Errorf("tried to get float (4 bytes) from %d bytes of data", len(b))
	}
	var f float32
	buf := bytes.NewReader(b[:4])
	err := binary.Read(buf, binary.BigEndian, &f)

	return f, b[4:], err
}

func popDouble(b []byte) (float64, []byte, error) {
	if len(b) < 8 {
		return 0, b, fmt.Errorf("tried to get double (8 bytes) from %d bytes of data", len(b))
	}
	var f float64
	buf := bytes.NewReader(b[:8])
	err := binary.Read(buf, binary.BigEndian, &f)

	return f, b[8:], err
}

func popString(b []byte, n uint16) (string, []byte, error) {
	if len(b) < int(n) {
		return "", b, fmt.Errorf("tried to get %d byte string from %d bytes of data", n, len(b))
	}

	return string(b[:n]), b[n:], nil
}

func popName(b []byte) (string, []byte, error) {
	lenName, b, err := popShort(b)
	if err != nil {
		return "", b, err
	}
	return popString(b, uint16(lenName))
}
