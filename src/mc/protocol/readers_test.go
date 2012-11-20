package protocol

import (
	"bytes"
	. "describe"
	"fmt"
	"testing"
)

func createProtocolReader() (*Reader, *bytes.Buffer) {
	b := bytes.NewBuffer([]byte{})
	r := NewReader(b, nil, ClientPacketMapper)
	return r, b
}

func TestReaderCanDispatchToTypeTable(t *testing.T) {
	r, b := createProtocolReader()
	err := writeBytes(b, int32(13), "default", CreativeMode, NetherDimension, NormalDifficulty, int8(0), int8(8))
	Expect(t, err, ToBeNil)

	p := LoginRequest{}
	err = r.ReadStruct(&p)
	Expect(t, err, ToBeNil)

	Expect(t, p.EntityID, ToEqual, int32(13))
	Expect(t, p.LevelType, ToEqual, "default")
	Expect(t, p.GameMode, ToEqual, CreativeMode)
	Expect(t, p.Dimension, ToEqual, NetherDimension)
	Expect(t, p.Difficulty, ToEqual, NormalDifficulty)
	Expect(t, p.MaxPlayers, ToEqual, int8(8))
}

func TestReaderCanParseDoublesAndBools(t *testing.T) {
	r, b := createProtocolReader()
	err := writeBytes(b, float64(13), float64(14), float64(1), float64(64), int8(1))
	Expect(t, err, ToBeNil)

	p := PlayerPosition{}
	err = r.ReadStruct(&p)
	Expect(t, err, ToBeNil)

	Expect(t, p.X, ToEqual, float64(13))
	Expect(t, p.Y, ToEqual, float64(14))
	Expect(t, p.Stance, ToEqual, float64(1))
	Expect(t, p.Z, ToEqual, float64(64))
	Expect(t, p.IsOnGround, ToEqual, true)
}

func TestReaderCanParseFloats(t *testing.T) {
	r, b := createProtocolReader()
	err := writeBytes(b, int16(13), int16(14), float32(3.5))
	Expect(t, err, ToBeNil)

	p := UpdateHealth{}
	err = r.ReadStruct(&p)
	Expect(t, err, ToBeNil)

	Expect(t, p.Health, ToEqual, int16(13))
	Expect(t, p.Food, ToEqual, int16(14))
	Expect(t, p.Saturation, ToEqual, float32(3.5))
}

/////////////////////////////////////////////////////////////////////
func TestProtocolHandshakeReader(t *testing.T) {
	r, b := createProtocolReader()
	err := writeBytes(b, byte(47), "_AlexM", "localhost", int32(25565))
	Expect(t, err, ToBeNil)

	var h Handshake
	err = r.ReadStruct(&h)
	Expect(t, err, ToBeNil)
	Expect(t, h.Version, ToEqual, byte(47))
	Expect(t, h.Username, ToEqual, "_AlexM")
	Expect(t, h.Hostname, ToEqual, "localhost")
	Expect(t, h.Port, ToEqual, int32(25565))
}

func TestProtocolEncryptionKeyResponseReader(t *testing.T) {
	r, b := createProtocolReader()
	err := writeBytes(b, int16(2), byte(1), byte(3), int16(2), byte(4), byte(5))
	Expect(t, err, ToBeNil)

	var ekr EncryptionKeyResponse
	err = r.ReadStruct(&ekr)
	Expect(t, err, ToBeNil)
	Expect(t, ekr.SharedSecret, ToDeeplyEqual, []byte{1, 3})
	Expect(t, ekr.VerifyToken, ToDeeplyEqual, []byte{4, 5})
}

func TestProtocolEncryptionKeyRequestReader(t *testing.T) {
	r, b := createProtocolReader()
	err := writeBytes(b, "MyServerName", int16(2), byte(1), byte(3), int16(2), byte(4), byte(5))
	Expect(t, err, ToBeNil)

	var ekr EncryptionKeyRequest
	err = r.ReadStruct(&ekr)
	Expect(t, err, ToBeNil)
	Expect(t, ekr.ServerID, ToEqual, "MyServerName")
	Expect(t, ekr.PublicKey, ToDeeplyEqual, []byte{1, 3})
	Expect(t, ekr.VerifyToken, ToDeeplyEqual, []byte{4, 5})
}

func TestReadPacket(t *testing.T) {
	r, b := createProtocolReader()
	err := writeBytes(b, byte(0xFD), "MyServerName", int16(2), byte(1), byte(3), int16(2), byte(4), byte(5))
	Expect(t, err, ToBeNil)

	obj, err := r.ReadPacket()
	Expect(t, err, ToBeNil)
	ekr, ok := obj.(*EncryptionKeyRequest)
	fmt.Printf("Type: %#v\n", ekr)
	Expect(t, ok, ToBeTrue)
	Expect(t, ekr.ServerID, ToEqual, "MyServerName")
	Expect(t, ekr.PublicKey, ToDeeplyEqual, []byte{1, 3})
	Expect(t, ekr.VerifyToken, ToDeeplyEqual, []byte{4, 5})
}
