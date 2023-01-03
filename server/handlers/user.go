package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	dto "server/dto/result"
	userdto "server/dto/user"
	"server/pkg/bcrypt"
	"server/repositories"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/golang-jwt/jwt/v4"
)

// SETUP HANLDER STRUCT
type handlerUser struct {
	UserRepository repositories.UserRepository
}

// SETUP HANDLER FUNCTION
func HandlerUser(UserRepository repositories.UserRepository) *handlerUser {
	return &handlerUser{UserRepository}
}

// FUNCTION GET USER
func (h *handlerUser) GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// GET USER ID FROM JWT TOKEN
	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	userID := int(userInfo["id"].(float64))

	// RUN REPOSITORY GET USER
	user, err := h.UserRepository.GetUser(userID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	// CHECK IMAGE
	var Image string
	if user.Image != "" {
		Image = os.Getenv("PATH_FILE") + user.Image
	}

	// CREATE RESPONSE
	userResponse := userdto.UserResponse{
		ID:      user.ID,
		Name:    user.Name,
		Email:   user.Email,
		Image:   Image,
		Phone:   user.Phone,
		Address: user.Address,
	}

	// WRITE RESPONSE
	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: "success", Data: userResponse}
	json.NewEncoder(w).Encode(response)
}

// FUNCTION UPDATE USER
func (h *handlerUser) UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// GET USER ID FROM JWT TOKEN
	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	userID := int(userInfo["id"].(float64))

	//GET UPLOAD IMAGE
	dataContext := r.Context().Value("dataFile")
	filepath := dataContext.(string)

	// GET PARAMETER FROM ENV
	var ctx = context.Background()
	var CLOUD_NAME = os.Getenv("CLOUD_NAME")
	var API_KEY = os.Getenv("API_KEY")
	var API_SECRET = os.Getenv("API_SECRET")

	// SETUP CLOUDINARY CREDENTIALS
	cld, _ := cloudinary.NewFromParams(CLOUD_NAME, API_KEY, API_SECRET)

	// UPLOAD FILE TO CLOUDINARY
	respImage, _ := cld.Upload.Upload(ctx, filepath, uploader.UploadParams{Folder: "waysbeans"})

	var requestImage string
	if respImage == nil {
		requestImage = filepath
	} else {
		requestImage = respImage.SecureURL
	}

	// GET UPDATE REQUEST
	updateUser := userdto.UpdateUserRequest{
		Name:     r.FormValue("name"),
		Password: r.FormValue("password"),
		Image:    requestImage,
		Phone:    r.FormValue("phone"),
		Address:  r.FormValue("address"),
	}

	// RUN REPOSITORY GET USER
	user, err := h.UserRepository.GetUser(userID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	// CHECK UPDATE VALUE
	if updateUser.Name != "" {
		user.Name = updateUser.Name
	}

	if updateUser.Password != "" {
		password, _ := bcrypt.HashingPassword(updateUser.Password)
		user.Password = password
	}

	if filepath != "false" {
		user.Image = updateUser.Image
	}

	if updateUser.Address != "" {
		user.Address = updateUser.Address
	}

	if updateUser.Phone != "" {
		user.Phone = updateUser.Phone
	}

	// CHANGE VALUE UPDATE TIME
	user.UpdateAt = time.Now()

	//RUN REPOSITORY UPDATE USER
	data, err := h.UserRepository.UpdateUser(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	// CREATE RESPONSE
	updateResponse := userdto.UpdateUserResponse{
		Name:    data.Name,
		Image:   data.Image,
		Phone:   data.Phone,
		Address: data.Address,
	}

	// WRITE RESPONSE
	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: "success", Data: updateResponse}
	json.NewEncoder(w).Encode(response)
}
