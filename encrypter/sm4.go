package encrypter

import (
	"encoding/hex"
	"github.com/tjfoc/gmsm/sm4"
	"log"
)

type Sm4EcbEncrypter struct {
	Encrypter
	key []byte
}

func NewSm4EcbEncrypter(hexKey string) *Sm4EcbEncrypter {
	key, err := hex.DecodeString(hexKey)
	if err != nil {
		log.Panic(err)
	}
	return &Sm4EcbEncrypter{key: key}
}

func (e *Sm4EcbEncrypter) Encrypt(plain string) (string, error) {
	out, err := sm4.Sm4Ecb(e.key, []byte(plain), true)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(out), nil
}

func (e *Sm4EcbEncrypter) Decrypt(hexCipher string) (string, error) {
	bytes, err := hex.DecodeString(hexCipher)
	if err != nil {
		return "", err
	}

	out, err := sm4.Sm4Ecb(e.key, bytes, false)
	if err != nil {
		return "", err
	}
	return string(out), nil
}
