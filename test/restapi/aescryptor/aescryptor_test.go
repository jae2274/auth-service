package aescryptor

import (
	"testing"
	"time"
	"userService/usersvc/restapi/aescryptor"
	"userService/usersvc/utils"

	"github.com/stretchr/testify/require"
)

func TestJsonAesCryptor(t *testing.T) {
	type sampleStruct struct {
		Name           string    `json:"name"`
		Age            int       `json:"age"`
		Birth          time.Time `json:"birth"`
		WillBeIgnored  string    `json:"-"`
		willBeIgnored2 string    `json:"willBeIgnored2"`
	}

	sample := &sampleStruct{
		Name:           "John Doe",
		Age:            30,
		Birth:          time.Now(),
		WillBeIgnored:  "This will be ignored",
		willBeIgnored2: "This will be ignored too",
	}

	t.Run("encrypt and decrypt", func(t *testing.T) {
		key := utils.CreateHash("testkey")
		cryptor, err := aescryptor.NewJsonAesCryptor(key)
		require.NoError(t, err)

		encrypted, err := cryptor.Encrypt(sample)
		require.NoError(t, err)

		decrypted := &sampleStruct{}
		err = cryptor.Decrypt(encrypted, decrypted)
		require.NoError(t, err)

		require.Equal(t, sample.Name, decrypted.Name)
		require.Equal(t, sample.Age, decrypted.Age)
		require.Equal(t, sample.Birth.Format(time.RFC3339), decrypted.Birth.Format(time.RFC3339))
		require.Empty(t, decrypted.WillBeIgnored)
		require.Empty(t, decrypted.willBeIgnored2)
	})

	t.Run("encrypt and decrypt with different key", func(t *testing.T) {

		key := utils.CreateHash("testkey")
		cryptor, err := aescryptor.NewJsonAesCryptor(key)
		require.NoError(t, err)

		encrypted, err := cryptor.Encrypt(sample)
		require.NoError(t, err)

		anotherKey := utils.CreateHash("anotherkey")
		anotherCryptor, err := aescryptor.NewJsonAesCryptor(anotherKey)
		require.NoError(t, err)

		decrypted := &sampleStruct{}
		err = anotherCryptor.Decrypt(encrypted, decrypted)
		require.Error(t, err)
	})
}
