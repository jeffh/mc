package protocol

import (
	"bytes"
	. "github.com/jeffh/goexpect"
	"testing"
)

/////////////////////////////////////////////////////////////////////
func createProtocolWriter() (*Writer, *bytes.Buffer) {
	b := bytes.NewBuffer([]byte{})
	w := NewWriter(b, ClientPacketMapper, nil, nil)
	return w, b
}

func TestWriterCanWriteFloats(t *testing.T) {
	w, buf := createProtocolWriter()
	p := UpdateHealth{
		Health:     float32(13),
		Food:       int16(14),
		Saturation: float32(3.5),
	}
	err := w.WriteStruct(&p)
	Expect(t, err, ToBeNil)

	var Food int16
	var Health, Saturation float32
	err = readBytes(buf, &Health, &Food, &Saturation)
	Expect(t, err, ToBeNil)
	Expect(t, Health, ToEqual, float32(13))
	Expect(t, Food, ToEqual, int16(14))
	Expect(t, Saturation, ToEqual, float32(3.5))
}

func TestProtocolHandshakeWriter(t *testing.T) {
	w, b := createProtocolWriter()
	h := Handshake{
		Version:  47,
		Username: "_AlexM",
		Hostname: "localhost",
		Port:     int32(25565),
	}
	err := w.WriteStruct(&h)
	Expect(t, err, ToBeNil)

	var Version int8
	var Username, Hostname string
	var Port int32
	err = readBytes(b, &Version, &Username, &Hostname, &Port)
	Expect(t, err, ToBeNil)
	Expect(t, Version, ToEqual, int8(47))
	Expect(t, Username, ToEqual, "_AlexM")
	Expect(t, Hostname, ToEqual, "localhost")
	Expect(t, Port, ToEqual, int32(25565))
}

func TestWriterCanWriteDoublesAndBools(t *testing.T) {
	w, b := createProtocolWriter()
	p := PlayerPosition{
		X:          float64(13),
		Y:          float64(14),
		Stance:     float64(1),
		Z:          float64(64),
		IsOnGround: true,
	}
	err := w.WriteStruct(&p)
	Expect(t, err, ToBeNil)

	var X, Y, Z, Stance float64
	var IsOnGround int8
	readBytes(b, &X, &Y, &Stance, &Z, &IsOnGround)

	Expect(t, X, ToEqual, float64(13))
	Expect(t, Y, ToEqual, float64(14))
	Expect(t, Stance, ToEqual, float64(1))
	Expect(t, Z, ToEqual, float64(64))
	Expect(t, IsOnGround, ToEqual, int8(1))
}
