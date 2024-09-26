package config

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"log"

	"github.com/thrasher-corp/gocryptotrader/common"
	"github.com/thrasher-corp/gocryptotrader/common/crypto"
	"golang.org/x/crypto/scrypt"
)

const (
	// SaltRandomLength is the number of random bytes to append after the prefix string
	SaltRandomLength = 12
)

var (
	errAESBlockSize = errors.New("config file data is too small for the AES required block size")
	errNoPrefix     = errors.New("data does not start with Encyption Prefix")

	// EncryptConfirmPrefix is a prefix to tell us the file is encrypted
	EncryptConfirmPrefix = []byte("THORS-HAMMER")
	SaltPrefix           = []byte("~GCT~SO~SALTY~")
)

// promptForConfigEncryption asks for encryption confirmation
// returns true if encryption was desired, false otherwise
func promptForConfigEncryption() (bool, error) {
	log.Println("Would you like to encrypt your config file (y/n)?")

	input := ""
	if _, err := fmt.Scanln(&input); err != nil {
		return false, err
	}

	return common.YesOrNo(input), nil
}

// Unencrypted provides the default key provider implementation for unencrypted files
func Unencrypted() ([]byte, error) {
	return nil, errors.New("encryption key was requested, no key provided")
}

// PromptForConfigKey asks for configuration key
// if initialSetup is true, the password needs to be repeated
func PromptForConfigKey(initialSetup bool) ([]byte, error) {
	var cryptoKey []byte

	for {
		log.Println("Please enter in your password: ")
		pwPrompt := func(i *[]byte) error {
			_, err := fmt.Scanln(i)
			return err
		}

		var p1 []byte
		err := pwPrompt(&p1)
		if err != nil {
			return nil, err
		}

		if !initialSetup {
			cryptoKey = p1
			break
		}

		var p2 []byte
		log.Println("Please re-enter your password: ")
		err = pwPrompt(&p2)
		if err != nil {
			return nil, err
		}

		if bytes.Equal(p1, p2) {
			cryptoKey = p1
			break
		}
		log.Println("Passwords did not match, please try again.")
	}
	return cryptoKey, nil
}

// EncryptConfigFile encrypts configuration data that is parsed in with a key
// and returns it as a byte array with an error
func EncryptConfigFile(configData, key []byte) ([]byte, error) {
	sessionDK, salt, err := makeNewSessionDK(key)
	if err != nil {
		return nil, err
	}
	c := &Config{
		sessionDK:  sessionDK,
		storedSalt: salt,
	}
	return c.encryptConfigFile(configData)
}

// encryptConfigFile encrypts configuration data that is parsed in with a key
// and returns it as a byte array with an error
func (c *Config) encryptConfigFile(configData []byte) ([]byte, error) {
	block, err := aes.NewCipher(c.sessionDK)
	if err != nil {
		return nil, err
	}

	ciphertext := make([]byte, aes.BlockSize+len(configData))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], configData)

	appendedFile := append(bytes.Clone(EncryptConfirmPrefix), c.storedSalt...)
	appendedFile = append(appendedFile, ciphertext...)
	return appendedFile, nil
}

// DecryptConfigFile decrypts configuration data with the supplied key and
// returns the un-encrypted data as a byte array with an error
func DecryptConfigFile(d, key []byte) ([]byte, error) {
	return (&Config{}).decryptConfigData(d, key)
}

// decryptConfigData decrypts configuration data with the supplied key and
// returns the un-encrypted data as a byte array with an error
func (c *Config) decryptConfigData(d []byte, key []byte) ([]byte, error) {
	if !bytes.HasPrefix(d, EncryptConfirmPrefix) {
		return d, errNoPrefix
	}

	d = bytes.TrimPrefix(d, EncryptConfirmPrefix)

	sessionDK, storedSalt, err := makeNewSessionDK(key)
	if err != nil {
		return nil, err
	}

	if bytes.HasPrefix(d, SaltPrefix) {
		salt := make([]byte, len(SaltPrefix)+SaltRandomLength)
		salt = d[0:len(salt)]

		var err error
		key, err = getScryptDK(key, salt)
		if err != nil {
			return nil, err
		}

		d = d[len(salt):]
	}

	blockDecrypt, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(d) < aes.BlockSize {
		return nil, errAESBlockSize
	}

	iv, d := d[:aes.BlockSize], d[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(blockDecrypt, iv)
	stream.XORKeyStream(d, d)

	c.sessionDK, c.storedSalt = sessionDK, storedSalt

	return d, nil
}

func getScryptDK(key, salt []byte) ([]byte, error) {
	if len(key) == 0 {
		return nil, errors.New("key is empty")
	}
	return scrypt.Key(key, salt, 32768, 8, 1, 32)
}

func makeNewSessionDK(key []byte) (dk, storedSalt []byte, err error) {
	storedSalt, err = crypto.GetRandomSalt([]byte(SaltPrefix), SaltRandomLength)
	if err != nil {
		return nil, nil, err
	}

	dk, err = getScryptDK(key, storedSalt)
	if err != nil {
		return nil, nil, err
	}

	return dk, storedSalt, nil
}
