package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Kardainn/accountability/backend/config"
	_ "github.com/lib/pq"
)

type Database struct {
	ctx      context.Context
	database *sql.DB
}

func ConnectDB(ctx context.Context) *Database {
	// connection string
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", config.FromContext(ctx).DatabaseHost, config.FromContext(ctx).DatabasePort, config.FromContext(ctx).DatabaseUser, config.FromContext(ctx).DatabasePassword, config.FromContext(ctx).DatabaseName)

	// open database
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		return nil
	}

	// close database
	defer db.Close()

	// check db
	err = db.Ping()
	if err != nil {
		return nil
	}

	fmt.Println("Connected!")

	database := &Database{
		ctx:      ctx,
		database: db,
	}

	return database
}
