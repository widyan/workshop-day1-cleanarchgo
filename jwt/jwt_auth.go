package jwt

import (
	"encoding/json"
	"github.com/gorilla/context"
	"github.com/sangianpatrick/devoria-article-service/entity"
	"github.com/sangianpatrick/devoria-article-service/response"
	"net/http"
	"strings"
)

type JwtMiddleware interface {
	VerifyToken(next http.HandlerFunc) http.HandlerFunc
}

type JwtToken struct {
	jsonWebToken JSONWebToken
}

func (j *JwtToken) VerifyToken(next http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var resp response.Response
		tokenString := request.Header.Get("Authorization")
		strArr := strings.Split(tokenString, " ")
		if len(strArr) == 2 {
			tokenString = strArr[1]
		} else {
			resp = response.Error(response.StatusUnauthorized, nil, nil)
			resp.JSON(writer)
			return
		}

		var claims entity.AccountStandardJWTClaims
		err := j.jsonWebToken.Parse(request.Context(),tokenString,&claims)
		if err != nil{
			resp = response.Error(response.StatusUnauthorized, nil, err)
			resp.JSON(writer)
			return
		}

		byt, _ := json.Marshal(claims)
		context.Set(request, "bind", byt)
		next.ServeHTTP(writer, request)
	}
}

func NewJwtToken(jsonWebToken JSONWebToken) JwtMiddleware {
	return &JwtToken{jsonWebToken:jsonWebToken}
}
