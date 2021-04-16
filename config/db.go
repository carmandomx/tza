package config

import (
	"context"
	"fmt"
	"os"

	"github.com/go-pg/pg/v10"
	"github.com/joho/godotenv"
)

var PgUrl string
var ZoomToken string

func SetupDb() (*pg.DB, error) {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
	}
	PgUrl = os.Getenv("DB_URL")
	ZoomToken = os.Getenv("ZOOM_TOKEN")
	opt, err := pg.ParseURL(PgUrl)

	if err != nil {
		panic(err)
	}

	db := pg.Connect(opt)
	ctx := context.Background()

	if err := db.Ping(ctx); err != nil {
		panic(err)
	}

	return db, nil
}
