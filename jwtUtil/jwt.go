package jwtutil

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

var mySecret = []byte("hello1213")

type Error struct {
	IsError bool   `json:"isError"`
	Message string `json:"message"`
}

type Token string

func GenerateToken(uid uint) string {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id": uid,
	})

	tokenString, err := token.SignedString(mySecret)

	if err != nil {
		log.Fatal(err)
	}
	return tokenString
}
func SetError(err Error, message string) Error {
	err.IsError = true
	err.Message = message
	return err
}
func Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Token"] == nil {
			var err Error
			err = SetError(err, "Token not found!")
			json.NewEncoder(w).Encode(err)
			return
		}

		token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("There was an error in parsing")
			}
			return mySecret, nil
		})

		if err != nil {
			json.NewEncoder(w).Encode(map[string]string{"err": err.Error()})
			return
		}

		if token.Valid {
			next.ServeHTTP(w, r)
			return
		}
	}

}
