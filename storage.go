package main // import "dirba.io/micropass"
import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"golang.org/x/crypto/scrypt"
)

/*
Databse format:
Offset Len Explanation
  0x00   4 Magic code "upss"
  0x04   4 Version number, little endian (currently 1)
  0x08   * Encrypted data
*/

func databaseLoad(filename, password string) (*Database, error) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	if len(file) < 52 {
		return nil, fmt.Errorf("file is too short")
	}

	magic := file[0:8]
	if !bytes.Equal(magic, []byte("upss\x01\x00\x00\x00")) {
		return nil, fmt.Errorf("invalid magic")
	}
	salt := file[8:40]
	nonce := file[40:52]
	ciphertext := file[52:]

	key, err := scrypt.Key([]byte(password), salt, 524288, 8, 1, 32)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	var db *Database
	json.Unmarshal(plaintext, db)
	return db, nil
}

func databaseSave(filename string, password string, db *Database) error {
	// Generate random salt
	salt := make([]byte, 32)
	_, err := rand.Read(salt)
	if err != nil {
		return err
	}

	key, err := scrypt.Key([]byte(password), salt, 524288, 8, 1, 32)
	if err != nil {
		return err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	nonce := make([]byte, 12)
	_, err = rand.Read(nonce)
	if err != nil {
		return err
	}

	plaintext, err := json.Marshal(db)
	if err != nil {
		return err
	}

	file := []byte("upss\x01\x00\x00\x00")
	file = append(file, salt...)
	file = append(file, nonce...)
	file = aesgcm.Seal(file, nonce, plaintext, nil)

	err = ioutil.WriteFile(filename, file, 0644)
	if err != nil {
		return err
	}

	return nil
}
