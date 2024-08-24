package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	gorilla_mux "github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID          string
	HashedToken string
	LastIP      string
}

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

var users []User

func GenerateAccessToken(user_ip string) string {
	jwt_token := jwt.New(jwt.SigningMethodHS512)

	claims := jwt_token.Claims.(jwt.MapClaims)
	claims["ip"] = user_ip

	access_token, _ := jwt_token.SignedString([]byte("secret"))
	return access_token
}

func GenerateRefreshToken() string {
	randomBytes := make([]byte, 32)
	if _, err := rand.Read(randomBytes); err != nil {
		return ""
	}

	return base64.StdEncoding.EncodeToString(randomBytes)
}

func CheckRefreshToken() string {
	return "1"
}

func IssuingTokensHandler(w http.ResponseWriter, r *http.Request) {
	vars := gorilla_mux.Vars(r)
	user_guid := vars["guid"]
	user_ip := r.RemoteAddr

	access_token := GenerateAccessToken(user_ip)
	refresh_token := GenerateRefreshToken()

	tokens := Tokens{
		AccessToken:  access_token,
		RefreshToken: refresh_token,
	}

	hashed_token, _ := bcrypt.GenerateFromPassword([]byte(refresh_token), bcrypt.DefaultCost)
	users = append(users, User{ID: user_guid, HashedToken: string(hashed_token), LastIP: user_ip})
	json.NewEncoder(w).Encode(tokens)

}

func RefreshTokensHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// vars := gorilla_mux.Vars(r)
	// user_guid := vars["guid"]

	decoder := json.NewDecoder(r.Body)
	var user_tokens Tokens
	if err := decoder.Decode(&user_tokens); err != nil {
		panic(err)
	}

	for _, user := range users {
		if err := bcrypt.CompareHashAndPassword([]byte(user.HashedToken), []byte(user_tokens.RefreshToken)); err == nil {
			fmt.Println("SUCCESS!")
		}
	}
}

func main() {
	mux := gorilla_mux.NewRouter()

	mux.HandleFunc("/auth/{guid}", IssuingTokensHandler).Methods("GET")
	mux.HandleFunc("/refresh", RefreshTokensHandler).Methods("POST")

	http.Handle("/", mux)
	http.ListenAndServe(":8080", nil)
}
