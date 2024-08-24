package services

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/MuhaFAH/AuthCheck/pkg/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"
)

func GenerateAccessToken(userIp string) (string, error) {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error when loading .env file: %s", err.Error())
	}
	validHours, err := strconv.Atoi(os.Getenv("ACCESS_TOKEN_VALID_HOURS"))
	if err != nil {
		return "", err
	}

	jwtToken := jwt.New(jwt.SigningMethodHS512)
	claims := jwtToken.Claims.(jwt.MapClaims)
	claims["ip"] = userIp
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(validHours)).Unix()

	accessToken, err := jwtToken.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func GenerateRefreshToken() (string, error) {
	randomBytes := make([]byte, 32)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(randomBytes), nil
}

func GetNewTokens(userIp string) (models.Tokens, string, error) {
	accessToken, err := GenerateAccessToken(userIp)
	if err != nil {
		return models.Tokens{}, "", err
	}
	refreshToken, err := GenerateRefreshToken()
	if err != nil {
		return models.Tokens{}, "", err
	}

	tokens := models.Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	hashedToken, _ := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)

	return tokens, string(hashedToken), nil
}

func SendResponse(writer http.ResponseWriter, response models.Response) error {
	err := json.NewEncoder(writer).Encode(response)
	if err != nil {
		return err
	}
	return nil
}

func SendErrorResponse(writer http.ResponseWriter, errorMessage string) error {
	err := SendResponse(writer, models.Response{Answer: "access denied", Reason: errorMessage, HTTPStatus: 403})
	if err != nil {
		return err
	}
	return nil
}

func GetRequestUserInfo(r *http.Request) (string, string, error) {
	vars := mux.Vars(r)
	userGUID := vars["guid"]
	userIp, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", "", err
	}
	return userGUID, userIp, nil
}

func SendEmail(mail models.Email) error {
	filename := fmt.Sprintf("email/email_%s.txt", time.Now().Format("20060102_150405"))
	err := os.MkdirAll("email", os.ModePerm)
	if err != nil {
		return err
	}
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	fullEmailText := fmt.Sprintf("От: %s\nДля: %s\n\nЗдравствуйте!\nПри последней выдаче нового токена, ваш IP не совпадал с тем, что был у Вас ранее. Спасибо!", mail.From, mail.To)
	if _, err := file.WriteString(fullEmailText); err != nil {
		return err
	}
	return nil
}
