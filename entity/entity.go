package entity

import "github.com/dgrijalva/jwt-go"

// CustomerStandardJWTClaims is a model.
type AccountStandardJWTClaims struct {
	jwt.StandardClaims
	Email string `json:"email"`
}
