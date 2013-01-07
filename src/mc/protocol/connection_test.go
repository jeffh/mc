package protocol

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	. "describe"
	"fmt"
	"testing"
)

func notSoRandomNumber() (int64, error) {
	return 1, nil
}

type packetBuffer struct {
	packets   []interface{}
	wUpgrader WriterFactory
	rUpgrader ReaderFactory
}

func NewPacketBuffer() *packetBuffer {
	return &packetBuffer{
		packets: make([]interface{}, 0),
	}
}

func (p *packetBuffer) IsEmpty() bool {
	return len(p.packets) == 0
}

func (p *packetBuffer) UpgradeWriter(f WriterFactory) { p.wUpgrader = f }
func (p *packetBuffer) UpgradeReader(f ReaderFactory) { p.rUpgrader = f }

func (p *packetBuffer) ReadPacket() (interface{}, error) {
	if len(p.packets) < 1 {
		return nil, fmt.Errorf("No packets to read!")
	}
	obj := p.packets[0]
	p.packets = p.packets[1:]
	return obj, nil
}

func (p *packetBuffer) WritePacket(v interface{}) error {
	p.packets = append(p.packets, v)
	return nil
}

////////////////////////////////////////////////////////////////////////////
func ToWritePacket(buf WritePacketer, packet interface{}) (string, bool) {
	err := buf.WritePacket(packet)
	if err != nil {
		return fmt.Sprintf("Expected nil from WritePacket, got: %#v", err), false
	}
	return "", true
}

func ToReadPacket(buf ReadPacketer, expected interface{}) (string, bool) {
	v, err := buf.ReadPacket()
	if err != nil {
		return err.Error(), false
	}
	return ToEqual(v, expected)
}

////////////////////////////////////////////////////////////////////////////
func createConnection() (*Connection, *packetBuffer, *packetBuffer) {
	wbuf := NewPacketBuffer()
	rbuf := NewPacketBuffer()
	c := NewConnection(rbuf, wbuf)
	return c, rbuf, wbuf
}

func createPPK() (*rsa.PrivateKey, []byte, error) {
	priv, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		return nil, nil, err
	}
	pub, err := x509.MarshalPKIXPublicKey(&priv.PublicKey)
	return priv, pub, err
}

////////////////////////////////////////////////////////////////////////////

func TestConnectionCanWrite(t *testing.T) {
	c, _, wbuf := createConnection()
	handshake := &Handshake{
		Version:  47,
		Username: "Joe Smoe",
		Hostname: "localhost",
		Port:     25565,
	}
	Expect(t, c, ToWritePacket, handshake)
	Expect(t, wbuf, ToReadPacket, handshake)
}

func TestConnectionCanRead(t *testing.T) {
	c, rbuf, _ := createConnection()
	handshake := &Handshake{
		Version:  47,
		Username: "Joe Smoe",
		Hostname: "localhost",
		Port:     25565,
	}
	err := rbuf.WritePacket(handshake)
	Expect(t, err, ToBeNil)
	Expect(t, c, ToReadPacket, handshake)
}

func TestConnectionIsNotEncryptedByDefault(t *testing.T) {
	c, _, _ := createConnection()
	Expect(t, c.IsEncrypted(), Not(ToBeTrue))
}

func TestConnectionIsEncryptedWhenSecretExists(t *testing.T) {
	c, _, _ := createConnection()
	c.Encryption.SharedKey = []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	Expect(t, c.IsEncrypted(), ToBeTrue)
}

////////////////////////////////////////////////////////////////////////////
func TestCanNegociateEncryptedConnection(t *testing.T) {
	c, rbuf, wbuf := createConnection()
	priv, pub, err := createPPK()
	Expect(t, err, ToBeNil)
	verifyToken := []byte{1, 2, 3, 4}

	// server gets handshake, sends EKRequest
	Expect(t, rbuf, ToWritePacket, &EncryptionKeyRequest{
		ServerID:    "-", // no user verification server ?
		PublicKey:   pub,
		VerifyToken: verifyToken,
	})
	// server gets EKResponse, returns empty EKResponse
	Expect(t, rbuf, ToWritePacket, &EncryptionKeyResponse{
		SharedSecret: []byte{},
		VerifyToken:  []byte{},
	})
	// server promotes to encrypted socket

	handshake := &Handshake{
		Version:  47,
		Username: "Joe Smoe",
		Hostname: "localhost",
		Port:     25565,
	}
	secret := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	err = EstablishEncryptedConnection(c, handshake, secret)

	Expect(t, err, ToBeNil)

	// client should send handshake
	Expect(t, wbuf, ToReadPacket, handshake)
	// client should send EKResponse
	p, err := wbuf.ReadPacket()
	Expect(t, err, ToBeNil)
	ekRes, ok := p.(*EncryptionKeyResponse)
	Expect(t, ok, ToBeTrue)

	// we need to decrypt the fields
	dVerifyToken, err := rsa.DecryptPKCS1v15(rand.Reader, priv, ekRes.VerifyToken)
	Expect(t, err, ToBeNil)
	Expect(t, dVerifyToken, ToEqual, verifyToken)
	dSecret, err := rsa.DecryptPKCS1v15(rand.Reader, priv, ekRes.SharedSecret)
	Expect(t, err, ToBeNil)
	Expect(t, dSecret, ToEqual, secret)

	// shouldn't have any extra data
	Expect(t, rbuf.IsEmpty(), ToBeTrue)
	Expect(t, wbuf.IsEmpty(), ToBeTrue)
	// connection should be modified
	Expect(t, c.IsEncrypted(), ToBeTrue)
	Expect(t, c.ServerID, ToBe, "-")
	Expect(t, c.Encryption.SharedKey, ToEqual, secret)
}
