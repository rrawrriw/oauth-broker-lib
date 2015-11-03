package broker

import (
	"crypto/rand"
	"crypto/sha1"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var (
	TokenRequestTable = "TokenRequest"
)

type (
	TokenRequest struct {
		ID    string
		Token string
	}
)

func NewID() (string, error) {
	blob := make([]byte, 250)
	_, err := rand.Read(blob)
	if err != nil {
		return "", err
	}

	hash := fmt.Sprintf("%x", sha1.Sum(blob))
	return hash, nil
}

func InitTokenRequest(db *sql.DB) (string, error) {
	newID, err := NewID()
	if err != nil {
		return "", err
	}

	err = db.Ping()
	if err != nil {
		return "", err
	}

	q := fmt.Sprintf("INSERT INTO %v VALUES (?,?)", TokenRequestTable)
	_, err = db.Exec(q, newID, "")
	if err != nil {
		return "", err
	}

	return newID, nil

}

func ReadTokenRequest(db *sql.DB, id string) (TokenRequest, error) {
	err := db.Ping()
	if err != nil {
		return TokenRequest{}, err
	}

	var token string
	q := fmt.Sprintf("SELECT Token FROM %v WHERE ID = ?", TokenRequestTable)
	err = db.QueryRow(q, id).Scan(&token)
	if err != nil {
		return TokenRequest{}, err
	}

	tokenRequest := TokenRequest{
		ID:    id,
		Token: token,
	}

	return tokenRequest, nil
}

func AppendToken(db *sql.DB, id, token string) error {
	err := db.Ping()
	if err != nil {
		return err
	}

	q := fmt.Sprintf("UPDATE %v SET Token = ? WHERE ID = ?", TokenRequestTable)
	_, err = db.Exec(q, token, id)
	if err != nil {
		return err
	}

	return nil
}

func RemoveTokenRequest(db *sql.DB, id string) error {
	err := db.Ping()
	if err != nil {
		return err
	}

	q := fmt.Sprintf("DELETE FROM %v WHERE ID = ?", TokenRequestTable)
	_, err = db.Exec(q, id)
	if err != nil {
		return err
	}

	return nil
}
