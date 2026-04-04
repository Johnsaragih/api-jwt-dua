package middleware

import (
	"api-jwt/configs"
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
			http.Error(w, "Unauthorized - No Token", http.StatusUnauthorized)
			return
		}
		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(configs.AppConfig.JWT.Secret), nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized - Invalid Token", http.StatusUnauthorized)
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Unauthorized Invalid Claims", http.StatusUnauthorized)
			return
		}
		//simpan ke contex biiar bisa dipakai di handler
		ctx := context.WithValue(r.Context(), UserKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
