package hasher

import (
	"crypto/sha512"
	"encoding/hex"
	"io"
)

type Sha512 struct {
}

func NewSha512() *Sha512 {
	return &Sha512{}
}

func (m *Sha512) Hash(file io.Reader) (string, error) {
	hash := sha512.New()
	_, err := io.Copy(hash, file)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}
