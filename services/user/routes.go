package user

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Waris-Shaik/todo/services/auth"
	"github.com/Waris-Shaik/todo/types"
	"github.com/Waris-Shaik/todo/utils"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}
}

type Handler struct {
	store types.UserStore
}

func NewHandler(store types.UserStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/", h.handleRoot).Methods(http.MethodGet)
	router.HandleFunc("/register", h.handleRegister).Methods(http.MethodPost)
	router.HandleFunc("/login", h.handleLogin).Methods(http.MethodPost)
	router.HandleFunc("/logout", h.handleLogout).Methods(http.MethodPost)
	router.HandleFunc("/me", auth.WithJWTAuth(h.handleMyProfile, h.store)).Methods(http.MethodGet)

}

func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	// get the JSON payload from req.body and parse it
	var payload types.RegisterUserPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		log.Println("Error in parsing PAYLOAD:", err)
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate the payload
	if err := utils.ValidateRegisterUserPayload(&payload); err != nil {
		log.Println("Error while validating payload:", err)
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// match the password criteria
	if err := utils.MatchPasswordCriteria(&payload.Password); err != nil {
		log.Println("Error in matching password criteria", err)
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// check the user if exsits in db
	_, err := h.store.GetUserByEmail(payload.Email)
	if err == nil {
		log.Println("User with these credentials already exists")
		utils.WriteError(w, http.StatusConflict, fmt.Errorf("user already exists, please login"))
		return
	}

	// hash the password
	hashedPassword, err := auth.HashPassword(&payload.Password)
	if err != nil {
		log.Println("Error hashing password:", err)
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// create a new user
	userID, err := h.store.CreateUser(types.User{
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		UserName:  payload.UserName,
		Email:     payload.Email,
		Password:  hashedPassword,
	})

	if err != nil {
		log.Println("Error while creating a user:", err)
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// get the cookie
	token, err := auth.CreateJWT(userID)
	if err != nil {
		log.Println("Error in generating JWT token:", err)
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// set the JWT as cookie
	cookie := &http.Cookie{
		Name:     "token",
		Value:    token,
		HttpOnly: true,
		Path:     "/",
		Expires:  time.Now().Add(15 * time.Minute),
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
	}

	if utils.GetNodeENV("NODE_ENV") == "Development" {
		cookie.SameSite = http.SameSiteLaxMode
		cookie.Secure = false
	}

	http.SetCookie(w, cookie)

	// return the response
	reponse := struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}{
		Success: true,
		Message: "user created successfully",
	}

	utils.WriteJSON(w, http.StatusCreated, reponse)

}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	// get the JSON payload from req.body and parse it
	var payload types.LoginUserPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		log.Println("Error in parsing PAYLOAD:", err)
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate the payload
	if err := utils.ValidateLoginUserPayload(&payload); err != nil {
		log.Println("Error occured in validated payload:", err)
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// check the user if exsits in db
	user, err := h.store.GetUserByEmail(payload.Text)
	if err != nil {
		log.Println("Error user not found, please register", err)
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("user not found, please register"))
		return
	}

	// check the password matches or does not
	if !auth.MatchPassword(&user.Password, &payload.Password) {
		log.Println("Error password does not match:", !auth.MatchPassword(&user.Password, &payload.Password))
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid passsword"))
		return
	}

	// set the JWT token
	jwtToken, err := auth.CreateJWT(user.ID)
	if err != nil {
		log.Println("Error in generating JWT token:", err)
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	// set the JWT as cookie
	cookie := &http.Cookie{
		Name:     "token",
		Value:    jwtToken,
		HttpOnly: true,
		Path:     "/",
		Expires:  time.Now().Add(15 * time.Minute),
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
	}

	if utils.GetNodeENV("NODE_ENV") == "Development" {
		cookie.SameSite = http.SameSiteLaxMode
		cookie.Secure = false
	}

	http.SetCookie(w, cookie)

	// return the response
	response := struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}{
		Success: true,
		Message: fmt.Sprintf("welcome back %v", user.UserName),
	}

	utils.WriteJSON(w, http.StatusOK, response)

}

func (h *Handler) handleRoot(w http.ResponseWriter, r *http.Request) {
	response := struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}{
		Success: true,
		Message: "welcome to go-lang-server",
	}
	// return the response
	utils.WriteJSON(w, http.StatusOK, response)
}

func (h *Handler) handleLogout(w http.ResponseWriter, r *http.Request) {
	// clear the JWT token cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		HttpOnly: true,
		Path:     "/",
		Expires:  time.Unix(0, 0),
	})
	response := struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}{
		Success: true,
		Message: "successfully logged out",
	}

	utils.WriteJSON(w, http.StatusOK, response)
}

func (h *Handler) handleMyProfile(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserIDFromContext(r.Context())
	if userID == -1 {
		log.Println("Something went wrong", userID)
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("something went wrong"))
		return
	}

	user, err := h.store.GetUserByID(userID)
	if err != nil {
		log.Println("Error while retreiving user from db", err)
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// return the response
	response := struct {
		Success bool        `json:"success"`
		User    *types.User `json:"user"`
	}{
		Success: true,
		User:    user,
	}

	utils.WriteJSON(w, http.StatusOK, response)
}
