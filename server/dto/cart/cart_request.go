package cartdto

type CartRequest struct {
	ProductID int32 `json:"id"`
	Quantity  int32 `json:"orderQuantity"`
}
