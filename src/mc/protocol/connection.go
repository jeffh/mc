package protocol

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"fmt"
	"io"
)

// Handles the handshake to a minecraft server. Uses a plaintext connection.
func EstablishPlaintextConnection(c *Connection, h *Handshake) (err error) {
	err = c.WritePacket(h)
	if err != nil {
		return
	}

	p, err := c.ReadPacket()
	if err != nil {
		return
	}

	_, ok := p.(*EncryptionKeyRequest)
	if !ok {
		err = fmt.Errorf("Expected EncryptionKeyRequest packet, but got: %#v", p)
		return
	}

	// just force the spawn
	return
}

// Handles the handshake to a minecraft server. Secret is the shared key
// used by both parties for encryption.
//
// The encryption upgrading of the socket stream has to be done
// immediately after the connection has been established without errors.
//
// Currently broken. Use EstablishPlaintextConnection for now.
func EstablishEncryptedConnection(c *Connection, h *Handshake, secret []byte) (err error) {
	if len(secret) != 16 {
		panic("Secret must be 16-bytes")
	}

	err = c.WritePacket(h)
	if err != nil {
		return
	}

	p, err := c.ReadPacket()
	if err != nil {
		return
	}

	ekReq, ok := p.(*EncryptionKeyRequest)
	if !ok {
		err = fmt.Errorf("Expected EncryptionKeyRequest packet, but got: %#v", p)
		return
	}

	publicKey, err := x509.ParsePKIXPublicKey(ekReq.PublicKey)
	if err != nil {
		return
	}
	fmt.Printf("PublicKey: %#v\n", publicKey)

	c.Encryption.PublicKey = publicKey
	encSecret, err := c.Encryption.Encrypt(secret)
	if err != nil {
		return
	}

	encToken, err := c.Encryption.Encrypt(ekReq.VerifyToken)
	if err != nil {
		return
	}

	fmt.Printf("Secret: %#v\n", secret)
	fmt.Printf("Token: %#v\n", ekReq.VerifyToken)
	err = c.WritePacket(&EncryptionKeyResponse{
		SharedSecret: encSecret,
		VerifyToken:  encToken,
	})
	if err != nil {
		return
	}
	// expect response

	fmt.Println("Awaiting for Encryption Key Response...")
	p, err = c.ReadPacket()
	if err != nil {
		return
	}

	_, ok = p.(*EncryptionKeyResponse)
	if !ok {
		err = fmt.Errorf("Expected EncryptionKeyResponse packet, but got: %#v", p)
		return
	}

	// encrypt connection
	c.ServerID = ekReq.ServerID
	c.Encryption.SharedKey = secret

	return
}

//////////////////////////////////////////////////////////

type EncryptionProtocol struct {
	PublicKey interface{}
	SharedKey []byte
}

func (e *EncryptionProtocol) Encrypt(d []byte) ([]byte, error) {
	pk, ok := e.PublicKey.(*rsa.PublicKey)
	if ok {
		return rsa.EncryptPKCS1v15(rand.Reader, pk, d)
	}
	return nil, fmt.Errorf("Unknown PublicKey: %#v", e.PublicKey)
}

//////////////////////////////////////////////////////////

type WriterFactory func(w io.Writer) io.Writer
type ReaderFactory func(w io.Reader) io.Reader

type WritePacketer interface {
	WritePacket(v interface{}) error
	UpgradeWriter(WriterFactory)
}

type ReadPacketer interface {
	ReadPacket() (interface{}, error)
	UpgradeReader(ReaderFactory)
}

type Connection struct {
	Writer     WritePacketer
	Reader     ReadPacketer
	Encryption EncryptionProtocol
	ServerID   string
}

func NewConnection(reader ReadPacketer, writer WritePacketer) *Connection {
	return &Connection{
		Writer:     writer,
		Reader:     reader,
		Encryption: EncryptionProtocol{},
	}
}

func (c *Connection) IsEncrypted() bool {
	return c.Encryption.SharedKey != nil
}

func (c *Connection) ReadPacket() (interface{}, error) {
	return c.Reader.ReadPacket()
}

func (c *Connection) WritePacket(v interface{}) error {
	return c.Writer.WritePacket(v)
}

func (c *Connection) UpgradeWriter(f WriterFactory) {
	c.Writer.UpgradeWriter(f)
}

func (c *Connection) UpgradeReader(f ReaderFactory) {
	c.Reader.UpgradeReader(f)
}
