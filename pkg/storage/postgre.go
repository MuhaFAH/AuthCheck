package storage

import (
	"fmt"
	"github.com/MuhaFAH/AuthCheck/pkg/models"
	"github.com/jmoiron/sqlx"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
	SSLMode  string
}

func NewDB(cfg Config) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", cfg.Host, cfg.Port, cfg.Username, cfg.Database, cfg.Password, cfg.SSLMode))
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func AddUser(u models.User, db *sqlx.DB) error {
	query := "INSERT INTO users (user_id, last_ip, token_hash) VALUES (:user_id, :last_ip, :token_hash)"
	_, err := db.NamedExec(query, u)
	if err != nil {
		return err
	}
	return nil
}

func GetUser(u *models.User, db *sqlx.DB) error {
	query := "SELECT user_id, last_ip, token_hash FROM users WHERE user_id=$1"
	if err := db.Get(u, query, u.GUID); err != nil {
		return err
	}
	return nil
}

func UpdateUser(u models.User, db *sqlx.DB) error {
	query := "UPDATE users SET last_ip = :last_ip, token_hash = :token_hash WHERE user_id = :user_id"
	_, err := db.NamedExec(query, u)
	if err != nil {
		return err
	}
	return nil
}
