package nbt

import (
	"fmt"
)

func NewParser(data []byte) (root NbtTag, restData []byte, err error) {

	t, data, err := popType(data)
	if err != nil {
		return root, data, err
	}

	switch t {
	case Tag_Compound:
		root = Compound{}
	case Tag_List:
		root = List{}
	default:
		return root, data, fmt.Errorf("Found invalid root tag: %s", t)
	}

	rootName, data, err := popString(data)
	if err != nil {
		return root, data, err
	}
	if rootName != "" {
		fmt.Printf("dropped root name: %s\n", rootName)
	}

	return root.parse(data)
}

/*

	switch t {
	case Tag_End:

	case Tag_Byte:

	case Tag_Short:

	case Tag_Int:

	case Tag_Long:

	case Tag_Float:

	case Tag_Double:

	case Tag_Byte_Array:

	case Tag_String:

	case Tag_List:

	case Tag_Compound:

	case Tag_Int_Array:

	case Tag_Long_Array:

	default:

	}
*/
