package types

import (
	"crypto/sha256"
)

type IHashProvider interface {
	Name() string
	Digest(source []byte) (hash []byte)
	DigestSize() int
}

var DefaultHashProvider IHashProvider = &SHA256HashProvider{}

type SHA256HashProvider struct {
}

func (self *SHA256HashProvider) Name() string {
	return "SHA256"
}

func (self *SHA256HashProvider) Digest(source []byte) (hash []byte) {
	t := sha256.Sum256(source)
	return t[:]
}

func (self *SHA256HashProvider) DigestSize() int {
	return sha256.Size
}
