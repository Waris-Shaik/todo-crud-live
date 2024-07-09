package todo

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Waris-Shaik/todo/services/auth"
	"github.com/Waris-Shaik/todo/types"
	"github.com/Waris-Shaik/todo/utils"
	"github.com/gorilla/mux"
)

type Handler struct {
	store     types.TodoStore
	userstore types.UserStore
}

func NewHandler(store types.TodoStore, userstore types.UserStore) *Handler {
	return &Handler{store: store, userstore: userstore}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {

	router.HandleFunc("/todos", auth.WithJWTAuth(h.handleGetTodos, h.userstore)).Methods(http.MethodGet)
	router.HandleFunc("/todos/new", auth.WithJWTAuth(h.handleCreateTodo, h.userstore)).Methods(http.MethodPost)
	router.HandleFunc("/todos/{id}", auth.WithJWTAuth(h.handleGetTodo, h.userstore)).Methods(http.MethodGet)
	router.HandleFunc("/todos/update/{id}", auth.WithJWTAuth(h.handleUpdateTodo, h.userstore)).Methods(http.MethodPatch)
	router.Handle("/todos/delete/{id}", auth.WithJWTAuth(h.handleDeleteTodo, h.userstore)).Methods(http.MethodDelete)

}

func (h *Handler) handleCreateTodo(w http.ResponseWriter, r *http.Request) {
	// authenticate user
	userID := auth.GetUserIDFromContext(r.Context())

	// get the JSON payloaf from req.body and parse it
	var payload types.TodoPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		log.Println("Error parsing PAYLOAD:", err)
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate the payload
	if err := utils.ValidateTodoPayload(&payload); err != nil {
		log.Println("Error validating payload", err)
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// create the todo
	err := h.store.CreateTodo(types.Todo{
		Title:       payload.Title,
		Description: payload.Description,
		UserID:      userID,
	})

	if err != nil {
		log.Println("Error while creating todo", err)
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// return the response
	response := struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}{
		Success: true,
		Message: "todo created successfully",
	}

	utils.WriteJSON(w, http.StatusCreated, response)
}

func (h *Handler) handleGetTodos(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserIDFromContext(r.Context())
	todos, err := h.store.GetTodos(userID)
	if err != nil {
		log.Println("Error while retreiving todos", err)
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// retrun thr response
	response := struct {
		Success bool          `json:"success"`
		Todos   []*types.Todo `json:"todos"`
	}{
		Success: true,
		Todos:   todos,
	}
	utils.WriteJSON(w, http.StatusOK, response)
}

func (h *Handler) handleGetTodo(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserIDFromContext(r.Context())
	if userID == -1 {
		log.Println("Unauthorized access: Invalid user ID")
		utils.WriteError(w, http.StatusForbidden, fmt.Errorf("unauthorized access"))
		return
	}

	// extract userID from req.params
	vars := mux.Vars(r)
	idStr := vars["id"]

	if idStr == "" {
		log.Println("TodoID is required")
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("todoID is required"))
		return
	}

	todoID, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println("Failed to convert todoID:", err)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("failed to convert str to int"))
		return
	}

	// retreive todo from store
	todo, err := h.store.GetTodoByID(todoID)
	if err != nil {
		log.Println("invalid id todo not found", err)
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("todo not found"))
		return
	}

	// Ensure that the task belongs to the current user (if needed)
	if userID != todo.UserID {
		log.Println("Unauthorized access: Task does not belong to the user")
		utils.WriteError(w, http.StatusForbidden, fmt.Errorf("unauthorized access"))
		return
	}
	// return the success response
	response := struct {
		Success bool       `json:"success"`
		Task    types.Todo `json:"todo"`
	}{
		Success: true,
		Task:    *todo,
	}

	utils.WriteJSON(w, http.StatusOK, response)

}

func (h *Handler) handleUpdateTodo(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserIDFromContext(r.Context())
	if userID == -1 {
		log.Println("Unauthorized access: Invalid user ID")
		utils.WriteError(w, http.StatusForbidden, fmt.Errorf("unauthorized access"))
		return
	}

	// extract userID from req.params
	vars := mux.Vars(r)
	idStr := vars["id"]

	if idStr == "" {
		log.Println("TodoID is required")
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("todoID is required"))
		return
	}

	todoID, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println("Failed to convert todoID:", err)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("failed to convert str to int"))
		return
	}

	// retreive todo from store
	todo, err := h.store.GetTodoByID(todoID)
	if err != nil {
		log.Println("invalid id todo not found", err)
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("todo not found"))
		return
	}

	// Ensure that the task belongs to the current user (if needed)
	if userID != todo.UserID {
		log.Println("Unauthorized access: Task does not belong to the user")
		utils.WriteError(w, http.StatusForbidden, fmt.Errorf("unauthorized access"))
		return
	}
	// Update the todo status
	if err := h.store.UpdateTodo(todo.ID); err != nil {
		log.Println("Error updating task:", err)
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// Return success response
	response := struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}{
		Success: true,
		Message: "todo updated successfully",
	}

	utils.WriteJSON(w, http.StatusOK, response)
}

func (h *Handler) handleDeleteTodo(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserIDFromContext(r.Context())
	if userID == -1 {
		log.Println("Unauthorized access: Invalid user ID")
		utils.WriteError(w, http.StatusForbidden, fmt.Errorf("unauthorized access"))
		return
	}

	// extract userID from req.params
	vars := mux.Vars(r)
	idStr := vars["id"]

	if idStr == "" {
		log.Println("TodoID is required")
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("todoID is required"))
		return
	}

	todoID, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println("Failed to convert todoID:", err)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("failed to convert str to int"))
		return
	}

	// retreive todo from store
	todo, err := h.store.GetTodoByID(todoID)
	if err != nil {
		log.Println("invalid id todo not found", err)
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("todo not found"))
		return
	}

	// Ensure that the task belongs to the current user (if needed)
	if userID != todo.UserID {
		log.Println("Unauthorized access: Task does not belong to the user")
		utils.WriteError(w, http.StatusForbidden, fmt.Errorf("unauthorized access"))
		return
	}

	// delete the todo
	err = h.store.DeleteTodo(todo.ID)
	if err != nil {
		log.Println("Error while deleting todo:", err)
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// Return success response
	response := struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}{
		Success: true,
		Message: "todo deleted successfully",
	}

	utils.WriteJSON(w, http.StatusOK, response)

}
