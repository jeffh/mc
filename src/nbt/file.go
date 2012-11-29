// NBT - Provide Minecraft NBT file format parsing capabilities
package nbt

import (
    "io"
    "errors"
    "encoding/binary"
)

//////////////////////////////////////////////////////////////////

// Handles the translation to and from NBT files.
// Currently, only reading is supported
type File struct {
    reader io.Reader
    ByteOrder binary.ByteOrder
}

// Creates a new NBT File instance that accepts an io.Reader that is
// capable of processing Minecraft NBT file formats.
//
// To process NBT Files under GZip or Zlib, wrap the io.Reader
// with a reader from compress/gzip or compress/zlib as appropriate.
//
// If you wish to read PocketMinecraft NBT files, manually set the
// ByteOrder field to binary.LittleEndian:
//
//     r, err := gzip.NewReader(fileReader)
//     if err != nil { panic(err) }
//     f := NewFile(r)
//     f.ByteOrder = binary.LittleEndian
func NewFile(r io.Reader) *File {
    return &File{
        reader: r,
        ByteOrder: binary.BigEndian,
    }
}

// Parses the entire NBT format. Each data type is contained in a
// Tag struct. The Tag returned is always contains a Value of
// a Compound.
//
// It is recommended to use ReadInto() instead.
//
// Returns an error if parsing failed in any way.
//
func (f *File) Read() (Tag, error) {
    if f.reader == nil {
        return InvalidTag, errors.New("No Reader specified")
    }
    r := &reader{f.reader, f.ByteOrder}
    return r.Read()
}

// Parses the NBT format into a struct. This behaves similarly to
// encoding/xml or encoding/json with the tag prefix of 'nbt'.
//
// With the exception of slices and stirngs, all other types must
// be specified in fixed-size data types. So int is not a valid type,
// but int32 is valid.
//
// The nbt.Name type can be used to read Tag Names when necessary
//
// For example, the following structs can parse bigtest.nbt.
// See http://www.wiki.vg/NBT#bigtest.nbt for reference of the format
//
//    type itemData struct {
//        Name string `nbt:"name"`
//        Value float32 `nbt:"value"`
//    }
//    type compoundItem struct {
//        CreatedOn int64 `nbt:"created-on"`
//        Name string `nbt:"name"`
//    }
//    var data struct {
//        Name nbt.Name
//        IntTest int32 `nbt:"intTest"`
//        ByteTest byte `nbt:"byteTest"`
//        StringTest string `nbt:"stringTest"`
//        DoubleTest float64 `nbt:"doubleTest"`
//        FloatTest float32 `nbt:"floatTest"`
//        LongTest int64 `nbt:"longTest"`
//        ShortTest int16 `nbt:"shortTest"`
//        ByteArrayTest []byte `nbt:"byteArrayTest (the first 1000 values of (n*n*255+n*7)%100, starting with n=0 (0, 62, 34, 16, 8, ...))"`
//        LongList []int64 `nbt:"listTest (long)"`
//        CompoundTest struct {
//            Egg itemData `nbt:"egg"`
//            Ham itemData `nbt:"ham"`
//        } `nbt:"nested compound test"`
//    }
//
//    err := nbtFile.ReadInto(&data)
//
// Invalid fields are silently ignored. But returns any errors
// when processing the NBT file.
//
func (f *File) ReadInto(v interface{}) error {
    tag, err := f.Read()
    if err != nil { return err }
    m := &mapper{ tag }
    return m.Read(v)
}

