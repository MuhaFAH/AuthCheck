package models

type User struct {
	GUID        string `db:"user_id"`
	HashedToken string `db:"token_hash"`
	LastIP      string `db:"last_ip"`
}

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type Response struct {
	Answer     string `json:"answer"`
	Reason     string `json:"reason"`
	HTTPStatus int    `json:"status"`
	Tokens     `json:"tokens"`
}

type Email struct {
	From string
	To   string
}
