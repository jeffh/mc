package nbt

import (
    . "describe"
    "testing"
    "os"
    "path/filepath"
    "bytes"
    "io"
    "compress/gzip"
)

func fixtureAsReader(filename string) io.Reader {
    buf := bytes.NewBuffer([]byte{})
    filepath, err := filepath.Abs(filename)
    if err != nil { panic(err) }

    fi, err := os.Open(filepath)
    if err != nil { panic(err) }
    defer fi.Close()

    _, err = io.Copy(buf, fi)
    if err != nil { panic(err) }
    return buf
}

func TestReadingWithNoReaderIsAnError(t *testing.T) {
    nbt := NewFile(nil)
    _, err := nbt.Read()
    Expect(t, err.Error(), ToEqual, "No Reader specified")
}

func TestReadingSimpleNBT(t *testing.T) {
    r := fixtureAsReader("fixtures/test.nbt")
    reader, err := gzip.NewReader(r)
    Expect(t, err, ToBeNil)

    nbt := NewFile(reader)
    tag, err := nbt.Read()
    Expect(t, err, ToBeNil)
    Expect(t, tag, Not(ToEqual), InvalidTag)
    compound, ok := tag.Value.(Compound)
    Expect(t, ok, ToBeTrue)
    Expect(t, tag.Name, ToEqual, "hello world")
    Expect(t, compound["name"], ToDeeplyEqual, Tag{
        Type: TagTypeString,
        Name: "name",
        Value: "Bananrama",
    })
}

func TestReadingFullNBT(t *testing.T){
    r := fixtureAsReader("fixtures/bigtest.nbt")
    reader, err := gzip.NewReader(r)
    Expect(t, err, ToBeNil)

    nbt := NewFile(reader)
    _, err = nbt.Read()
    Expect(t, err, ToBeNil)
}
