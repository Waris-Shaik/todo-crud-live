package main

import (
	"log"
	"os"
	"strconv"

	"github.com/Waris-Shaik/todo/configs"
	"github.com/Waris-Shaik/todo/db"
	mySqlCfg "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	// Database Configuration
	cfg := mySqlCfg.Config{
		User:                 configs.Envs.DBUser,
		Passwd:               configs.Envs.DBPassword,
		Addr:                 configs.Envs.DBAddress,
		DBName:               configs.Envs.DBName,
		Net:                  "tcp",
		AllowNativePasswords: true,
		ParseTime:            true,
	}

	// Database Connection
	db, err := db.MyNewSQLStorage(cfg)
	if err != nil {
		log.Fatalf("could not connect to database: %v", err)
	}

	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		log.Fatalf("error creating database driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://cmd/migrate/migrations",
		"mysql",
		driver,
	)
	if err != nil {
		log.Fatalf("error creating migration instance: %v", err)
	}

	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <command> [args]", os.Args[0])
	}

	cmd := os.Args[1]

	switch cmd {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("error applying migrations: %v", err)
		}
	case "down":
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("error reverting migrations: %v", err)
		}
	case "force":
		if len(os.Args) < 3 {
			log.Fatalf("force command requires a version argument")
		}
		version, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatalf("invalid version: %v", err)
		}
		if err := m.Force(version); err != nil {
			log.Fatalf("error forcing migration version to %d: %v", version, err)
		}
	default:
		log.Fatalf("unknown command: %s", cmd)
	}
}
