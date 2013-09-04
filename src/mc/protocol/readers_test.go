package protocol

import (
	"bytes"
	"fmt"
	. "github.com/jeffh/goexpect"
	"testing"
)

func createProtocolReader() (*Reader, *bytes.Buffer) {
	b := bytes.NewBuffer([]byte{})
	r := NewReader(b, ClientPacketMapper, nil, nil)
	return r, b
}

func TestReaderCanDispatchToTypeTable(t *testing.T) {
	r, b := createProtocolReader()
	err := writeBytes(b, int32(13), "default",
		GameModeCreative, GameDimensionNether, GameDifficultyNormal, int8(0), int8(8))
	Expect(t, err, ToBeNil)

	p := LoginRequest{}
	err = r.ReadStruct(&p)
	Expect(t, err, ToBeNil)

	Expect(t, p.EntityID, ToBe, int32(13))
	Expect(t, p.LevelType, ToBe, LevelType("default"))
	Expect(t, p.GameMode, ToBe, GameModeCreative)
	Expect(t, p.Dimension, ToBe, GameDimensionNether)
	Expect(t, p.Difficulty, ToBe, GameDifficultyNormal)
	Expect(t, p.MaxPlayers, ToBe, int8(8))
}

func TestReaderCanParseDoublesAndBools(t *testing.T) {
	r, b := createProtocolReader()
	err := writeBytes(b, float64(13), float64(14), float64(1), float64(64), int8(1))
	Expect(t, err, ToBeNil)

	p := PlayerPosition{}
	err = r.ReadStruct(&p)
	Expect(t, err, ToBeNil)

	Expect(t, p.X, ToBe, float64(13))
	Expect(t, p.Y, ToBe, float64(14))
	Expect(t, p.Stance, ToBe, float64(1))
	Expect(t, p.Z, ToBe, float64(64))
	Expect(t, p.IsOnGround, ToBe, true)
}

func TestReaderCanParseFloats(t *testing.T) {
	r, b := createProtocolReader()
	err := writeBytes(b, float32(13), int16(14), float32(3.5))
	Expect(t, err, ToBeNil)

	p := UpdateHealth{}
	err = r.ReadStruct(&p)
	Expect(t, err, ToBeNil)

	Expect(t, p.Health, ToBe, float32(13))
	Expect(t, p.Food, ToBe, int16(14))
	Expect(t, p.Saturation, ToBe, float32(3.5))
}

/////////////////////////////////////////////////////////////////////
func TestProtocolHandshakeReader(t *testing.T) {
	r, b := createProtocolReader()
	err := writeBytes(b, byte(47), "_AlexM", "localhost", int32(25565))
	Expect(t, err, ToBeNil)

	var h Handshake
	err = r.ReadStruct(&h)
	Expect(t, err, ToBeNil)
	Expect(t, h.Version, ToBe, byte(47))
	Expect(t, h.Username, ToBe, "_AlexM")
	Expect(t, h.Hostname, ToBe, "localhost")
	Expect(t, h.Port, ToBe, int32(25565))
}

func TestProtocolEncryptionKeyResponseReader(t *testing.T) {
	r, b := createProtocolReader()
	err := writeBytes(b, int16(2), byte(1), byte(3), int16(2), byte(4), byte(5))
	Expect(t, err, ToBeNil)

	var ekr EncryptionKeyResponse
	err = r.ReadStruct(&ekr)
	Expect(t, err, ToBeNil)
	Expect(t, ekr.SharedSecret, ToEqual, []byte{1, 3})
	Expect(t, ekr.VerifyToken, ToEqual, []byte{4, 5})
}

func TestProtocolEncryptionKeyRequestReader(t *testing.T) {
	r, b := createProtocolReader()
	err := writeBytes(b, "MyServerName", int16(2), byte(1), byte(3), int16(2), byte(4), byte(5))
	Expect(t, err, ToBeNil)

	var ekr EncryptionKeyRequest
	err = r.ReadStruct(&ekr)
	Expect(t, err, ToBeNil)
	Expect(t, ekr.ServerID, ToBe, "MyServerName")
	Expect(t, ekr.PublicKey, ToEqual, []byte{1, 3})
	Expect(t, ekr.VerifyToken, ToEqual, []byte{4, 5})
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
	Expect(t, ekr.ServerID, ToBe, "MyServerName")
	Expect(t, ekr.PublicKey, ToEqual, []byte{1, 3})
	Expect(t, ekr.VerifyToken, ToEqual, []byte{4, 5})
}
