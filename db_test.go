package broker

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"path"
	"testing"
)

const (
	TestUser = "tochti"
	TestPass = "123"
	TestHost = "127.0.0.1"
	TestPort = 3306
	TestDB   = "testing"
)

func handleError(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

func equalTokenRequest(a, b TokenRequest) bool {
	if a.ID != b.ID ||
		a.Token != b.Token {
		return false
	}

	return true
}

func existsID(l []string, id string) bool {
	for _, i := range l {
		if i == id {
			return true
		}
	}

	return false
}

func readSQL(file string) (string, error) {
	p := path.Join("sql", file)
	content, err := ioutil.ReadFile(p)

	return string(content), err
}

func newTestEnv(t *testing.T) *sql.DB {
	url := fmt.Sprintf("%v:%v@tcp(%v:%v)/", TestUser, TestPass, TestHost, TestPort)
	db, err := sql.Open("mysql", url)
	handleError(t, err)

	_, err = db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %v", TestDB))
	handleError(t, err)

	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %v", TestDB))
	handleError(t, err)

	_, err = db.Exec(fmt.Sprintf("USE %v", TestDB))
	handleError(t, err)

	query, err := readSQL("new-env.sql")
	handleError(t, err)

	err = db.Ping()
	handleError(t, err)

	_, err = db.Exec(query)
	handleError(t, err)

	return db
}

func removeTestEnv(t *testing.T, db *sql.DB) {
	_, err := db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %v", TestDB))
	handleError(t, err)
	db.Close()
}

func Test_NewID_OK(t *testing.T) {
	l := []string{}
	for x := 0; x < 10; x++ {
		token, err := NewID()
		handleError(t, err)
		if existsID(l, token) {
			t.Fatal("Double token")
		}
		l = append(l, token)
	}
}

func Test_CRUDTokenRequest_OK(t *testing.T) {
	db := newTestEnv(t)
	defer removeTestEnv(t, db)

	id, err := InitTokenRequest(db)
	handleError(t, err)
	if len(id) <= 0 {
		t.Fatal("Expect token to be longer then 0")
	}

	req, err := ReadTokenRequest(db, id)
	handleError(t, err)
	expect := TokenRequest{
		ID: id,
	}
	if !equalTokenRequest(expect, req) {
		t.Fatal("Expect", expect, "was", req)
	}

	token := "1234"
	err = AppendToken(db, id, token)
	handleError(t, err)
	req, err = ReadTokenRequest(db, id)
	handleError(t, err)
	expect = TokenRequest{
		ID:    id,
		Token: token,
	}
	if !equalTokenRequest(expect, req) {
		t.Fatal("Expect", expect, "was", req)
	}

	err = RemoveTokenRequest(db, id)
	handleError(t, err)
	_, err = ReadTokenRequest(db, id)
	if err != sql.ErrNoRows {
		t.Fatal("Expect", sql.ErrNoRows, "was", err)
	}

}
