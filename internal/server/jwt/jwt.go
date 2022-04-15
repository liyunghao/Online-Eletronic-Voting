package jwt

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"log"
	"time"

	"github.com/golang-jwt/jwt"
)

// ALGO: HS384
var JWT_Secret_KEY []byte

type JWT_AuthTokenClaims struct {
	Name  string `json:"name"`
	Group string `json:"group"`
	jwt.StandardClaims
}

func InitJWT() {
	key := make([]byte, 48)
	_, err := rand.Read(key)

	log.Println("JWT_Secret_KEY: ", base64.StdEncoding.EncodeToString(key))

	// JWT_Secret_KEY = []byte(base64.StdEncoding.EncodeToString(key))
	JWT_Secret_KEY = []byte("JWT_Secret_KEY")

	if err != nil {
		log.Fatalf("Failed to generate JWT secret key. Something WRONG: %v\n", err)
	}
}

func VerifyToken(tokenString string) (string, error) {
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return JWT_Secret_KEY, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims["name"].(string), nil
	} else {
		return "", errors.New("Invalid User")
	}
}
func GenerateToken(name string, group string) ([]byte, error) {
	token := jwt.New(jwt.GetSigningMethod("HS384"))
	token.Claims = &JWT_AuthTokenClaims{
		name,
		group,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour).Unix(),
		},
	}

	tokenString, err := token.SignedString(JWT_Secret_KEY)
	if err != nil {
		return nil, err
	}
	return []byte(tokenString), nil
}
