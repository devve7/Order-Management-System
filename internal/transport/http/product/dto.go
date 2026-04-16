package product

import application_product "Order-Management-System/internal/application/product"

type ErrorResponse struct {
	Error string `json:"error"`
}

type CreateProductDTO struct {
	Name  string `json:"name"`
	Price *int64 `json:"price_cents"`
	Stock *int64 `json:"stock"`
}

type ProductIDDTO struct {
	ProductID int64 `json:"product_id"`
}

type ProductDTO struct {
	ID     int64
	Name   string
	Price  int64
	Stock  int64
	Active bool
}

func ToProductDTO(product *application_product.ProductDTO) ProductDTO {
	return ProductDTO{
		ID:     product.ID,
		Name:   product.Name,
		Price:  product.Price,
		Stock:  product.Stock,
		Active: product.Active,
	}
}

type PriceDTO struct {
	Price int64 `json:"price"`
}

type StockAmountDTO struct {
	StockAmount int64 `json:"amount"`
}
