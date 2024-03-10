package nbt

func Compose(tag NbtTag) []byte {
	var buf []byte
	buf = pushByte(buf, tag.Type())
	buf = pushString(buf, "") // root name not supported so its always empty

	buf = append(buf, tag.compose()...)
	return buf
}

func (t EndTag) compose() []byte {
	return []byte{}
}
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
	itemCap := len(t)
	if itemCap == 0 {
		buf := pushByte([]byte{}, Tag_End)
		return pushInt(buf, 0)
	}

	buf := pushByte([]byte{}, t[0].Type())
	buf = pushInt(buf, itemCap)
	for _, i := range t {
		buf = append(buf, i.compose()...)
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
