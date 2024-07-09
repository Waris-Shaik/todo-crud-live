package auth

import "golang.org/x/crypto/bcrypt"

func HashPassword(password *string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(*password), 10)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func MatchPassword(userPassword, logInPassword *string) bool {
	result := bcrypt.CompareHashAndPassword([]byte(*userPassword), []byte(*logInPassword))
	return result == nil
}
