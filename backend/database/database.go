package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Kardainn/accountability/backend/config"
	_ "github.com/lib/pq"
)

func ConnectDB(ctx context.Context) *sql.DB {
	// connection string
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", config.FromContext(ctx).DatabaseHost, config.FromContext(ctx).DatabasePort, config.FromContext(ctx).DatabaseUser, config.FromContext(ctx).DatabasePassword, config.FromContext(ctx).DatabaseName)

	// open database
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		return nil
	}

	// check db
	err = db.Ping()
	if err != nil {
		return nil
	}

	fmt.Println("Connected!")

	return db
}
