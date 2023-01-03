package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	productdto "server/dto/product"
	dto "server/dto/result"
	"server/models"
	"server/repositories"
	"strconv"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
)

// SETUP HANLDER STRUCT
type handlerProduct struct {
	ProductRepository repositories.ProductRepository
}

// SETUP HANDLER FUNCTION
func HandlerProduct(ProductRepository repositories.ProductRepository) *handlerProduct {
	return &handlerProduct{ProductRepository}
}

// FUNCTION FIND PRODUCTS
func (h *handlerProduct) FindProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// RUN REPOSITORY GET PRODUCTS
	products, err := h.ProductRepository.FindProducts()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	// CREATE RESPONSE
	var productsResponse []productdto.ProductResponse
	for _, p := range products {
		productResponse := productdto.ProductResponse{
			ID:          p.ID,
			Name:        p.Name,
			Price:       p.Price,
			Description: p.Description,
			Image:       os.Getenv("PATH_FILE") + p.Image,
			Stock:       p.Stock,
		}
		productsResponse = append(productsResponse, productResponse)
	}

	// WRITE RESPONSE
	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: "success", Data: productsResponse}
	json.NewEncoder(w).Encode(response)
}

// FUNCTION GET PRODUCT
func (h *handlerProduct) GetProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// GET PRODUCT ID FROM URL
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	// RUN REPOSITORY GET PRODUCT
	product, err := h.ProductRepository.GetProduct(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	// CREATE RESPONSE
	productResponse := productdto.ProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Price:       product.Price,
		Description: product.Description,
		Image:       os.Getenv("PATH_FILE") + product.Image,
		Stock:       product.Stock,
	}

	// WRITE RESPONSE
	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: "success", Data: productResponse}
	json.NewEncoder(w).Encode(response)
}

// FUNCTION CREATE PRODUCT
func (h *handlerProduct) CreateProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// GET USER ROLE FROM JWT TOKEN
	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	userRole := userInfo["role"]

	// CHECK ROLE ADMIN
	if userRole != "admin" {
		w.WriteHeader(http.StatusUnauthorized)
		response := dto.ErrorResult{Code: http.StatusUnauthorized, Message: "unauthorized"}
		json.NewEncoder(w).Encode(response)
		return
	}

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
	respImage, err := cld.Upload.Upload(ctx, filepath, uploader.UploadParams{Folder: "waysbeans"})

	if err != nil {
		fmt.Println(err.Error())
	}

	// CONVERT REQUEST STRING TO INT
	price, _ := strconv.Atoi(r.FormValue("price"))
	stock, _ := strconv.Atoi(r.FormValue("stock"))

	// GET REQUEST
	request := productdto.ProductRequest{
		Name:        r.FormValue("name"),
		Price:       int32(price),
		Image:       respImage.SecureURL,
		Description: r.FormValue("description"),
		Stock:       int32(stock),
	}

	// VALIDATE REQUEST INPUT
	validation := validator.New()
	err = validation.Struct(request)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	// SETUP FOR QUERY PRODUCT
	product := models.Product{
		Name:        request.Name,
		Price:       request.Price,
		Image:       request.Image,
		Description: request.Description,
		Stock:       request.Stock,
		CreateAt:    time.Now(),
		UpdateAt:    time.Now(),
	}

	// RUN REPOSITORY CREATE PRODUCT
	product, err = h.ProductRepository.CreateProduct(product)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	// CREATE RESPONSE
	productResponse := productdto.ProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Price:       product.Price,
		Description: product.Description,
		Image:       product.Image,
		Stock:       product.Stock,
	}

	// WRITE RESPONSE
	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: "success", Data: productResponse}
	json.NewEncoder(w).Encode(response)
}

// FUNCTION UPDATE PRODUCT
func (h *handlerProduct) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// GET PRODUCT ID FROM URL
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	// GET USER ROLE FROM JWT TOKEN
	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	userRole := userInfo["role"]

	// CHECK ROLE ADMIN
	if userRole != "admin" {
		w.WriteHeader(http.StatusUnauthorized)
		response := dto.ErrorResult{Code: http.StatusUnauthorized, Message: "You're not admin"}
		json.NewEncoder(w).Encode(response)
		return
	}

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

	// GET REQUEST
	price, _ := strconv.Atoi(r.FormValue("price"))
	stock, _ := strconv.Atoi(r.FormValue("stock"))
	request := productdto.UpdateProductRequest{
		Name:        r.FormValue("name"),
		Price:       int32(price),
		Image:       requestImage,
		Description: r.FormValue("description"),
		Stock:       int32(stock),
	}

	// RUN REPOSITORY GET PRODUCT
	product, err := h.ProductRepository.GetProduct(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	// CHECK UPDATE VALUE
	if request.Name != "" {
		product.Name = request.Name
	}

	if request.Price != 0 {
		product.Price = request.Price
	}

	if filepath != "false" {
		product.Image = request.Image
	}

	if request.Stock != 0 {
		product.Stock = request.Stock
	}

	// CHANGE VALUE UPDATE TIME
	product.UpdateAt = time.Now()

	//RUN REPOSITORY UPDATE PRODUCT
	product, err = h.ProductRepository.UpdateProduct(product)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	// CREATE RESPONSE
	updateResponse := productdto.ProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Price:       product.Price,
		Description: product.Description,
		Image:       product.Image,
		Stock:       product.Stock,
	}

	// WRITE RESPONSE
	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: "success", Data: updateResponse}
	json.NewEncoder(w).Encode(response)
}

// FUNCTION DELETE PRODUCT
func (h *handlerProduct) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// GET PRODUCT ID FROM URL
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	// GET USER ROLE FROM JWT TOKEN
	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	userRole := userInfo["role"]

	// CHECK ROLE ADMIN
	if userRole != "admin" {
		w.WriteHeader(http.StatusUnauthorized)
		response := dto.ErrorResult{Code: http.StatusUnauthorized, Message: "You're not admin"}
		json.NewEncoder(w).Encode(response)
		return
	}

	// RUN REPOSITORY GET PRODUCT
	delete, err := h.ProductRepository.GetProduct(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	// RUN REPOSITORY DELETE PRODUCT
	product, err := h.ProductRepository.DeleteProduct(delete)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	// CREATE RESPONSE
	deleteResponse := productdto.DeleteProductResponse{
		ID:   product.ID,
		Name: product.Name,
	}

	// WRITE RESPONSE
	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: "success", Data: deleteResponse}
	json.NewEncoder(w).Encode(response)
}
