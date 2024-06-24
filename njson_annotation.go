package nbtreader

import (
	"fmt"
	"reflect"
)

type TypeAnnotation uint8

const (
	NoAnnotation TypeAnnotation = iota
	InferenceAnnotation
	CompoundAnnotation
	ByteArrayAnnotation
	IntArrayAnnotation
	LongArrayAnnotation
	InferenceArrayAnnotation
	ByteAnnotation
	ShortAnnotation
	IntAnnotation
	LongAnnotation
	FloatAnnotation
	DoubleAnnotation
	StringAnnotation
)

func (ta TypeAnnotation) Characters() string {
	switch ta {
	default:
		return ""
	case CompoundAnnotation:
		return "C"
	case ByteArrayAnnotation:
		return "B;"
	case IntArrayAnnotation:
		return "I;"
	case LongArrayAnnotation:
		return "L;"
	case InferenceArrayAnnotation:
		return ";"
	case ByteAnnotation:
		return "b"
	case ShortAnnotation:
		return "s"
	case IntAnnotation:
		return "i"
	case LongAnnotation:
		return "L"
	case FloatAnnotation:
		return "f"
	case DoubleAnnotation:
		return "d"
	case StringAnnotation:
		return "S"
	}
}

func (ta TypeAnnotation) StringSubtype(sub string) string {
	if ta == NoAnnotation {
		return ""
	}

	if sub == "" {
		return fmt.Sprintf("<%s>", ta.Characters())
	}
	return fmt.Sprintf("<%s:%s>", ta.Characters(), sub)
}

func (ta TypeAnnotation) String() string {
	return ta.StringSubtype("")
}

func AnnotationOf(v reflect.Value) TypeAnnotation {
	switch v.Kind() {
	default:
		fmt.Printf("unknown type '%s' (%s)\n", v, v.Kind())
		return InferenceAnnotation
	case reflect.Map:
		return CompoundAnnotation
	case reflect.Slice:
		if t := v.Type(); t == reflect.TypeOf(List{}) {
			return NoAnnotation
		} else if t == reflect.TypeOf(ByteArray{}) {
			return NoAnnotation
		} else if t == reflect.TypeOf(IntArray{}) {
			return NoAnnotation
		} else if t == reflect.TypeOf(LongArray{}) {
			return NoAnnotation
		}
		switch t := v.Type().Elem(); t.Kind() {
		case reflect.Int8:
			return ByteArrayAnnotation
		case reflect.Int32, reflect.Uint16:
			return IntArrayAnnotation
		case reflect.Int64, reflect.Uint32:
			return LongArrayAnnotation
		default:
			fmt.Printf("unknown array type '%s' (%s)\n", t, t.Kind())
			return InferenceArrayAnnotation
		}
	case reflect.Int8:
		return ByteAnnotation
	case reflect.Int16, reflect.Uint8:
		return ShortAnnotation
	case reflect.Int32, reflect.Uint16:
		return IntAnnotation
	case reflect.Int64, reflect.Uint32:
		return LongAnnotation
	case reflect.Float32:
		return FloatAnnotation
	case reflect.Float64:
		return DoubleAnnotation
	case reflect.String:
		return StringAnnotation
	case reflect.Pointer, reflect.Interface:
		return AnnotationOf(v.Elem())
	}
}
