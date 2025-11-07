package auth

import (
	// "strconv"
	"time"

	"github.com/eugenius-watchman/ecom_go_rest_api/config"
	"github.com/golang-jwt/jwt/v4"
)

func CreateJWT(secret []byte, userID int) (string, error) {
	expiration := time.Second * time.Duration(config.Envs.
		JWTExpirationInSeconds)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		// "userID": strconv.Itoa(userID),
		"userID": userID,  // Keep as integer
		"exp": time.Now().Add(expiration).Unix(),
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}