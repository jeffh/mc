package nbt

import (
	"bytes"
	"compress/gzip"
	. "describe"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func fixtureAsReader(filename string) io.Reader {
	buf := bytes.NewBuffer([]byte{})
	filepath, err := filepath.Abs(filename)
	if err != nil {
		panic(err)
	}

	fi, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer fi.Close()

	_, err = io.Copy(buf, fi)
	if err != nil {
		panic(err)
	}
	return buf
}

func TestReadingWithNoReaderIsAnError(t *testing.T) {
	nbt := NewFile(nil)
	_, err := nbt.Read()
	Expect(t, err.Error(), ToBe, "No Reader specified")
}

func TestReadingSimpleNBTAsStruct(t *testing.T) {
	r := fixtureAsReader("fixtures/test.nbt")
	reader, err := gzip.NewReader(r)
	Expect(t, err, ToBeNil)

	var value struct {
		Name    Name
		TheName string `nbt:"name"`
	}

	nbt := NewFile(reader)
	err = nbt.ReadInto(&value)
	Expect(t, err, ToBeNil)
	Expect(t, value.Name, ToBe, Name("hello world"))
	Expect(t, value.TheName, ToBe, "Bananrama")
}

func TestReadingSimpleNBT(t *testing.T) {
	r := fixtureAsReader("fixtures/test.nbt")
	reader, err := gzip.NewReader(r)
	Expect(t, err, ToBeNil)

	nbt := NewFile(reader)
	tag, err := nbt.Read()
	Expect(t, err, ToBeNil)
	Expect(t, tag, Not(ToBe), InvalidTag)
	compound, ok := tag.Value.(Compound)
	Expect(t, ok, ToBeTrue)
	Expect(t, tag.Name, ToBe, "hello world")
	Expect(t, compound["name"], ToEqual, Tag{
		Type:  TagTypeString,
		Name:  "name",
		Value: "Bananrama",
	})
}

type itemData struct {
	Name  string  `nbt:"name"`
	Value float32 `nbt:"value"`
}
type compoundItem struct {
	CreatedOn int64  `nbt:"created-on"`
	Name      string `nbt:"name"`
}

func TestReadIntoFullNBT(t *testing.T) {
	r := fixtureAsReader("fixtures/bigtest.nbt")
	reader, err := gzip.NewReader(r)
	Expect(t, err, ToBeNil)

	var data struct {
		Name          Name
		IntTest       int32   `nbt:"intTest"`
		ByteTest      byte    `nbt:"byteTest"`
		StringTest    string  `nbt:"stringTest"`
		DoubleTest    float64 `nbt:"doubleTest"`
		FloatTest     float32 `nbt:"floatTest"`
		LongTest      int64   `nbt:"longTest"`
		ShortTest     int16   `nbt:"shortTest"`
		ByteArrayTest []byte  `nbt:"byteArrayTest (the first 1000 values of (n*n*255+n*7)%100, starting with n=0 (0, 62, 34, 16, 8, ...))"`
		LongList      []int64 `nbt:"listTest (long)"`
		CompoundTest  struct {
			Egg itemData `nbt:"egg"`
			Ham itemData `nbt:"ham"`
		} `nbt:"nested compound test"`
		ListCompound []compoundItem `nbt:"listTest (compound)"`
	}

	expectedListCompound := []compoundItem{
		{int64(1264099775885), "Compound tag #0"},
		{int64(1264099775885), "Compound tag #1"},
	}

	expectedBytes := []byte{}
	for i := 0; i < 1000; i++ {
		expectedBytes = append(expectedBytes, byte((i*i*255+i*7)%100))
	}

	nbt := NewFile(reader)
	err = nbt.ReadInto(&data)
	Expect(t, err, ToBeNil)
	Expect(t, data.Name, ToBe, Name("Level"))
	Expect(t, data.IntTest, ToBe, int32(2147483647))
	Expect(t, data.ByteTest, ToBe, byte(127))
	Expect(t, data.StringTest, ToBe, "HELLO WORLD THIS IS A TEST STRING \xc3\x85\xc3\x84\xc3\x96!")
	Expect(t, data.DoubleTest, ToBe, float64(0.49312871321823148))
	Expect(t, data.FloatTest, ToBe, float32(0.49823147058486938))
	Expect(t, data.LongTest, ToBe, int64(9223372036854775807))
	Expect(t, data.ShortTest, ToBe, int16(32767))
	Expect(t, data.ByteArrayTest, ToEqual, expectedBytes)
	Expect(t, data.LongList, ToEqual, []int64{11, 12, 13, 14, 15})
	Expect(t, data.CompoundTest.Egg, ToEqual, itemData{"Eggbert", 0.5})
	Expect(t, data.CompoundTest.Ham, ToEqual, itemData{"Hampus", 0.75})
	Expect(t, data.ListCompound, ToEqual, expectedListCompound)
}

func TestReadingFullNBT(t *testing.T) {
	r := fixtureAsReader("fixtures/bigtest.nbt")
	reader, err := gzip.NewReader(r)
	Expect(t, err, ToBeNil)

	nbt := NewFile(reader)
	tag, err := nbt.Read()
	Expect(t, err, ToBeNil)

	byteArrayTest := "byteArrayTest (the first 1000 values of (n*n*255+n*7)%100, starting with n=0 (0, 62, 34, 16, 8, ...))"
	expectedBytes := []byte{}
	for i := 0; i < 1000; i++ {
		expectedBytes = append(expectedBytes, byte((i*i*255+i*7)%100))
	}
	expectedTag := Tag{
		Name: "Level",
		Type: TagTypeCompound,
		Value: Compound{
			"nested compound test": {"nested compound test", TagTypeCompound, Compound{
				"egg": {"egg", TagTypeCompound, Compound{
					"name":  {"name", TagTypeString, "Eggbert"},
					"value": {"value", TagTypeFloat, float32(0.5)},
				}},
				"ham": {"ham", TagTypeCompound, Compound{
					"name":  {"name", TagTypeString, "Hampus"},
					"value": {"value", TagTypeFloat, float32(0.75)},
				}},
			}},
			"intTest":    {"intTest", TagTypeInt, int32(2147483647)},
			"byteTest":   {"byteTest", TagTypeByte, byte(127)},
			"stringTest": {"stringTest", TagTypeString, "HELLO WORLD THIS IS A TEST STRING \xc3\x85\xc3\x84\xc3\x96!"},
			"listTest (long)": {"listTest (long)", TagTypeList, List{
				TagTypeLong,
				[]interface{}{
					int64(11),
					int64(12),
					int64(13),
					int64(14),
					int64(15),
				},
			}},
			"doubleTest": {"doubleTest", TagTypeDouble, float64(0.49312871321823148)},
			"floatTest":  {"floatTest", TagTypeFloat, float32(0.49823147058486938)},
			"longTest":   {"longTest", TagTypeLong, int64(9223372036854775807)},
			"listTest (compound)": {"listTest (compound)", TagTypeList, List{
				TagTypeCompound,
				[]interface{}{
					Compound{
						"created-on": {"created-on", TagTypeLong, int64(1264099775885)},
						"name":       {"name", TagTypeString, "Compound tag #0"},
					},
					Compound{
						"created-on": {"created-on", TagTypeLong, int64(1264099775885)},
						"name":       {"name", TagTypeString, "Compound tag #1"},
					},
				},
			}},
			byteArrayTest: {byteArrayTest, TagTypeByteArray, expectedBytes},
			"shortTest":   {"shortTest", TagTypeShort, int16(32767)},
		},
	}
	Expect(t, tag, ToEqual, expectedTag)
}
