package mc

import (
	"bytes"
	"crypto/rand"
	"crypto/x509"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"math/big"
	"mc/protocol"
)

type SecureConnection struct {
	PlainInbox    chan interface{}
	PlainOutbox   chan interface{}
	SecuredInbox  chan interface{}
	SecuredOutbox chan interface{}
	Cert          *x509.Certificate
	ServerID      string
	SharedKey     []byte
	RandReader    io.Reader
}

func NewSecureConnection(plainInbox, plainOutbox chan interface{}) *SecureConnection {
	return &SecureConnection{
		PlainInbox:    plainInbox,
		PlainOutbox:   plainOutbox,
		SecuredInbox:  make(chan interface{}, 20),
		SecuredOutbox: make(chan interface{}, 20),
		RandReader:    rand.Reader,
	}
}

func (c *SecureConnection) generateInt() (int64, error) {
	max := big.NewInt(math.MaxInt64)
	bi, err := rand.Int(rand.Reader, max)
	return bi.Int64(), err
}

func (c *SecureConnection) generateKey() ([]byte, error) {
	sharedKey1, err := c.generateInt()
	if err != nil {
		return nil, err
	}
	sharedKey2, err := c.generateInt()
	if err != nil {
		return nil, err
	}
	b := bytes.NewBuffer([]byte{})
	err = binary.Write(b, binary.BigEndian, sharedKey1)
	if err != nil {
		return nil, err
	}
	err = binary.Write(b, binary.BigEndian, sharedKey2)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func (c *SecureConnection) Connect(hostname string, port int32, username string) error {
	c.PlainInbox <- protocol.Handshake{
		Version:  47,
		Username: username,
		Hostname: hostname,
		Port:     port,
	}
	packet := <-c.PlainOutbox
	p, ok := packet.(protocol.EncryptionKeyRequest)
	if !ok {
		return fmt.Errorf("Expected encryption key request, got: %v", packet)
	}
	c.ServerID = p.ServerID
	cert, err := x509.ParseCertificate(p.PublicKey)
	if err != nil {
		return fmt.Errorf("Failed to parse server cert: %s", err)
	} else {
		c.Cert = cert
	}

	// now we must generate random numbers
	c.SharedKey, err = c.generateKey()

	return nil
}

func (c *SecureConnection) ProcessInbox() {
	for {
		_ = <-c.SecuredInbox
	}
}
