package utils

import "crypto/sha256"

func CreateHash(key string) []byte {
	return createHashWithLen(key, 32)
}

func createHashWithLen(key string, length int) []byte {
	hasher := sha256.New()
	hasher.Write([]byte(key))
	return hasher.Sum(nil)[:length]
}
