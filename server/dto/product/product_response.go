package productdto

type ProductResponse struct {
	ID          int32  `json:"id"`
	Name        string `json:"name"`
	Price       int32  `json:"price"`
	Description string `json:"description"`
	Image       string `json:"image"`
	Stock       int32  `json:"stock"`
}

type DeleteProductResponse struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
}
