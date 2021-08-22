package service

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"github.com/hashicorp/go-hclog"
	"github.com/minio/sio"
	"golang.org/x/crypto/hkdf"
	"io"
)

type EncryptionManager interface {
	Encrypt(input io.Reader) (io.Reader, string, error)
	Decrypt(input io.Reader, nonceString string, output io.Writer) error
}

type SIOEncryptionManager struct {
	log       hclog.Logger
	masterKey string
}

func NewSIOEncryptionManager(l hclog.Logger, key string) *SIOEncryptionManager {
	return &SIOEncryptionManager{
		log:       l,
		masterKey: key,
	}
}

func (m *SIOEncryptionManager) Encrypt(input io.Reader) (io.Reader, string, error) {
	// the master key used to derive encryption keys
	masterkey := []byte(m.masterKey)

	// generate a random nonce to derive an encryption key from the master key
	// this nonce must be saved to be able to decrypt the data again
	var nonce [32]byte
	if _, err := io.ReadFull(rand.Reader, nonce[:]); err != nil {
		m.log.Error("Failed to read random data", err)
		return nil, "", err
	}

	// derive an encryption key from the master key and the nonce
	var key [32]byte
	kdf := hkdf.New(sha256.New, masterkey, nonce[:], nil)
	if _, err := io.ReadFull(kdf, key[:]); err != nil {
		m.log.Error("Failed to derive encryption key", err)
		return nil, "", err
	}

	encrypted, err := sio.EncryptReader(input, sio.Config{Key: key[:]})
	if err != nil {
		m.log.Error("Failed to encrypted reader", err)
		return nil, "", err
	}

	return encrypted, hex.EncodeToString(nonce[:]), nil
}

func (m *SIOEncryptionManager) Decrypt(input io.Reader, nonceString string, output io.Writer) error {
	// the master key used to derive encryption keys
	masterkey := []byte(m.masterKey)

	// the nonce used to derive the encryption key
	nonce, err := hex.DecodeString(nonceString)
	if err != nil {
		m.log.Error("Failed to decode hex string to byte", err)
		return err
	}

	// derive the encryption key from the master key and the nonce
	var key [32]byte
	kdf := hkdf.New(sha256.New, masterkey, nonce, nil)
	if _, err := io.ReadFull(kdf, key[:]); err != nil {
		m.log.Error("Failed to derive encryption key", err)
		return err
	}

	if _, err := sio.Decrypt(output, input, sio.Config{Key: key[:]}); err != nil {
		if _, ok := err.(sio.Error); ok {
			m.log.Error("Malformed encrypted data", err)
			return err
		}
		m.log.Error("Failed to decrypt data", err)
		return err
	}

	return nil
}
