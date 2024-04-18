package api

import (
	"log"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	db "github.com/sidsharma96/chatapp-backend/db/sqlc"
)

type authedHandler func(http.ResponseWriter, *http.Request, db.User)

func (server *Server) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		godotenv.Load(".env")
		signingKey := os.Getenv("SIGNING_KEY")
		if signingKey == "" {
			log.Fatal("signing key is not found in the environment!")
		}

		cookie, err := r.Cookie("jwt")
		if err != nil {
			errorHandler(w, http.StatusUnauthorized, "invalid authentication token")
			return
		}
		token, err := jwt.ParseWithClaims(cookie.Value, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
			return []byte(signingKey), nil
		})

		if err != nil {
			errorHandler(w, http.StatusUnauthorized, err.Error())
			return
		}

		claims := token.Claims.(*jwt.RegisteredClaims)
		user, err := server.DbQueries.GetUserByID(r.Context(), uuid.MustParse(claims.ID))
		if err != nil {
			errorHandler(w, http.StatusInternalServerError, err.Error())
			return
		}
		handler(w, r, user)
	}
}
