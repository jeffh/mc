package nbt

import (
	"encoding/binary"
	"fmt"
	"io"
	"reflect"
	"unicode/utf8"
)

func createTag(t TagType, v interface{}, name string) Tag {
	return Tag{
		Name:  name,
		Type:  t,
		Value: v,
	}
}

//////////////////////////////////////////////////////////////////

type reader struct {
	reader    io.Reader
	byteOrder binary.ByteOrder
}

func (r *reader) read(v interface{}) error {
	return binary.Read(r.reader, r.byteOrder, v)
}

func (r *reader) consumeString() (string, error) {
	var size uint16
	err := r.read(&size)
	if err != nil {
		return "", err
	}

	var byt byte
	raw := make([]byte, size)
	for i := uint16(0); i < size; i++ {
		err = r.read(&byt)
		if err != nil {
			return "", err
		}
		raw[int(i)] = byt
	}

	runes := make([]rune, 0)
	for len(raw) > 0 {
		r, s := utf8.DecodeRune(raw)
		if r == utf8.RuneError {
			return "", fmt.Errorf("Got invalid rune on bytes: %#v", raw)
		}
		runes = append(runes, r)
		raw = raw[s:]
	}
	return string(runes), nil
}

func (r *reader) consumeTagEnd() (Tag, error) {
	return Tag{Type: TagTypeEnd}, nil
}

func (r *reader) consumeStandardTag(t TagType, v interface{}, includesName bool) (Tag, error) {
	var str string
	var err error
	if includesName {
		str, err = r.consumeString()
		if err != nil {
			return InvalidTag, err
		}
	}

	value := reflect.New(reflect.TypeOf(v))
	err = r.read(value.Interface())
	if err != nil {
		return InvalidTag, nil
	}
	return createTag(t, value.Elem().Interface(), str), nil
}

func (r *reader) consumeTagByte(includesName bool) (Tag, error) {
	return r.consumeStandardTag(TagTypeByte, byte(0), includesName)
}

func (r *reader) consumeTagShort(includesName bool) (Tag, error) {
	return r.consumeStandardTag(TagTypeShort, int16(0), includesName)
}

func (r *reader) consumeTagInt(includesName bool) (Tag, error) {
	return r.consumeStandardTag(TagTypeInt, int32(0), includesName)
}

func (r *reader) consumeTagLong(includesName bool) (Tag, error) {
	return r.consumeStandardTag(TagTypeLong, int64(0), includesName)
}

func (r *reader) consumeTagFloat(includesName bool) (Tag, error) {
	return r.consumeStandardTag(TagTypeFloat, float32(0), includesName)
}

func (r *reader) consumeTagDouble(includesName bool) (Tag, error) {
	return r.consumeStandardTag(TagTypeDouble, float64(0), includesName)
}

func (r *reader) consumeTagString(includesName bool) (Tag, error) {
	var name string
	var err error
	if includesName {
		name, err = r.consumeString()
		if err != nil {
			return InvalidTag, err
		}
	}

	value, err := r.consumeString()
	if err != nil {
		return InvalidTag, err
	}
	return createTag(TagTypeString, value, name), nil
}

func (r *reader) consumeTagByteArray(includesName bool) (Tag, error) {
	var name string
	var err error
	if includesName {
		name, err = r.consumeString()
		if err != nil {
			return InvalidTag, err
		}
	}

	var size int32
	err = r.read(&size)
	if err != nil {
		return InvalidTag, err
	}

	bytes := make([]byte, size)
	for i := int32(0); i < size; i++ {
		var byt byte
		err = r.read(&byt)
		if err != nil {
			return InvalidTag, err
		}
		bytes[i] = byt
	}
	return createTag(TagTypeByteArray, bytes, name), nil
}

func (r *reader) consumeTagCompound(includesName bool) (Tag, error) {
	var str string
	var err error
	if includesName {
		str, err = r.consumeString()
		if err != nil {
			return InvalidTag, err
		}
	}

	compound := make(Compound)
	for {
		tag, err := r.next(true)
		if err != nil {
			return InvalidTag, err
		}
		if tag.Type == TagTypeEnd {
			break
		}

		compound[tag.Name] = tag
	}
	return createTag(TagTypeCompound, compound, str), nil
}

func (r *reader) consumeTagIntArray(includesName bool) (Tag, error) {
	var name string
	var err error
	if includesName {
		name, err = r.consumeString()
		if err != nil {
			return InvalidTag, err
		}
	}

	var size int32
	err = r.read(&size)
	if err != nil {
		return InvalidTag, err
	}

	numbers := make([]int32, size)
	for i := int32(0); i < size; i++ {
		var n int32
		err = r.read(&n)
		if err != nil {
			return InvalidTag, err
		}
		numbers[i] = n
	}
	return createTag(TagTypeIntArray, numbers, name), nil
}

