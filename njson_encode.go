package nbtreader

import (
	"bytes"
	"fmt"
	"reflect"
	"slices"
)

const startDetectingCyclesAfter = 1000

// Marshaler is the interface implemented by types that can marshal themselves into valid
// NJSON.
type Marshaler interface {
	MarshalNJSON() ([]byte, error)
}

func MarshalNJSON(v any) ([]byte, error) {
	if v, ok := v.(Marshaler); ok {
		return v.MarshalNJSON()
	}
	e := &njsonEncoderState{ptrSeen: map[any]struct{}{}}
	err := e.marshal(v)
	buf := slices.Clone(e.Bytes())
	return buf, err
}

// An UnsupportedValueError is returned by Marshal when attempting to encode an
// unsupported value.
type UnsupportedValueError struct {
	Str string
}

// Error implements [[error]]
func (e *UnsupportedValueError) Error() string {
	return "njson: unsupported value: " + e.Str
}

type njsonEncoderState struct {
	bytes.Buffer

	ptrLevel int
	ptrSeen  map[any]struct{}
}

// njsonError is an error wrapper type for internal use only.
// Panics with errors are wrapped in njsonError so that the top-level recover
// can distinguish intentional panics from this package.
type njsonError struct{ error }

// error aborts the encoding by panicking with err wrapped in njsonError.
func (e *njsonEncoderState) error(err error) {
	panic(njsonError{err})
}

func (e *njsonEncoderState) marshal(v any) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if je, ok := r.(njsonError); ok {
				err = je.error
			} else {
				panic(r)
			}
		}
	}()

	e.valueEncoder(reflect.ValueOf(v))
	return nil
}

func (e *njsonEncoderState) valueEncoder(v reflect.Value) {
	switch v.Type() {
	case reflect.TypeOf(Byte(0)), reflect.TypeOf(Short(0)), reflect.TypeOf(Int(0)), reflect.TypeOf(Long(0)):
		e.intEncoder(v)
		return
	case reflect.TypeOf(Float(0)), reflect.TypeOf(Double(0)):
		e.floatEncoder(v)
		return
	case reflect.TypeOf(String("")):
		e.stringEncoder(v)
		return
	case reflect.TypeOf(List{}):
		l := v.Interface().(List)
		if anno := l.TagType.Annotation(); anno != NoAnnotation {
			// allocation friendly prepending of type annotation
			l.Elements = append(l.Elements, nil)
			copy(l.Elements[1:], l.Elements)
			l.Elements[0] = String(anno.String())
		}
		v = reflect.ValueOf(l.Elements)
		fallthrough
	case reflect.TypeOf(ByteArray{}), reflect.TypeOf(IntArray{}), reflect.TypeOf(LongArray{}):
		e.arrayEncoder(v)
		return
	}

	switch kind := v.Kind(); kind {
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		e.intEncoder(v)
	case reflect.Float32, reflect.Float64:
		e.floatEncoder(v)
	case reflect.String:
		e.stringEncoder(v)
	case reflect.Slice:
		e.arrayEncoder(v)
	case reflect.Map:
		e.objectEncoder(v)

	case reflect.Pointer, reflect.Interface:
		e.valueEncoder(e.realValue(v))
	default:
		e.error(&UnsupportedValueError{fmt.Sprintf("encoding '%s' isn't supported yet: %s", kind, v.Type())})
	}
}

func (e *njsonEncoderState) intEncoder(v reflect.Value) {
	e.WriteString(fmt.Sprintf("%d", v.Int()))
}

func (e *njsonEncoderState) floatEncoder(v reflect.Value) {
	e.WriteString(fmt.Sprintf("%f", v.Float()))
}

func (e *njsonEncoderState) stringEncoder(v reflect.Value) {
	e.WriteString("\"" + e.realValue(v).String() + "\"")
}

func (e *njsonEncoderState) arrayEncoder(v reflect.Value) {
	e.WriteByte('[')
	for i := 0; i < v.Len(); i++ {
		val := v.Index(i)
		e.valueEncoder(val)
		if i < v.Len()-1 {
			e.WriteByte(',')
		}
	}
	e.WriteByte(']')
}

func (e *njsonEncoderState) objectEncoder(v reflect.Value) {
	if v.Type().Key().Kind() != reflect.String {
		e.error(&UnsupportedValueError{fmt.Sprintf("invalid map key type '%s'", v.Type().Key())})
	}

	e.WriteByte('{')
	iter := v.MapRange()
	for i := 0; iter.Next(); i++ {
		key, val := e.realValue(iter.Key()), e.realValue(iter.Value())
		e.WriteByte('"')
		e.WriteString(key.String())
		e.WriteString(AnnotationOf(val).String())
		e.WriteByte('"')
		e.WriteByte(':')
		e.valueEncoder(val)
		if i < v.Len()-1 {
			e.WriteByte(',')
		}
	}
	e.WriteByte('}')
}

func (e *njsonEncoderState) realValue(v reflect.Value) reflect.Value {
	switch v.Kind() {
	case reflect.Interface, reflect.Pointer:
	default:
		return v
	}

	if e.ptrLevel++; e.ptrLevel > startDetectingCyclesAfter {
		ptr := v.Interface()
		if _, ok := e.ptrSeen[ptr]; ok {

		}
		e.ptrSeen[ptr] = struct{}{}
	}

	v = e.realValue(v.Elem())
	e.ptrLevel--
	return v
}
