package services

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
)

type Crypt struct {
	key []byte
}

type CryptInterface interface {
	GenerateRand() (string, error)
	Hash(content string) string
	IsEqual(content string, compareSum string) bool
}

func NewCrypt(key string) CryptInterface {
	return Crypt{
		key: []byte(key),
	}
}

func (c Crypt) GenerateRand() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return ``, err
	}

	return hex.EncodeToString(b), nil
}

func (c Crypt) Hash(content string) string {
	sign := hmac.New(sha256.New, c.key)
	sign.Write([]byte(content))
	sum := sign.Sum(nil)

	return hex.EncodeToString(sum)
}

func (c Crypt) IsEqual(content string, compareSum string) bool {
	sign := hmac.New(sha256.New, c.key)
	sign.Write([]byte(content))
	sum := sign.Sum(nil)

	compareByte, _ := hex.DecodeString(compareSum)

	return hmac.Equal(sum, compareByte)
}
