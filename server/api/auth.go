package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"

	"github.com/google/uuid"
	"github.com/lib/pq"
	database "github.com/sidsharma96/chatapp-backend/db/sqlc"
	"github.com/sidsharma96/chatapp-backend/utils"
)

type Auth struct {
	Server *Server
}

func (a Auth) router(server *Server) {
	a.Server = server

	authRouter := chi.NewRouter()
	server.ChiRouter.Mount("/auth", authRouter)
	authRouter.Post("/login", a.login)
	authRouter.Post("/logout", a.logout)
	authRouter.Post("/register", a.register)
}

type UserParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SuccessMessage struct {
	SuccessMsg string `json:"message"`
}

func (a *Auth) register(w http.ResponseWriter, r *http.Request) {
	userParams := UserParams{}
	if err := json.NewDecoder(r.Body).Decode(&userParams); err != nil {
		errorHandler(w, http.StatusBadRequest, fmt.Sprintf("Error decoding %v", err))
		return
	}

	if userParams.Username == "" || userParams.Password == "" {
		errorHandler(w, http.StatusBadRequest, fmt.Sprintln("Username and password cannot be empty!"))
		return
	}

	hashedPassword, err := utils.GenerateHashedPassword(userParams.Password)
	if err != nil {
		errorHandler(w, http.StatusInternalServerError, fmt.Sprint(err.Error()))
		return
	}

	args := database.CreateUserParams{
		ID:             uuid.New(),
		Username:       userParams.Username,
		HashedPassword: hashedPassword,
		CreatedAt:      time.Now().UTC(),
		UpdatedAt:      time.Now().UTC(),
	}

	user, err := a.Server.DbQueries.CreateUser(r.Context(), args)
	if err != nil {
		if pgError, ok := err.(*pq.Error); ok {
			if pgError.Code == "23505" {
				errorHandler(w, http.StatusBadRequest, "User already exists!")
				return
			}
		}
		errorHandler(w, http.StatusInternalServerError, fmt.Sprint(err.Error()))
		return
	}
	jsonHandler(w, http.StatusOK, UserResponse{}.toUserResponse(user))
}

func (a *Auth) login(w http.ResponseWriter, r *http.Request) {
	userParams := UserParams{}
	if err := json.NewDecoder(r.Body).Decode(&userParams); err != nil {
		errorHandler(w, http.StatusBadRequest, fmt.Sprintf("Error decoding %v", err))
		return
	}

	if userParams.Username == "" || userParams.Password == "" {
		errorHandler(w, http.StatusBadRequest, "Username and password cannot be empty!")
		return
	}

	user, err := a.Server.DbQueries.GetUserByUsername(r.Context(), userParams.Username)
	if err != nil {
		errorHandler(w, http.StatusBadRequest, "No username exists with the one provided")
		return
	}

	if err := utils.VerifyPassword(userParams.Password, user.HashedPassword); err != nil {
		errorHandler(w, http.StatusBadRequest, "Incorrect password!")
		return
	}

	godotenv.Load(".env")
	signingKey := os.Getenv("SIGNING_KEY")
	if signingKey == "" {
		log.Fatal("signing key is not found in the environment!")
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ID:        user.ID.String(),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
	})

	token, err := claims.SignedString([]byte(signingKey))
	if err != nil {
		errorHandler(w, http.StatusInternalServerError, err.Error())
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(time.Hour * 12),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	})

	success := SuccessMessage{
		SuccessMsg: "Login succesful!",
	}
	jsonHandler(w, http.StatusOK, success)
}

func (a *Auth) logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    "",
		Path:     "/",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	})

	success := SuccessMessage{
		SuccessMsg: "Logout succesful!",
	}
	jsonHandler(w, http.StatusOK, success)
}
