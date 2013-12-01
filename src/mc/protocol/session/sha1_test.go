package session

import (
	. "github.com/jeffh/goexpect"
	"testing"
)

func TestSha1(t *testing.T) {
	it := NewIt(t)
	it.Expects(sha1HexDigest([]byte("Notch")), ToEqual, "4ed1f46bbe04bc756bcb17c0c7ce3e4632f06a48")
	it.Expects(sha1HexDigest([]byte("jeb_")), ToEqual, "-7c9d5b0044c130109a5d7b5fb5c317c02b4e28c1")
	it.Expects(sha1HexDigest([]byte("simon")), ToEqual, "88e16a1019277b15d58faf0541e11910eb756f6")
}
