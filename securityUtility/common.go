package securityutility

import (
	"fmt"
	"net/http"
	"strings" //  "github.com/gorilla/mux"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"github.com/heroeexceso/golang/httputility"
	"github.com/mitchellh/mapstructure"
)

// User ... estructura de credenciales
type User struct {
	UserName     string `json:"username"`
	UserPassword string `json:"userpassword"`
}

/*
type JwtToken struct {
	Toke string `json:"token"`
}
*/

// validateMiddleware ... algo
func validateMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		authorizationHeader := httputility.GetHeaderKey(req, "authorization")
		if authorizationHeader != "" {
			bearerToken := strings.Split(authorizationHeader, " ")
			if len(bearerToken) == 2 {
				token, error := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("There was an error")
					}
					return []byte("secret"), nil
				})
				if error != nil {
					httputility.GetJsonResponseMessage(w, error.Error())
					return
				}
				if token.Valid {
					context.Set(req, "decoded", token.Claims)
					next(w, req)
				} else {
					httputility.GetJsonResponseMessage(w, "Invalid authorization token")
				}
			}
		} else {
			httputility.GetJsonResponseMessage(w, "An authorization header is required")
		}
	})
}

// GetTokenString ... algo
func GetTokenString(userName string, userPassword string) (string, bool) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": userName,
		"password": userPassword,
	})

	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		return "", false
	}

	return tokenString, true
}

// PostProtected ... algo
func PostProtected(tokenString string) (string, string, bool) {
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Error!")
		}
		return []byte("secret"), nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		var user User

		mapstructure.Decode(claims, &user)

		return user.UserName, user.UserPassword, true
	} else {
		return "", "", false
	}
}
