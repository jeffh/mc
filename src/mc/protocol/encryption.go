package protocol

import (
	"crypto/rand"
	"io"
	"mc/protocol/cfb8"
)

// Securely generates a random series of bytes of the given size.
func randomBytes(size int) ([]byte, error) {
	b := make([]byte, size)
	_, err := io.ReadFull(rand.Reader, b)
	return b, err
}

// Generates a secret key for opening encrypted connections
func GenerateSecretKey() ([]byte, error) {
	return randomBytes(16)
}

// Promotes the given connection to be encrypted.
func EncryptConnection(c *Connection) {
	key := c.Encryption.SharedKey
	c.UpgradeReader(aesCfbReader(key, key))
	c.UpgradeWriter(aesCfbWriter(key, key))
}

////////////////////////////////////////////////////////////
// golang doesn't support CFB8 from the stdlib

func aesCfbReader(key, iv []byte) ReaderFactory {
	return func(r io.Reader) io.Reader {
		return cfb8.NewReader(r, key)
	}
}

func aesCfbWriter(key, iv []byte) WriterFactory {
	return func(w io.Writer) io.Writer {
		return cfb8.NewWriter(w, key)
	}
}

/*
func aesCfbReader(key, iv []byte) ReaderFactory {
	return func(r io.Reader) io.Reader {
		block, err := aes.NewCipher(key)
		if err != nil {
			panic(err)
		}
		stream := cipher.NewCFBDecrypter(block, key)
		return &cipher.StreamReader{
			S: stream,
			R: r,
		}
	}
}

func aesCfbWriter(key, iv []byte) WriterFactory {
	return func(w io.Writer) io.Writer {
		block, err := aes.NewCipher(key)
		if err != nil {
			panic(err)
		}
		stream := cipher.NewCFBEncrypter(block, key)
		return &cipher.StreamWriter{
			S: stream,
			W: w,
		}
	}
}
*/
