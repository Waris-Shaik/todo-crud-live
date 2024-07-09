package api

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/Waris-Shaik/todo/services/todo"
	"github.com/Waris-Shaik/todo/services/user"
	"github.com/Waris-Shaik/todo/utils"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file", err)
	}
}

type APIServer struct {
	addr string
	db   *sql.DB
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{addr: addr, db: db}
}

func (s *APIServer) Run() error {
	router := mux.NewRouter().StrictSlash(true)
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	// user-store
	userStore := user.NewStore(s.db)
	// user-handler
	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(subrouter)

	// todo-store
	todoStore := todo.NewStore(s.db)
	// todo-handler
	todoHandler := todo.NewHandler(todoStore, userStore)
	todoHandler.RegisterRoutes(subrouter)

	// CORS configuration
	frontendURL := os.Getenv("FRONTEND_URL")
	corsOptions := handlers.AllowedOrigins([]string{frontendURL})
	corsMethods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "PATCH", "DELETE"})
	corsHeaders := handlers.AllowedHeaders([]string{"Content-Type", "Authorization"})
	corsCredentials := handlers.AllowCredentials()

	// Apply CORS middleware
	handler := handlers.CORS(corsOptions, corsMethods, corsHeaders, corsCredentials)(router)

	log.Printf("Server is listening on PORT %v in %v mode⚡⚡⚡ \n", s.addr, utils.GetNodeENV("NODE_ENV"))
	return http.ListenAndServe(s.addr, handler)
}
