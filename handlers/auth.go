package handlers

import (
	"encoding/json"
	"github.com/MuhaFAH/AuthCheck/pkg/models"
	"github.com/MuhaFAH/AuthCheck/pkg/storage"
	"github.com/MuhaFAH/AuthCheck/services"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"os"
)

type App struct {
	DB *sqlx.DB
}

func (app *App) IssuingTokensHandler(w http.ResponseWriter, r *http.Request) {
	userGUID, userIP, err := services.GetRequestUserInfo(r)
	if err != nil {
		log.Printf("error when getting user info: %s", err.Error())
	}

	tokens, hashedToken, err := services.GetNewTokens(userIP)
	if err != nil {
		log.Fatalf("error when generating tokens: %s", err.Error())
		return
	}

	err = storage.AddUser(models.User{GUID: userGUID, LastIP: userIP, HashedToken: hashedToken}, app.DB)
	if err != nil {
		if err := services.SendErrorResponse(w, "invalid user guid"); err != nil {
			log.Printf("error when sending error response: %s", err.Error())
		}
		return
	}
	if err := services.SendResponse(w, models.Response{Answer: "OK", HTTPStatus: 200, Tokens: tokens}); err != nil {
		log.Printf("error when sending issuing-response: %s", err.Error())
	}

}

func (app *App) RefreshTokensHandler(w http.ResponseWriter, r *http.Request) {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error when loading .env file: %s", err.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	userGUID, userIP, err := services.GetRequestUserInfo(r)
	if err != nil {
		log.Printf("error when getting user info: %s", err.Error())
	}

	decoder := json.NewDecoder(r.Body)
	var userTokens models.Tokens
	if err := decoder.Decode(&userTokens); err != nil {
		panic(err)
	}
	user := models.User{GUID: userGUID}
	if err := storage.GetUser(&user, app.DB); err != nil {
		if err := services.SendErrorResponse(w, "invalid user guid"); err != nil {
			log.Printf("error when sending response: %s", err.Error())
		}
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedToken), []byte(userTokens.RefreshToken)); err != nil {
		if err := services.SendErrorResponse(w, "invalid refresh token"); err != nil {
			log.Printf("error when compare user refresh token: %s", err.Error())
		}
		return
	}

	tokens, hashedToken, err := services.GetNewTokens(user.LastIP)
	if err != nil {
		log.Printf("error when generating new tokens: %s", err.Error())
		return
	}
	user.HashedToken = hashedToken
	if user.LastIP != userIP {
		err := services.SendEmail(models.Email{From: os.Getenv("EMAIL_SENDER"), To: "test_user@gmail.ru"})
		if err != nil {
			log.Printf("error when send email: %s", err.Error())
			return
		}
	}
	user.LastIP = userIP

	if err := storage.UpdateUser(user, app.DB); err != nil {
		log.Printf("error when updating user: %s", err.Error())
		return
	}
	if err := services.SendResponse(w, models.Response{Answer: "OK", HTTPStatus: 200, Tokens: tokens}); err != nil {
		log.Printf("error when sending refresh-response: %s", err.Error())
		return
	}

}
