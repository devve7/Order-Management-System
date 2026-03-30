package order

type Repository interface {
	Save(order *Order) error

	Get(id OrderID) (*Order, error)

	Update(order *Order) error

	NextID() (OrderID, error)

	GetAll() ([]*Order, error)
}
