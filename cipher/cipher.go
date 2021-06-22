package cipher

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type CustomClaims struct {
	Message string `json:"message"`
	jwt.StandardClaims
}

type Cipher struct {
	TokenExpiredHours int
	HashCost          int
	PassPhrase        string
}

func DefaultCipher() Cipher {
	// PrivateKey used to encrypt and decrypt, lenght must be 32
	return Cipher{
		TokenExpiredHours: 24,
		HashCost:          11,
		PassPhrase:        "default passphrase",
	}
}

// NewCipher new a cihper
func NewCipher(tokenExpiredHours, hashCost int, passPhrase string) Cipher {
	// PrivateKey used to encrypt and decrypt, lenght must be 32
	return Cipher{
		TokenExpiredHours: tokenExpiredHours,
		HashCost:          hashCost,
		PassPhrase:        passPhrase,
	}
}

// EncodePassword encode plain text password
func (c *Cipher) EncodePassword(plaintext string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(plaintext), c.HashCost)
	return string(bytes)
}

// ComparePassword compare the cipher text password is correct or not
func (c *Cipher) ComparePassword(ciphertext string, plaintext string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(ciphertext), []byte(plaintext))
	return err == nil
}

// GetJWT convert message string to token string
func (c *Cipher) GetJWT(message string) string {
	claims := CustomClaims{
		message,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(c.TokenExpiredHours) * time.Hour).Unix(), // 1 day
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(c.PassPhrase))
	return tokenString
}

// ParseJWT convert token string to message string
func (c *Cipher) ParseJWT(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(c.PassPhrase), nil
	})
	if token != nil {
		if token.Valid {
			if claims, ok := token.Claims.(*CustomClaims); ok {
				return claims.Message, nil
			} else {
				return "", errors.New("invalid token")
			}
		} else {
			return "", err
		}
	} else {
		return "", errors.New("invalid token")
	}

}
