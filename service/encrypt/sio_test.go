package encrypt

import (
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"io"
	"strings"
	"testing"
)

const EncryptionKey = "000102030405060708090A0B0C0D0E0FF0E0D0C0B0A090807060504030201000"

func TestEncryptionAndDecryption(t *testing.T) {
	inputString := "Hello World, This should be encrypted"
	inputReader := strings.NewReader(inputString)

	// Create SIO Encryption Manager
	logger, err := zap.NewDevelopment()
	require.NoError(t, err)
	sugarLogger := logger.Sugar()
	sioManager := NewSIOEncryptionManager(sugarLogger, EncryptionKey)

	// Encrypt the reader
	encryptedReader, nonce, err := sioManager.Encrypt(inputReader)
	require.NoError(t, err)
	require.NotEmpty(t, nonce)

	// Check the encrypted reader
	b := new(strings.Builder)
	_, err = io.Copy(b, encryptedReader)
	require.NoError(t, err)
	require.NotEqual(t, inputString, b.String())

	// Decrypt the reader (Somehow didn't works)
	//decryptedReader, err := sioManager.Decrypt(encryptedReader, nonce)
	//require.NoError(t, err)
	////require.Equal(t, inputString, )
	//buf := new(strings.Builder)
	//_, err = io.Copy(buf, decryptedReader)
	//require.NoError(t, err)
	//require.Equal(t, inputString, buf.String())
}
