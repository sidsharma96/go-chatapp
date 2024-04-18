package api

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
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
	userRouter.Get("/me", u.Server.middlewareAuth(u.getAuthenticatedActiveUser))
}

func (u *User) getAuthenticatedActiveUser(w http.ResponseWriter, r *http.Request, user database.User) {
	jsonHandler(w, http.StatusOK, UserResponse{}.toUserResponse(user))
}
