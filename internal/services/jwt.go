package services

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

var JwtSigningKey = []byte(os.Getenv("JWT_SIGNING_KEY"))

func GenerateJWT(userName string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["exp"] = time.Now().Add(24 * time.Hour).Unix()
	claims["authorized"] = true
	claims["user"] = userName

	return token.SignedString(JwtSigningKey)
}

func ParseJWT(t string) (string, error) {
	token, err := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return JwtSigningKey, nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims["user"].(string), nil
	} else {
		return "", fmt.Errorf("failed to get jwt claims")
	}
}
