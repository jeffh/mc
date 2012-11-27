package nbt

import (
    "io"
    "errors"
    "encoding/binary"
)


//////////////////////////////////////////////////////////////////

type File struct {
    reader io.Reader
    ByteOrder binary.ByteOrder
}

func NewFile(r io.Reader) *File {
    return &File{
        reader: r,
        ByteOrder: binary.BigEndian,
    }
}

func (f *File) Read() (Tag, error) {
    if f.reader == nil {
        return InvalidTag, errors.New("No Reader specified")
    }
    r := NewReader(f.reader, f.ByteOrder)
    return r.Next()
}

//////////////////////////////////////////////////////////////////

