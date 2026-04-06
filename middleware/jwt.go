package middleware

import (
	"api-jwt-dua/configs"
	"api-jwt-dua/utils"
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type ContextKey string

const UserKey ContextKey = "user" // exported supaya bisa diakses controllers

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.JSONResponse(w, http.StatusUnauthorized, "Unauthorized", "", "")
			//		http.Error(w, "Unauthorized - No Token", http.StatusUnauthorized)
			return
		}
		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(configs.AppConfig.JWT.Secret), nil
		})
		if err != nil || !token.Valid {
			utils.JSONResponse(w, http.StatusUnauthorized, "Token Expired", "", tokenString)
			//http.Error(w, "Unauthorized - Invalid Token", http.StatusUnauthorized)
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			utils.JSONResponse(w, http.StatusUnauthorized, "Unauthorized Invalid Claims", "", tokenString)
			//http.Error(w, "Unauthorized Invalid Claims", http.StatusUnauthorized)
			return
		}
		//simpan ke contex biiar bisa dipakai di handler
		ctx := context.WithValue(r.Context(), UserKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
