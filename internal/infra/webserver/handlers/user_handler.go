package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/jailtonjunior94/go-products/internal/dto"
	"github.com/jailtonjunior94/go-products/internal/entity"
	"github.com/jailtonjunior94/go-products/internal/infra/database"

	"github.com/go-chi/jwtauth"
)

type Error struct {
	Message string `json:"message"`
}

type UserHandler struct {
	UserDB        database.UserInterface
	JWT           *jwtauth.JWTAuth
	JWTExperiesIn int
}

func NewUserHandler(db database.UserInterface, jwt *jwtauth.JWTAuth, jwtExperiesIn int) *UserHandler {
	return &UserHandler{
		UserDB:        db,
		JWT:           jwt,
		JWTExperiesIn: jwtExperiesIn,
	}
}

// GetJWT user godoc
// @Summary     Get a user JTW
// @Description Get a user JTW
// @Tags        Authorization
// @Accept      json
// @Produce     json
// @Param       request body     dto.GetJWTInput true "user request"
// @Success     200     {object} dto.TokenOutput
// @Failure     401     {object} Error
// @Failure     500     {object} Error
// @Router      /token [post]
func (h *UserHandler) GetJWT(w http.ResponseWriter, r *http.Request) {
	var user dto.GetJWTInput
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	u, err := h.UserDB.FindByEmail(user.Email)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if !u.ValidatePassword(user.Password) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	_, tokenString, _ := h.JWT.Encode(map[string]interface{}{
		"sub": u.ID.String(),
		"exp": time.Now().Add(time.Second * time.Duration(h.JWTExperiesIn)).Unix(),
	})

	accessToken := dto.TokenOutput{AccessToken: tokenString}

	w.Header().Set("Contenty-Type", "application/json")
	json.NewEncoder(w).Encode(accessToken)
	w.WriteHeader(http.StatusOK)
}

// Create user godoc
// @Summary     Create user
// @Description Create user
// @Tags        Users
// @Accept      json
// @Produce     json
// @Param       request body dto.CreateUserInput true "user request"
// @Success     201
// @Failure     500 {object} Error
// @Router      /users [post]
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user dto.CreateUserInput
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	u, err := entity.NewUser(user.Name, user.Email, user.Password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.UserDB.Create(u)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
