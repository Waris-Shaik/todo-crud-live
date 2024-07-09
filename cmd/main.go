package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/Waris-Shaik/todo/cmd/api"
	"github.com/Waris-Shaik/todo/configs"
	"github.com/Waris-Shaik/todo/db"
	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func init() {
	// load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetOutput(os.Stdout)

	// get the port from .env file
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	port = fmt.Sprintf(":%v", port)

	db, err := db.MyNewSQLStorage(mysql.Config{
		User:                 configs.Envs.DBUser,
		Passwd:               configs.Envs.DBPassword,
		Addr:                 configs.Envs.DBAddress,
		DBName:               configs.Envs.DBName,
		Net:                  "tcp",
		AllowNativePasswords: true,
		ParseTime:            true,
	})

	if err != nil {
		log.Fatal("Could not connect to database successfully", err)
	}

	initStorage(db)

	// server-instance
	server := api.NewAPIServer(port, db)
	if err := server.Run(); err != nil {
		log.Fatal("Could not start server:", err)
	}

}

func initStorage(db *sql.DB) {
	err := db.Ping()
	if err != nil {
		log.Fatalf("Error while pinging database: %v", err)
	}
	log.Printf("Connected to database successfully on host:%v🔥🔥🔥 \n", configs.Envs.DBAddress)
}
