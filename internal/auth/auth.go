package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (hashedPW string, err error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(b), err
}

func CheckPasswordHash(hash, password string) (err error) {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	tk := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.RegisteredClaims{
			Issuer:    "chirpy",
			Subject:   userID.String(),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		},
	)
	tkString, err := tk.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}
	return tkString, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	tk, err := jwt.ParseWithClaims(
		tokenString,
		jwt.MapClaims{},
		func(tk *jwt.Token) (interface{}, error) {
			if tk.Method != jwt.SigningMethodHS256 {
				return nil, fmt.Errorf("error incorrect signing method given: %v", tk.Method.Alg())
			}
			return []byte(tokenSecret), nil
		},
	)
	if err != nil {
		return uuid.Nil, err
	}
	issuer, err := tk.Claims.GetIssuer()
	if err != nil {
		return uuid.Nil, err
	}
	if issuer != string("chirpy") {
		return uuid.Nil, fmt.Errorf("invalid issuer")
	}
	idStr, err := tk.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}
	id, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid user id")
	}
	return id, nil
}
