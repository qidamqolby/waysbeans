package productdto

type ProductRequest struct {
	Name        string `json:"name" form:"name" validate:"required"`
	Price       int32  `json:"price" form:"price" validate:"required"`
	Image       string `json:"image" form:"image" validate:"required"`
	Description string `json:"description" form:"description" validate:"required"`
	Stock       int32  `json:"stock" form:"stock" validate:"required"`
}

type UpdateProductRequest struct {
	Name        string `json:"name" form:"name"`
	Price       int32  `json:"price" form:"price"`
	Image       string `json:"image" form:"image"`
	Description string `json:"description" form:"description"`
	Stock       int32  `json:"stock" form:"stock"`
}