func (r *reader) consumeTagList(includesName bool) (Tag, error) {
	var name string
	var err error
	if includesName {
		name, err = r.consumeString()
		if err != nil {
			return InvalidTag, err
		}
	}

	var list List
	err = r.read(&list.Type)
	if err != nil {
		return InvalidTag, err
	}

	var size int32
	err = r.read(&size)
	if err != nil {
		return InvalidTag, err
	}

	list.Values = make([]interface{}, size)
	for i := int32(0); i < size; i++ {
		val, err := r.nextType(list.Type, false)
		if err != nil {
			return InvalidTag, err
		}
		list.Values[i] = val.Value
	}
	return createTag(TagTypeList, list, name), err
}

func (r *reader) nextType(typ TagType, includesName bool) (Tag, error) {
	switch typ {
	case TagTypeEnd:
		return r.consumeTagEnd()
	case TagTypeByte:
		return r.consumeStandardTag(TagTypeByte, byte(0), includesName)
	case TagTypeShort:
		return r.consumeStandardTag(TagTypeShort, int16(0), includesName)
	case TagTypeInt:
		return r.consumeStandardTag(TagTypeInt, int32(0), includesName)
	case TagTypeLong:
		return r.consumeStandardTag(TagTypeLong, int64(0), includesName)
	case TagTypeFloat:
		return r.consumeStandardTag(TagTypeFloat, float32(0), includesName)
	case TagTypeDouble:
		return r.consumeStandardTag(TagTypeDouble, float64(0), includesName)
	case TagTypeByteArray:
		return r.consumeTagByteArray(includesName)
	case TagTypeString:
		return r.consumeTagString(includesName)
	case TagTypeList:
		return r.consumeTagList(includesName)
	case TagTypeCompound:
		return r.consumeTagCompound(includesName)
	case TagTypeIntArray:
		return r.consumeTagIntArray(includesName)
	}

	return InvalidTag, fmt.Errorf("Unknown Tag Type: %#v", typ)
}

func (r *reader) next(includesName bool) (Tag, error) {
	var typ TagType
	err := r.read(&typ)
	if err != nil {
		return InvalidTag, err
	}
	return r.nextType(typ, includesName)
}

func (r *reader) Read() (Tag, error) {
	return r.next(true)
}

///////////////////////////////////////////////////////////////////

type mapper struct {
	tag Tag
}

func (m *mapper) intoStruct(tag Tag, value reflect.Value) error {
	derefValue := value.Elem()
	size := derefValue.NumField()
	for i := 0; i < size; i++ {
		f := derefValue.Field(i)
		fType := derefValue.Type().Field(i)
		typ := f.Type()
		if typ == nameType {
			f.Set(reflect.ValueOf(Name(tag.Name)))
		} else {
			meta := fType.Tag.Get("nbt")
			if meta == "" {
				meta = fType.Name
			}

			compound := tag.Value.(Compound)
			item, ok := compound[meta]
			if !ok {
				continue
			}

			//m.into(item, f.Addr())
			data := reflect.New(fType.Type)
			err := m.into(item, data)
			if err != nil {
				return err
			}
			f.Set(data.Elem())
		}
	}
	return nil
}

func (m *mapper) intoSlice(tag Tag, value reflect.Value) error {
	// byte array
	_, ok := tag.Value.([]uint8)
	// value is a ptr to a slice
	if ok && value.Type().Elem().Elem().Kind() == reflect.Uint8 {
		value.Elem().Set(reflect.ValueOf(tag.Value))
		return nil
	}

	list, ok := tag.Value.(List)
	if !ok {
		return fmt.Errorf("Expected NBT List, got %s instead", reflect.TypeOf(tag.Value))
	}

	// list of any data type
	size := len(list.Values)
	slice := reflect.MakeSlice(value.Elem().Type(), size, size)
	for i, val := range list.Values {
		err := m.into(Tag{Type: list.Type, Value: val}, slice.Index(i).Addr())
		if err != nil {
			return err
		}
	}
	value.Elem().Set(slice)
	return nil
}

func (m *mapper) into(tag Tag, value reflect.Value) error {
	derefValue := value.Elem()
	switch derefValue.Kind() {
	case reflect.Struct: // == compound
		fmt.Printf("Struct Set: %#v = %#v\n", tag, value.Interface())
		return m.intoStruct(tag, value)
	case reflect.Slice, reflect.Array:
		return m.intoSlice(tag, value)
	case reflect.String, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Uint8, reflect.Float32, reflect.Float64:
		derefValue.Set(reflect.ValueOf(tag.Value))
		return nil
	}
	return fmt.Errorf("Unknown type: %s", derefValue.Kind().String())
}

func (m *mapper) Read(v interface{}) error {
	return m.into(m.tag, reflect.ValueOf(v))
}
