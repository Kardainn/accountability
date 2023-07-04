package database

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
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

type userDB struct {
	Username string    `json:"username"`
	Id       uint16    `json:"id"`
	Name     string    `json:"name"`
	Surname  string    `json:"surname"`
	LastLog  time.Time `json:"lastLog"`
	Dob      string    `json:"dob"`
	Creation string    `json:"creation"`
	IsActive bool      `json:"isActive"`
	Password string    `json:"password"`
	BigId    string    `json:"bigId"`
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
	exists, err := CheckUserBigID(ctx, bigid)
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

	if userCreation.Username == "" || userCreation.Password == "" {
		w.WriteHeader(http.StatusExpectationFailed)
		w.Write([]byte("mandatory field not present"))
		return
	}

	if userCreation.DOB == "" {
		userCreation.DOB = "01-01-1970"
	}

	_, err = ctxDatabase.ExecContext(ctx, `INSERT INTO users ("bigId", "username", "id", "name", "surname", "dob", "creation", "isActive", "lastLog", "password") values($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);`,
		bigid, strings.ToLower(userCreation.Username), userCreation.ID, userCreation.Name, userCreation.Surname, userCreation.DOB, time.Now(), true, time.Now(), userCreation.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User created succesfully"))
}

func CheckUserBigID(ctx context.Context, bigid uint32) (bool, error) {
	var exists bool = true

	// connexion to DB
	ctxDatabase := ConnectDB(ctx)
	if ctxDatabase == nil {
		return true, fmt.Errorf("no database connected")
	}

	err := ctxDatabase.QueryRowContext(ctx, `SELECT(EXISTS(SELECT * FROM users WHERE "bigId" = $1)) as "exists";`, bigid).Scan(&exists)
	if err != nil {
		return true, err
	}

	return exists, nil
}

func CheckUsernameID(ctx context.Context, username string, id uint16) (bool, error) {
	var exists bool = true

	// connexion to DB
	ctxDatabase := ConnectDB(ctx)
	if ctxDatabase == nil {
		return true, fmt.Errorf("no database connected")
	}

	err := ctxDatabase.QueryRowContext(ctx, `SELECT(EXISTS(SELECT * FROM users WHERE "username" = $1 AND "id" = $2)) as "exists";`, username, id).Scan(&exists)
	if err != nil {
		return true, fmt.Errorf("an error has occured while interogating db")
	}

	return exists, nil
}

func GetUser(ctx context.Context, w http.ResponseWriter, r *http.Request, idResquested uint32) {
	// connexion to DB
	ctxDatabase := ConnectDB(ctx)
	if ctxDatabase == nil {
		return
	}

	fmt.Println(idResquested)

	var userDB userDB

	row := ctxDatabase.QueryRowContext(ctx, `SELECT * FROM users WHERE "bigId" = $1`, idResquested)
	err := row.Scan(&userDB.Username,
		&userDB.Id,
		&userDB.Name,
		&userDB.Surname,
		&userDB.Dob,
		&userDB.Creation,
		&userDB.IsActive,
		&userDB.LastLog,
		&userDB.Password,
		&userDB.BigId)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(userDB)

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userDB)
}

func PatchUser(ctx context.Context, w http.ResponseWriter, r *http.Request, id uint32) {
	// todo implement
}

func DeleteUser(ctx context.Context, w http.ResponseWriter, r *http.Request, id uint32) {
	// todo implement
}
