package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	authdto "server/dto/auth"
	dto "server/dto/result"
	"server/models"
	"server/pkg/bcrypt"
	jwtToken "server/pkg/jwt"
	"server/repositories"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
)

// SETUP HANLDER STRUCT
type handlerAuth struct {
	AuthRepository repositories.AuthRepository
}

// SETUP HANDLER FUNCTION
func HandlerAuth(AuthRepository repositories.AuthRepository) *handlerAuth {
	return &handlerAuth{AuthRepository}
}

// FUNCTION REGISTER USER
func (h *handlerAuth) Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// GET REQUEST AND DECODING JSON
	request := new(authdto.RegisterRequest)
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	// VALIDATE REQUEST INPUT
	validation := validator.New()
	err := validation.Struct(request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	// ENCRYPT PASSWORD
	password, err := bcrypt.HashingPassword(request.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
	}

	// SETUP FOR QUERY REGISTER
	user := models.User{
		Name:     request.Name,
		Email:    request.Email,
		Password: password,
		Role:     "user",
		CreateAt: time.Now(),
		UpdateAt: time.Now(),
	}

	// RUN REPOSITORY REGISTER
	data, err := h.AuthRepository.Register(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: "email is registered"}
		json.NewEncoder(w).Encode(response)
		return
	}

	// CREATE RESPONSE
	registerResponse := authdto.RegisterResponse{
		ID:    data.ID,
		Name:  data.Name,
		Email: data.Email,
	}

	// WRITE RESPONSE
	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: "success", Data: registerResponse}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerAuth) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// GET REQUEST AND DECODING JSON
	request := new(authdto.LoginRequest)
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	// RUN REPOSITORY LOGIN
	user, err := h.AuthRepository.Login(request.Email)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: "wrong email or password"}
		json.NewEncoder(w).Encode(response)
		return
	}

	// CHECK PASSWORD REQUEST
	isValid := bcrypt.CheckPasswordHash(request.Password, user.Password)
	if !isValid {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: "wrong email or password"}
		json.NewEncoder(w).Encode(response)
		return
	}

	// GENERATE TOKEN
	claims := jwt.MapClaims{}
	claims["id"] = user.ID
	claims["role"] = user.Role
	claims["exp"] = time.Now().Add(time.Hour * 2).Unix() // TIME LIMIT 2 HOURS

	token, errGenerateToken := jwtToken.GenerateToken(&claims)
	if errGenerateToken != nil {
		log.Println(errGenerateToken)
		fmt.Println("Unauthorize")
		return
	}

	// CREATE RESPONSE
	loginResponse := authdto.LoginResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
		Token: token,
	}

	// WRITE RESPONSE
	w.Header().Set("Content-Type", "application/json")
	response := dto.SuccessResult{Code: "success", Data: loginResponse}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerAuth) CheckAuth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// GET USERID FROM TOKEN
	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	userID := int(userInfo["id"].(float64))

	// RUN REPOSITORY GET USER AUTH
	user, err := h.AuthRepository.GetUserAuth(userID)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		response := dto.ErrorResult{Code: http.StatusUnauthorized, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	// CREATE RESPONSE
	checkAuth := authdto.CheckAuthResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}

	// WRITE RESPONSE
	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: "success", Data: checkAuth}
	json.NewEncoder(w).Encode(response)
}
