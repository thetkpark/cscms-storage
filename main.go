package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/google/uuid"
	"github.com/minio/sio"
	"github.com/thanhpk/randstr"
	"golang.org/x/crypto/hkdf"
)

func main() {
	var EncryptionKey = randstr.Bytes(32)
	app := fiber.New(fiber.Config{
		BodyLimit: 150 << 20,
	})

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

	app.Post("/api/file", func(c *fiber.Ctx) error {
		fmt.Println("Start")
		PrintMemUsage()
		t1 := time.Now()
		fileHeader, err := c.FormFile("file")
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "unable to get file from form-data", err.Error())
		}

		fileId, err := uuid.NewRandom()
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "unable to create new filename", err.Error())
		}

		fmt.Println("Before open file")
		PrintMemUsage()
		ts := time.Now()
		file, err := fileHeader.Open()
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "unable to open multipart file header", err.Error())
		}
		saveFileDuration := time.Now().Sub(ts)
		fmt.Println("After open file")
		PrintMemUsage()

		fmt.Println("Before encrypt")
		PrintMemUsage()
		ts = time.Now()
		encrypted := encryptFile(file)
		encryptFileDuration := time.Now().Sub(ts)
		fmt.Println("After encrypt")
		PrintMemUsage()

		// Save file to disk
		fmt.Println("Before save file")
		PrintMemUsage()
		encryptFilePath := fmt.Sprintf("%s/%s.enc", "tmp", fileId.String())
		encryptFile, err := os.Create(encryptFilePath)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "unable to create new file on disk", err.Error())
		}
		if _, err := io.Copy(encryptFile, encrypted); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "unable to copy data to file", err.Error())
		}
		fmt.Println("After save file")
		PrintMemUsage()

		return c.JSON(fiber.Map{
			"id":               fileId,
			"open_duration":    saveFileDuration.String(),
			"encrypt_duration": encryptFileDuration.String(),
			"total_time":       time.Since(t1).String(),
		})
	})

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

func encryptFile(input io.Reader) io.Reader {
	// the master key used to derive encryption keys
	// this key must be keep secret
	masterkey, err := hex.DecodeString("000102030405060708090A0B0C0D0E0FF0E0D0C0B0A090807060504030201000") // use your own key here
	if err != nil {
		log.Fatalln("unable to decode string", err)
	}

	// generate a random nonce to derive an encryption key from the master key
	// this nonce must be saved to be able to decrypt the data again - it is not
	// required to keep it secret
	var nonce [32]byte
	if _, err = io.ReadFull(rand.Reader, nonce[:]); err != nil {
		log.Fatalf("Failed to read random data: %v", err)
	}

	// derive an encryption key from the master key and the nonce
	var key [32]byte
	kdf := hkdf.New(sha256.New, masterkey, nonce[:], nil)
	if _, err = io.ReadFull(kdf, key[:]); err != nil {
		log.Fatalf("Failed to derive encryption key: %v", err)
	}

	encrypted, err := sio.EncryptReader(input, sio.Config{Key: key[:]})
	if err != nil {
		log.Fatalf("Failed to encrypted reader: %v", err)
	}

	return encrypted

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
