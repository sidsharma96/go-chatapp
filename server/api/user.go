package api

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	database "github.com/sidsharma96/chatapp-backend/db/sqlc"
)

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u UserResponse) toUserResponse(dbUser database.User) UserResponse {
	return UserResponse{
		ID:        dbUser.ID,
		Username:  dbUser.Username,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
	}
}

type User struct {
	Server *Server
}

func (u User) router(server *Server) {
	u.Server = server

	userRouter := chi.NewRouter()
	server.ChiRouter.Mount("/users", userRouter)
	userRouter.Get("/me", u.getActiveUser)
}

func (u *User) getActiveUser(w http.ResponseWriter, r *http.Request) {
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
	user, err := u.Server.DbQueries.GetUserByID(r.Context(), uuid.MustParse(claims.ID))
	if err != nil {
		errorHandler(w, http.StatusInternalServerError, err.Error())
	}
	jsonHandler(w, http.StatusOK, UserResponse{}.toUserResponse(user))
}
