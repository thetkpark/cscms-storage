package service

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"github.com/hashicorp/go-hclog"
	"io"
)

type EncryptionManager interface {
	Encrypt(data *[]byte, key []byte) (*[]byte, error)
	Decrypt(ciphertext *[]byte, key []byte) (*[]byte, error)
}

type AESEncryptionManager struct {
	log hclog.Logger
}

func NewAESEncryptionManager(log hclog.Logger) *AESEncryptionManager {
	return &AESEncryptionManager{log: log}
}

func (m *AESEncryptionManager) Encrypt(data *[]byte, key []byte) (*[]byte, error) {
	// generate GCM
	gcm, err := m.generateGCM(key)
	if err != nil {
		return nil, err
	}

	// creates a new byte array the size of the nonce which must be passed to Seal
	nonce := make([]byte, gcm.NonceSize())
	// populates our nonce with a cryptographically secure random sequence
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		m.log.Error("cannot generate random nonce", err.Error())
		return nil, err
	}

	// Encrypt bytes using the Seal function
	// Seal encrypts and authenticates plaintext, authenticates the
	// additional data and appends the result to dst, returning the updated
	// slice. The nonce must be NonceSize() bytes long and unique for all
	// time, for a given key.
	ciphertext := gcm.Seal(nonce, nonce, *data, nil)
	return &ciphertext, nil
}

func (m *AESEncryptionManager) Decrypt(data *[]byte, key []byte) (*[]byte, error) {
	// generate GCM
	gcm, err := m.generateGCM(key)
	if err != nil {
		return nil, err
	}

	// Check if nonce size is longer than ciphertext
	nonceSize := gcm.NonceSize()
	if len(*data) < nonceSize {
		m.log.Error("nonce size if longer than ciphertext")
		return nil, fmt.Errorf("nonce size if longer than ciphertext")
	}

	// Decrypt the data
	var ci = *data
	nonce, ciphertext := ci[:nonceSize], ci[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		m.log.Error("cannot decrypt the ciphertext", err.Error())
		return nil, err
	}
	return &plaintext, nil
}

func (m *AESEncryptionManager) generateGCM(key []byte) (cipher.AEAD, error) {
	// generate a new aes cipher using 32 byte long key
	c, err := aes.NewCipher(key)
	if err != nil {
		m.log.Error("cannot generate new aes cipher", err.Error())
		return nil, err
	}

	// Create gcm or Galois/Counter Mode, a mode of operation for symmetric key cryptographic block ciphers
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		m.log.Error("cannot create new GCM", err.Error())
		return nil, err
	}

	return gcm, nil
}
