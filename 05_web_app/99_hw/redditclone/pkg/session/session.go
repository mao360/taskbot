package session

import (
	"context"
	"errors"
	"github.com/dgrijalva/jwt-go"
	redditclone "reddit_clone"
)

type sessKey string

var Key sessKey = "session"
var secretKey = []byte("ueFGhgui7T8OLg8")

var (
	errBadSignMethod = errors.New("bad sign method")
	errNoPayload     = errors.New("no payload")
	errBadToken      = errors.New("bad token")
)

func CreateSess(user *redditclone.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":   user.UserID,
		"password": user.Password,
		"username": user.UserName,
	})
	tokenStr, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}

func CheckSess(inToken string) (map[string]interface{}, error) {
	hashSecretGetter := func(token *jwt.Token) (interface{}, error) {
		method, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok || method.Alg() != "HS256" {
			return nil, errBadSignMethod
		}
		return secretKey, nil
	}
	token, err := jwt.Parse(inToken, hashSecretGetter)
	if err != nil || !token.Valid {
		return nil, errBadToken
	}
	payload, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errNoPayload
	}

	return map[string]interface{}{
		"userID":   payload["userID"],
		"password": payload["password"],
		"username": payload["username"],
	}, nil
}

func SessFromContext(ctx context.Context) (map[string]interface{}, error) {
	userMap, ok := ctx.Value(Key).(map[string]interface{})
	if !ok || userMap == nil {
		return nil, errBadToken
	}

	return userMap, nil
}
