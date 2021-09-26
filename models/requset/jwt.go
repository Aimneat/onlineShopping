package requset

import "github.com/dgrijalva/jwt-go"

type CustomClaims struct {
	Telephone string
	UserName  string
	jwt.StandardClaims
}
