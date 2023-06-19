package encrypter

type Encrypter interface {
	Encrypt(plain string) (string, error)
	Decrypt(hexCipher string) (string, error)
}
