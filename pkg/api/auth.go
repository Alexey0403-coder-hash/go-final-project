package api

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

type Password struct {
	Password string `json:"password"`
}

type JWT struct {
	Token string `json:"token"`
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	var passwordUser Password

	err := json.NewDecoder(r.Body).Decode(&passwordUser)
	if err != nil {
		writeJSONError(w, "ошибка чтения из JSON", http.StatusBadRequest)
		return
	}

	passwordEnv := os.Getenv("TODO_PASSWORD")
	if passwordEnv == "" {
		writeJSONError(w, "аутентификация отключена", http.StatusForbidden)
		return
	}

	if passwordUser.Password != passwordEnv {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		errResponse := ErrResponse{
			Error: "неправильный пароль",
		}
		err := json.NewEncoder(w).Encode(errResponse)
		if err != nil {
			log.Printf("ошибка отправки JSON: %v", err)
		}
		return
	}

	hash := sha256.Sum256([]byte(passwordEnv))
	hashString := fmt.Sprintf("%x", hash)

	claims := jwt.MapClaims{
		"hash": hashString,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(passwordEnv))
	if err != nil {
		log.Printf("ошибка подписи токена: %v", err)
		writeJSONError(w, "ошибка на стороне сервера", http.StatusInternalServerError)
		return
	}

	jwtResponse := JWT{
		Token: tokenString,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(jwtResponse)
	if err != nil {
		log.Printf("ошибка отправки JSON: %v", err)
		return
	}
}

func auth(next http.HandlerFunc) http.HandlerFunc {
	password := os.Getenv("TODO_PASSWORD")

	return func(w http.ResponseWriter, r *http.Request) {

		if password == "" {
			next(w, r)
			return
		}

		cookie, err := r.Cookie("token")
		if err != nil || cookie.Value == "" {
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(
			cookie.Value,
			func(token *jwt.Token) (interface{}, error) {
				return []byte(password), nil
			},
			jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		)

		if err != nil || !token.Valid {
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}

		hash := sha256.Sum256([]byte(password))
		hashString := fmt.Sprintf("%x", hash)

		tokenHash, ok := claims["hash"].(string)
		if !ok || tokenHash != hashString {
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}
