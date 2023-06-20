package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Kardainn/accountability/backend/config"
	_ "github.com/lib/pq"
)

func ConnectDB(ctx context.Context) {
	// connection string
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", config.FromContext(ctx).DatabaseHost, config.FromContext(ctx).DatabasePort, config.FromContext(ctx).DatabaseUser, config.FromContext(ctx).DatabasePassword, config.FromContext(ctx).DatabaseName)

	// open database
	db, err := sql.Open("postgres", psqlconn)
	CheckError(err)

	// close database
	defer db.Close()

	// check db
	err = db.Ping()
	CheckError(err)

	fmt.Println("Connected!")
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}
