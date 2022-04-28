package hasher

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
)

type Sha256 struct {
}

func NewSha256() *Sha256 {
	return &Sha256{}
}

func (m *Sha256) Hash(file io.Reader) (string, error) {
	hash := sha256.New()
	_, err := io.Copy(hash, file)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}
