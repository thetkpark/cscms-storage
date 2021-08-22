package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"github.com/hashicorp/go-hclog"
	"github.com/thetkpark/cscms-temp-storage/handlers"
	"github.com/thetkpark/cscms-temp-storage/service"
	"io"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/thanhpk/randstr"
)

func main() {
	logger := hclog.Default()
	var EncryptionKey = randstr.Bytes(32)
	app := fiber.New(fiber.Config{
		BodyLimit: 150 << 20,
	})

	encryptionManager := service.NewAESEncryptionManager(logger)

	fileRouteHandler := handlers.NewFileRoutesHandler(logger, encryptionManager)

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET POST",
	}))

	app.Get("/api/ping", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success":   true,
			"timestamp": time.Now(),
		})
	})

	app.Post("/api/file", fileRouteHandler.UploadFile)

	app.Get("/api/file/:fileId", func(c *fiber.Ctx) error {
		t1 := time.Now()
		fileId := c.Params("fileId")
		if len(fileId) == 0 {
			return fiber.NewError(fiber.StatusBadRequest, "no fileId present")
		}

		encryptedFilePath := fmt.Sprintf("%s/%s.enc", "tmp", fileId)
		if _, err := os.Stat(encryptedFilePath); os.IsNotExist(err) {
			return fiber.NewError(fiber.StatusNotFound, "file not found")
		}

		ts := time.Now()
		_ = decryptFile("tmp", fileId, EncryptionKey)
		decryptDuration := time.Since(ts)

		decryptFilePath := fmt.Sprintf("%s/%s", "tmp", fileId)
		//c.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, "test.zip"))

		fmt.Sprintln(fileId)
		fmt.Printf("decrypt duration: %s\n", decryptDuration.String())
		fmt.Printf("total duration: %s\n", time.Since(t1).String())
		return c.SendFile(decryptFilePath, true)
		//return c.Send(*byteData)
	})

	app.Static("/", "./client/build")

	fmt.Println("Before anything else")
	PrintMemUsage()
	err := app.Listen(":5000")
	if err != nil {
		log.Fatalln("unable to start server", err)
	}
}

func encryptFile(filePath string, key []byte) string {
	fmt.Printf("Start encryption\n")

	fmt.Println("Before read unencrypted file")
	PrintMemUsage()
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalln("unable reading file", err)
	}
	fmt.Println("After read unencrypted file")
	PrintMemUsage()

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
	fmt.Println("Before seal")
	PrintMemUsage()
	encryptedBytes := gcm.Seal(nonce, nonce, file, nil)
	file = nil
	fmt.Println("After seal")
	PrintMemUsage()

	encryptedFilePath := fmt.Sprintf("%s.enc", filePath)
	encryptedFile, err := os.Create(encryptedFilePath)
	if err != nil {
		log.Fatalln("unable to create new file on disk", err)
	}
	defer encryptedFile.Close()

	fmt.Println("Before write file to disk")
	PrintMemUsage()
	byteWritten, err := encryptedFile.Write(encryptedBytes)
	if err != nil {
		log.Fatalln("unable to write bytes to file", err)
	}
	fmt.Println("After write file to disk")
	PrintMemUsage()

	fmt.Printf("Written %d bytes to disk\n", byteWritten)
	return encryptedFilePath
}

func decryptFile(filePath string, fileId string, key []byte) *[]byte {
	startTimestamp := time.Now()
	fmt.Println("Start decrypting")

	ciphertext, err := ioutil.ReadFile(fmt.Sprintf("%s/%s.enc", filePath, fileId))
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

	decryptedFile, err := os.Create(fmt.Sprintf("%s/%s", filePath, fileId))
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
	return &decryptedByte
}

func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
