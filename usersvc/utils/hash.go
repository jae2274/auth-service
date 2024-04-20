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

func TinyIntToBool(tinyInt int8) bool {
	return tinyInt != 0
}

func BoolToTinyInt(b bool) int8 {
	if b {
		return 1
	}
	return 0
}
