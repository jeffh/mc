// This package handles the low-level details of how the minecraft protocol
// operations on the socket level. This subpackage provides the following features:
//
//    - Parsing the latest minecraft protocol as documented at: http://mc.kev009.com/Protocol
//    - Provides structs for all the types of messages the minecraft protocol utilizes
//    - Can (de)serialize minecraft protocol bytes into structs
//
//
// mc/protocol is extremely low-level and probably should not be used outside of
// the mc package, unless you understand the protocol directly.
//
// While this package handles the serialization to an io.Reader or io.Writer, it does
// NOT automatically handle higher-level connection or parsing concerns such as:
//
//    - Regularly sending keep alives
//    - Promoting connections to be fully encrypted
//    - minecraft server session authentication (http://www.wiki.vg/Session)
//    - Parsing NBT or Chunk data
//
// The majority of types in this package deal with the various minecraft
// messages (aka, packets) that can be serialized.
//
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
	}
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
	encSecret, err := c.Encryption.encrypt(secret)
	if err != nil {
		return
	}

	encToken, err := c.Encryption.encrypt(ekReq.VerifyToken)
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

// The struct that holds the information for opening an encrypted connection.
// This is part of the Connection struct.
//
// Use connection.IsEncrypted() to check if this struct is used or not.
type EncryptionProtocol struct {
	PublicKey interface{}
	SharedKey []byte
}

// Internal. Encrypts the given bytes.
func (e *EncryptionProtocol) encrypt(d []byte) ([]byte, error) {
	pk, ok := e.PublicKey.(*rsa.PublicKey)
	if ok {
		return rsa.EncryptPKCS1v15(rand.Reader, pk, d)
	}
	return nil, fmt.Errorf("Unknown PublicKey: %#v", e.PublicKey)
}

//////////////////////////////////////////////////////////

// A function that creates an io.Writer from another io.Writer.
// This is used to promote connections from plaintext to encrypted.
type WriterFactory func(w io.Writer) io.Writer

// A function that creates an io.Reader from another io.Reader.
// This is used to promote connections from plaintext to encrypted.
type ReaderFactory func(w io.Reader) io.Reader

// The interface that wraps writing to a socket.
//
// It can optionally be given a WriterFactory to promote the internally
// used io.Writer to be encrypted after the encryption handshake
// has been completed.
type WritePacketer interface {
	WritePacket(v interface{}) error
	UpgradeWriter(WriterFactory)
}

// The interface that wraps reading from a socket.
//
// It can optionally be given a ReaderFactory to promote the internally
// used io.Reader to be encrypted after the encryption handshake.
type ReadPacketer interface {
	ReadPacket() (interface{}, error)
	UpgradeReader(ReaderFactory)
}

// Represents the minecraft connection. It allows consumers of this type
// to send and receive packets
type Connection struct {
	Writer     WritePacketer
	Reader     ReadPacketer
	Encryption EncryptionProtocol
	ServerID   string
}

// Creates a new consumer-level connection from the given packet readers
// and writers.
func NewConnection(reader ReadPacketer, writer WritePacketer) *Connection {
	return &Connection{
		Writer:     writer,
		Reader:     reader,
		Encryption: EncryptionProtocol{},
	}
}

// Returns true if the current connection is encrypted.
func (c *Connection) IsEncrypted() bool {
	return c.Encryption.SharedKey != nil
}

// Reads a minecraft packet off the connection.
func (c *Connection) ReadPacket() (interface{}, error) {
	return c.Reader.ReadPacket()
}

// Writes a minecraft packet to the connection
func (c *Connection) WritePacket(v interface{}) error {
	return c.Writer.WritePacket(v)
}

// You can pass a WriterFactory function to wrap the io.Writer mid-connection.
// This is used by EncryptConnection to encrypt an open connection.
func (c *Connection) UpgradeWriter(f WriterFactory) {
	c.Writer.UpgradeWriter(f)
}

// You can pass a ReaderFactory function to wrap the io.Writer mid-connection.
// This is used by EncryptConnection to encrypt an open connection.
func (c *Connection) UpgradeReader(f ReaderFactory) {
	c.Reader.UpgradeReader(f)
}
