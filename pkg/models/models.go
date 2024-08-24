package models

type User struct {
	ID          string
	HashedToken string
	LastIP      string
}

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
