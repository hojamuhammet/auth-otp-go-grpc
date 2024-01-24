package utils

import (
	"log/slog"
	"math/rand"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}

func GenerateOTP() int {
	const minOTP = 100000
	const maxOTP = 999999

	source := rand.NewSource(time.Now().UnixNano())
	randomGenerator := rand.New(source)

	otpCode := randomGenerator.Intn(maxOTP-minOTP+1) + minOTP

	return otpCode
}

func GenerateJWTToken(phoneNumber string, jwtSecret string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["phone_number"] = phoneNumber
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // Token expiration time (1 day)

	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
