package hasher

import (
	"crypto/sha256"
	"fmt"
	"strings"

	"golang.org/x/crypto/pbkdf2"
)

type PBKDF2Hasher struct {
	alg        string
	iterations int
	secret     []byte
}

func NewPBKDF2Hasher(iterations int, secret []byte) *PBKDF2Hasher {
	return &PBKDF2Hasher{
		alg:        "pbkdf2",
		iterations: iterations,
		secret:     secret,
	}
}

func (h *PBKDF2Hasher) GenerateFromPassword(password string) (string, error) {
	key := pbkdf2.Key([]byte(password), []byte(h.secret), h.iterations, sha256.Size, sha256.New)
	stringPass := fmt.Sprintf("%s$%d$%x$%x", h.alg, h.iterations, h.secret, key)
	fmt.Println(stringPass)
	return stringPass, nil
}

func (h *PBKDF2Hasher) CompareHashAndPassword(hashedPassword, password string) error {
	hashedPasswordSplit := strings.Split(hashedPassword, "$")
	if len(hashedPasswordSplit) != 4 {
		return fmt.Errorf("invalid hashed password")
	}
	iterations := hashedPasswordSplit[1]
	secret := hashedPasswordSplit[2]
	key := pbkdf2.Key([]byte(password), []byte(secret), h.iterations, sha256.Size, sha256.New)
	newHashedPassword := fmt.Sprintf("%s$%s$%s$%x", h.alg, iterations, secret, key)
	if hashedPassword != newHashedPassword {
		return fmt.Errorf("passwords do not match")
	}
	return nil
}
