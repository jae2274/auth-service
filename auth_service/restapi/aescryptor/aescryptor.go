package aescryptor

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io"
)

type JsonAesCryptor struct {
	gcm cipher.AEAD
}

func NewJsonAesCryptor(key []byte) (*JsonAesCryptor, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return &JsonAesCryptor{
		gcm: gcm,
	}, nil
}

func (ja *JsonAesCryptor) Encrypt(obj any) (string, error) {
	data, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}

	return ja.encrypt(data)
}

func (ja *JsonAesCryptor) Decrypt(encrypted string, obj any) error {
	data, err := ja.decrypt(encrypted)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, obj)
}

func (ja *JsonAesCryptor) encrypt(data []byte) (string, error) {
	nonce := make([]byte, ja.gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := ja.gcm.Seal(nonce, nonce, data, nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func (ja *JsonAesCryptor) decrypt(encrypted string) ([]byte, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return nil, err
	}

	nonceSize := ja.gcm.NonceSize()
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return ja.gcm.Open(nil, nonce, ciphertext, nil)
}
