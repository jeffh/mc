package protocol

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

func RandomBytes(size int) ([]byte, error) {
	b := make([]byte, size)
	_, err := io.ReadFull(rand.Reader, b)
	return b, err
}

func GenerateSecretKey() ([]byte, error) {
	return RandomBytes(16)
}

func EncryptConnection(c *Connection) {
	key := c.Encryption.SharedKey
	c.UpgradeReader(AesCfbReader(key, key))
	c.UpgradeWriter(AesCfbWriter(key, key))
}

////////////////////////////////////////////////////////////

func AesCfbReader(key, iv []byte) ReaderFactory {
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

func AesCfbWriter(key, iv []byte) WriterFactory {
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
