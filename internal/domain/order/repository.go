package order

type Repository interface {
	Save(order *Order) error

	Get(id OrderID) (*Order, error)

	NextID() (OrderID, error)
}
