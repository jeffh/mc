package protocol

import (
	. "describe"
	"testing"
)

func TestProtocolBoolWriter(t *testing.T) {
	w, buf := createProtocolWriter()
	value := true
	err := ProtocolWriteBool(w, value)
	Expect(t, err, ToBeNil)

	var result byte
	err = readBytes(buf, &result)
	Expect(t, err, ToBeNil)
	Expect(t, result, ToEqual, byte(1))
}

func TestProtocolStringWriter(t *testing.T) {
	w, buf := createProtocolWriter()
	data := "foobar"
	err := ProtocolWriteString(w, data)
	Expect(t, err, ToBeNil)

	var result string
	err = readBytes(buf, &result)
	Expect(t, err, ToBeNil)
	Expect(t, result, ToEqual, data)
}

func TestProtocolBytesWriter(t *testing.T) {
	w, buf := createProtocolWriter()
	data := []byte{1, 2, 3}
	err := ProtocolWriteByteSlice(w, data)
	Expect(t, err, ToBeNil)

	var size int16
	err = readBytes(buf, &size)
	Expect(t, err, ToBeNil)
	Expect(t, size, ToEqual, int16(len(data)))
	for _, ch := range data {
		var b byte
		err = readBytes(buf, &b)
		Expect(t, err, ToBeNil)
		Expect(t, ch, ToEqual, b)
	}
}
