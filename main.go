package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"
)

const EncryptionKey = "@i90$5NEWTpF@%rSZlovn@CQETD2FbA2"

func main() {
	app := fiber.New(fiber.Config{
		BodyLimit: 150 << 20,
	})

	app.Get("/api/ping", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success":   true,
			"timestamp": time.Now(),
		})
	})

	app.Post("/api/upload", func(c *fiber.Ctx) error {
		fileHeader, err := c.FormFile("file")
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "unable to get file from form-data", err.Error())
		}

		fileId, err := uuid.NewRandom()
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "unable to create new filename", err.Error())
		}
		unencryptedFilePath := fmt.Sprintf("%s/%s", "tmp", fileId.String())
		ts := time.Now()
		err = c.SaveFile(fileHeader, unencryptedFilePath)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "unable to save unencrypt file to disk", err.Error())
		}
		saveFileDuration := time.Now().Sub(ts)

		// Encrypt the file
		ts = time.Now()
		encryptFilePath := encryptFile(unencryptedFilePath)
		encryptFileDuration := time.Now().Sub(ts)

		//return c.
		return c.JSON(fiber.Map{
			"path":             encryptFilePath,
			"save_duration":    saveFileDuration.String(),
			"encrypt_duration": encryptFileDuration.String(),
		})
	})

	err := app.Listen(":5000")
	if err != nil {
		log.Fatalln("unable to start server", err)
	}
}

func encryptFile(filePath string) string {
	fmt.Printf("Start encryption")

	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalln("unable reading file", err)
	}
	key := []byte(EncryptionKey)

	// generate a new aes cipher using our 32 byte long key
	c, err := aes.NewCipher(key)
	// if there are any errors, handle them
	if err != nil {
		log.Fatalln("unable to create new cipher", err)
	}

	// gcm or Galois/Counter Mode, is a mode of operation
	// for symmetric key cryptographic block ciphers
	gcm, err := cipher.NewGCM(c)
	// if any error generating new GCM
	// handle them
	if err != nil {
		log.Fatalln("unable to create new GCM", err)
	}

	// creates a new byte array the size of the nonce
	// which must be passed to Seal
	nonce := make([]byte, gcm.NonceSize())
	// populates our nonce with a cryptographically secure
	// random sequence
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		log.Fatalln("unable to generate nonce", err)
	}

	// here we encrypt our text using the Seal function
	// Seal encrypts and authenticates plaintext, authenticates the
	// additional data and appends the result to dst, returning the updated
	// slice. The nonce must be NonceSize() bytes long and unique for all
	// time, for a given key.
	encryptedBytes := gcm.Seal(nonce, nonce, file, nil)

	encryptedFilePath := fmt.Sprintf("%s.enc", filePath)
	encryptedFile, err := os.Create(encryptedFilePath)
	if err != nil {
		log.Fatalln("unable to create new file on disk", err)
	}
	defer encryptedFile.Close()

	byteWritten, err := encryptedFile.Write(encryptedBytes)
	if err != nil {
		log.Fatalln("unable to write bytes to file", err)
	}

	fmt.Printf("Written %d bytes to disk\n", byteWritten)
	return encryptedFilePath
}

func decryptFile() {
	startTimestamp := time.Now()
	fmt.Println("Start decrypting")

	key := []byte(EncryptionKey)
	ciphertext, err := ioutil.ReadFile("encrypted")
	// if our program was unable to read the file
	// print out the reason why it can't
	if err != nil {
		log.Fatalln("unable to read encrypted file", err)
	}

	c, err := aes.NewCipher(key)
	if err != nil {
		log.Fatalln("unable to create new cipher", err)
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		log.Fatalln("unable to create new gcm", err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		log.Fatalln("ciphertext less than nonce size")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	decryptedByte, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		log.Fatalln("unable to decrypt the file", err)
	}

	decryptedFile, err := os.Create("decrypt.zip")
	if err != nil {
		log.Fatalln("unable to create new file on disk", err)
	}
	defer decryptedFile.Close()

	byteWritten, err := decryptedFile.Write(decryptedByte)
	if err != nil {
		log.Fatalln("unable to write bytes to file", err)
	}

	fmt.Printf("Written %d bytes to disk\n", byteWritten)
	endTimestamp := time.Now()
	fmt.Printf("Time used %v ms\n", endTimestamp.Sub(startTimestamp).Milliseconds())
}
