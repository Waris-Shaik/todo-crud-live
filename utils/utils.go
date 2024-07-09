package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/Waris-Shaik/todo/types"
	"github.com/go-playground/validator/v10"
)

func ParseJSON(r *http.Request, payload any) error {
	if r.Body == nil {
		return fmt.Errorf("missing request body")
	}
	return json.NewDecoder(r.Body).Decode(payload)
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)
}

func WriteError(w http.ResponseWriter, status int, err error) {
	response := struct {
		Success bool   `json:"success"`
		Error   string `json:"error"`
	}{
		Success: false,
		Error:   err.Error(),
	}

	WriteJSON(w, status, response)
}

func ValidateRegisterUserPayload(payload *types.RegisterUserPayload) error {
	err := validator.New().Struct(*payload)
	if err != nil {
		errorMessages := make([]string, 0)

		for _, validationError := range err.(validator.ValidationErrors) {
			errorMessages = append(errorMessages, validationError.Field())
		}

		if len(errorMessages) == 0 {
			fmt.Println("this field is required:", errorMessages)
		} else if len(errorMessages) > 1 {
			fmt.Println("these fields are required:", errorMessages)
		} else {
			fmt.Println("validation error", errorMessages)
		}
		return fmt.Errorf("%v is required", errorMessages[0])
	}
	return nil
}

func ValidateLoginUserPayload(payload *types.LoginUserPayload) error {
	err := validator.New().Struct(*payload)
	if err != nil {
		errorMessages := make([]string, 0)

		for _, validationError := range err.(validator.ValidationErrors) {
			errorMessages = append(errorMessages, validationError.Field())
		}

		if len(errorMessages) == 0 {
			fmt.Println("this field is required:", errorMessages)
		} else if len(errorMessages) > 1 {
			fmt.Println("these fields are required:", errorMessages)
		} else {
			fmt.Println("validation error", errorMessages)
		}
		return fmt.Errorf("%v is required", errorMessages[0])
	}
	return nil
}

func ValidateTodoPayload(payload *types.TodoPayload) error {
	err := validator.New().Struct(*payload)
	if err != nil {
		errorMessages := make([]string, 0)

		for _, validationError := range err.(validator.ValidationErrors) {
			errorMessages = append(errorMessages, validationError.Field())
		}

		if len(errorMessages) == 0 {
			fmt.Println("this field is required:", errorMessages)
		} else if len(errorMessages) > 1 {
			fmt.Println("these fields are required:", errorMessages)
		} else {
			fmt.Println("validation error", errorMessages)
		}
		return fmt.Errorf("%v is required", errorMessages[0])
	}
	return nil
}

func MatchPasswordCriteria(password *string) error {
	const (
		minPasswordLen = 6
		maxPasswordLen = 12
	)

	if len(*password) < minPasswordLen {
		return fmt.Errorf("password atleast should contain %d charcters", minPasswordLen)
	}
	if len(*password) > maxPasswordLen {
		return fmt.Errorf("password must not be greater than %d characters", maxPasswordLen)
	}
	return nil
}

func GetNodeENV(key string) string {
	return os.Getenv(key)
}
