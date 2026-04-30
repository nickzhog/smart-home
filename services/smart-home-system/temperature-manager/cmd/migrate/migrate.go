package main

import (
	"embed"
	"log"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

func main() {
	time.Sleep(time.Second * 2)
	connString, _ := os.LookupEnv("DATABASE_URI")
	if connString == "" {
		panic("conn string empty")
	}
	err := Migrate(connString)
	if err != nil {
		log.Printf("cant migrate err: %v\n", err)
		os.Exit(255)
	}
	log.Printf("ok\n")
}

//go:embed migrations/*
var migrations embed.FS

func Migrate(connString string) error {
	src, err := iofs.New(migrations, "migrations")
	if err != nil {
		return err
	}
	m, err := migrate.NewWithSourceInstance("iofs", src, connString)
	if err != nil {
		return err
	}
	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}
