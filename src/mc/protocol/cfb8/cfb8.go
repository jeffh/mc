package cfb8

// Taken from toquetoes:
// https://github.com/toqueteos/minero/blob/master/util/crypto/cfb8/cfb8.go
//
// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// Based on crypto/cipher's StreamReader and StreamWriter.

import (
	"crypto/aes"
	"crypto/cipher"
	"io"
)

type Reader struct {
	RW     io.Reader
	Sr, Sw cipher.Stream
}

type Writer struct {
	RW     io.Writer
	Sr, Sw cipher.Stream
}

func NewReader(rw io.Reader, secret []byte) *Reader {
	block, _ := aes.NewCipher(secret)
	iv := secret[:block.BlockSize()]
	return &Reader{
		RW: rw,
		Sr: decrypt(block, iv),
		Sw: encrypt(block, iv),
	}
}

func NewWriter(rw io.Writer, secret []byte) *Writer {
	block, _ := aes.NewCipher(secret)
	iv := secret[:block.BlockSize()]
	return &Writer{
		RW: rw,
		Sr: decrypt(block, iv),
		Sw: encrypt(block, iv),
	}
}

func (b Reader) Read(s []byte) (n int, err error) {
	n, err = b.RW.Read(s)
	b.Sr.XORKeyStream(s[:n], s[:n])
	return
}

func (b Writer) Write(s []byte) (n int, err error) {
	d := make([]byte, len(s))
	b.Sw.XORKeyStream(d, s)
	return b.RW.Write(d)
}

type cfb8 struct {
	c         cipher.Block
	blockSize int
	iv, tmp   []byte
	de        bool
}

func encrypt(c cipher.Block, iv []byte) *cfb8 {
	cp := make([]byte, len(iv))
	copy(cp, iv)
	return &cfb8{
		c:         c,
		blockSize: c.BlockSize(),
		iv:        cp,
		tmp:       make([]byte, c.BlockSize()),
		de:        false,
	}
}

func decrypt(c cipher.Block, iv []byte) *cfb8 {
	cp := make([]byte, len(iv))
	copy(cp, iv)
	return &cfb8{
		c:         c,
		blockSize: c.BlockSize(),
		iv:        cp,
		tmp:       make([]byte, c.BlockSize()),
		de:        true,
	}
}

func (cf *cfb8) XORKeyStream(dst, src []byte) {
	for i := 0; i < len(src); i++ {
		val := src[i]
		copy(cf.tmp, cf.iv)
		cf.c.Encrypt(cf.iv, cf.iv)
		val = val ^ cf.iv[0]

		copy(cf.iv, cf.tmp[1:])
		if cf.de {
			cf.iv[15] = src[i]
		} else {
			cf.iv[15] = val
		}

		dst[i] = val
	}
}
