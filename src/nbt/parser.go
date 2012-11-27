package nbt

import (
    "io"
    "encoding/binary"
    "unicode/utf8"
    "fmt"
    "reflect"
)

func createTag(t TagType, v interface{}, name string) Tag {
    return Tag {
        Name: name,
        Type: t,
        Value: v,
    }
}

//////////////////////////////////////////////////////////////////

type Reader struct {
    r io.Reader
    bo binary.ByteOrder
}

func NewReader(r io.Reader, bo binary.ByteOrder) *Reader {
    return &Reader{r: r, bo: bo}
}

func (r *Reader) read(v interface{}) error {
    return binary.Read(r.r, r.bo, v)
}

func (r *Reader) consumeString() (string, error) {
    var size uint16
    err := r.read(&size)
    if err != nil { return "", err }

    var byt byte
    raw := make([]byte, size)
    for i := uint16(0); i<size; i++ {
        err = r.read(&byt)
        if err != nil { return "", err }
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

func (r *Reader) consumeTagEnd() (Tag, error) {
    return Tag { Type: TagTypeEnd }, nil
}

func (r *Reader) consumeStandardTag(v interface{}) (Tag, error) {
    str, err := r.consumeString()
    if err != nil { return InvalidTag, err }

    value := reflect.New(reflect.TypeOf(v))
    err = r.read(value.Interface())
    if err != nil { return InvalidTag, nil }
    return createTag(TagTypeByte, value.Elem().Interface(), str), nil
}

func (r *Reader) consumeTagByte() (Tag, error) {
    return r.consumeStandardTag(byte(0))
}

func (r *Reader) consumeTagShort() (Tag, error) {
    return r.consumeStandardTag(int16(0))
}

func (r *Reader) consumeTagInt() (Tag, error) {
    return r.consumeStandardTag(int32(0))
}

func (r *Reader) consumeTagLong() (Tag, error) {
    return r.consumeStandardTag(int64(0))
}

func (r *Reader) consumeTagFloat() (Tag, error) {
    return r.consumeStandardTag(float32(0))
}

func (r *Reader) consumeTagDouble() (Tag, error) {
    return r.consumeStandardTag(float64(0))
}

func (r *Reader) consumeTagString() (Tag, error) {
    name, err := r.consumeString()
    if err != nil { return InvalidTag, err }

    value, err := r.consumeString()
    if err != nil { return InvalidTag, err }
    return createTag(TagTypeString, value, name), nil
}

func (r *Reader) consumeTagByteArray() (Tag, error) {
    name, err := r.consumeString()
    if err != nil { return InvalidTag, err }

    var size int32
    err = r.read(&size)
    if err != nil { return InvalidTag, err }

    bytes := make([]byte, size)
    for i := int32(0); i<size; i++ {
        var byt byte
        err = r.read(&byt)
        if err != nil { return InvalidTag, err }
        bytes[i] = byt
    }
    return createTag(TagTypeByteArray, bytes, name), nil
}

func (r *Reader) consumeTagCompound() (Tag, error) {
    str, err := r.consumeString()
    if err != nil { return InvalidTag, err }

    compound := make(Compound)
    for {
        tag, err := r.Next()
        if err != nil { return InvalidTag, err }
        if tag.Type == TagTypeEnd { break }

        compound[tag.Name] = tag
    }
    return createTag(TagTypeCompound, compound, str), nil
}

func (r *Reader) consumeTagIntArray() (Tag, error) {
    name, err := r.consumeString()
    if err != nil { return InvalidTag, err }

    var size int32
    err = r.read(&size)
    if err != nil { return InvalidTag, err }

    numbers := make([]int32, size)
    for i := int32(0); i<size; i++ {
        var n int32
        err = r.read(&n)
        if err != nil { return InvalidTag, err }
        numbers[i] = n
    }
    return createTag(TagTypeIntArray, numbers, name), nil
}

func (r *Reader) Next() (Tag, error) {
    var typ TagType
    err := r.read(&typ)
    if err != nil { return InvalidTag, err }

    switch typ {
    case TagTypeEnd:
        return r.consumeTagEnd()
    case TagTypeByte:
        return r.consumeTagByte()
    case TagTypeShort:
        return r.consumeTagShort()
    case TagTypeInt:
        return r.consumeTagInt()
    case TagTypeLong:
        return r.consumeTagLong()
    case TagTypeFloat:
        return r.consumeTagFloat()
    case TagTypeDouble:
        return r.consumeTagDouble()
    case TagTypeByteArray:
        return r.consumeTagByteArray()
    case TagTypeString:
        return r.consumeTagString()
    case TagTypeList:
    case TagTypeCompound:
        return r.consumeTagCompound()
    case TagTypeIntArray:
        return r.consumeTagIntArray()
    }

    return InvalidTag, fmt.Errorf("Unknown Tag Type: %#v", typ)
}
