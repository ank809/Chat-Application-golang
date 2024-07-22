package models

import "github.com/dgrijalva/jwt-go"

type Claims struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	jwt.StandardClaims
}
