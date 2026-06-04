package product

type ProductDTO struct {
	ID     int64
	Name   string
	Price  int64
	Stock  int64
	Active bool
}

type ProductListParamsDTO struct {
	Cursor *int64
	Limit  *int64
}
