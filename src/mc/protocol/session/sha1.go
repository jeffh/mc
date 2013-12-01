package session

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"strings"
)

func sha1HexDigest(data []byte) string {
	hasher := sha1.New()
	hasher.Write(data)
	hash := hasher.Sum(nil)

	negative := (hash[0] & 0x80) > 0

	if negative {
		carry := true
		for i := len(hash) - 1; i >= 0; i-- {
			hash[i] = ^hash[i]
			if carry {
				carry = (hash[i] == 0xff)
				hash[i]++
			}
		}
	}

	result := base64.URLEncoding.EncodeToString(hash)
	result = strings.TrimLeft(fmt.Sprintf("%x", hash), "0")

	if negative {
		return "-" + result
	}
	return result
}
