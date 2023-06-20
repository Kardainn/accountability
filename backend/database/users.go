package database

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type userCreation struct {
	Username string `json:"username"`
	ID       uint16 `json:"id"`
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	DOB      string `json:"dob"`
	Password string `json:"password"`
}

func CreateUser(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var err error = nil

	// connexion to DB
	ctxDatabase := ConnectDB(ctx)
	if ctxDatabase == nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("no database connected"))
		return
	}

	// cheking request
	bodyBytes, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	// checking content type
	ct := r.Header.Get("content-type")
	if ct != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(fmt.Sprintf("expected content-type to be 'application/json' but got '%s' instead", ct)))
		return
	}

	// creating uuid
	bigid := uuid.New().ID()

	// checking if uuid is already attribuated
	exists, err := ctxDatabase.CheckUserBigID(bigid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	if exists {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("cannot create user id '%c' is already used", bigid)))
		return
	}

	// reading request
	var userCreation userCreation
	err = json.Unmarshal(bodyBytes, &userCreation)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	if userCreation.Username == "" || userCreation.ID == 0 || userCreation.Password == "" {
		w.WriteHeader(http.StatusExpectationFailed)
		w.Write([]byte("mandatory field not present"))
	}

	_, err = ctxDatabase.database.ExecContext(ctx, `INSERT INTO users ("bigId", "username", "id", "name", "surname", "dob", "creation", "isActive", "lastLog", "password") values($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		bigid, userCreation.Username, userCreation.ID, userCreation.Name, userCreation.Surname, userCreation.DOB, time.Now(), true, time.Now(), userCreation.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User created succesfully"))
	return
}

func (db *Database) CheckUserBigID(bigid uint32) (bool, error) {
	var exists bool = true

	err := db.database.QueryRowContext(db.ctx, `SELECT(EXISTS(SELECT * FROM users WHERE "bigId" = $1))`, bigid).Scan(&exists)
	if err != nil {
		return true, fmt.Errorf("an error has occured while interogating db")
	}

	return exists, nil
}

func (db *Database) CheckUsernameID(username string, id uint16) (bool, error) {
	var exists bool = true

	err := db.database.QueryRowContext(db.ctx, `SELECT(EXISTS(SELECT * FROM users WHERE "username" = $1 AND "id" = $2))`, username, id).Scan(&exists)
	if err != nil {
		return true, fmt.Errorf("an error has occured while interogating db")
	}

	return exists, nil
}
