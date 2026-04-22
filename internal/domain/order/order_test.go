package order

import (
	"testing"
	"time"
)

func mustOrderID(t *testing.T, id int64) OrderID {
	t.Helper()
	v, err := NewOrderID(id)
	if err != nil {
		t.Fatalf("NewOrderID(%d) error = %v", id, err)
	}
	return v
}

func mustCustomerID(t *testing.T, id int64) CustomerID {
	t.Helper()
	v, err := NewCustomerID(id)
	if err != nil {
		t.Fatalf("NewCustomerID(%d) error = %v", id, err)
	}
	return v
}

func mustProductID(t *testing.T, id int64) ProductID {
	t.Helper()
	v, err := NewProductID(id)
	if err != nil {
		t.Fatalf("NewProductID(%d) error = %v", id, err)
	}
	return v
}

func mustPrice(t *testing.T, amount int64) Price {
	t.Helper()
	v, err := NewPrice(amount)
	if err != nil {
		t.Fatalf("NewPrice(%d) error = %v", amount, err)
	}
	return v
}

func mustQuantity(t *testing.T, q int64) Quantity {
	t.Helper()
	v, err := NewQuantity(q)
	if err != nil {
		t.Fatalf("NewQuantity(%d) error = %v", q, err)
	}
	return v
}

func mustProductName(t *testing.T, s string) ProductName {
	t.Helper()
	v, err := NewProductName(s)
	if err != nil {
		t.Fatalf("NewProductName(%q) error = %v", s, err)
	}
	return v
}

func mustItemID(t *testing.T, id int64) ItemID {
	t.Helper()
	v, err := NewItemID(id)
	if err != nil {
		t.Fatalf("NewItemID(%d) error = %v", id, err)
	}
	return v
}

func mustVersion(t *testing.T, v int64) OrderVersion {
	t.Helper()
	res, err := NewOrderVersion(v)
	if err != nil {
		t.Fatalf("NewOrderVersion(%d) error = %v", v, err)
	}
	return res
}

func newTestOrder(t *testing.T) *Order {
	t.Helper()

	return NewOrder(
		mustOrderID(t, 1),
		mustCustomerID(t, 10),
		StatusCreated,
		time.Date(2026, 4, 22, 12, 0, 0, 0, time.UTC),
	)
}

func addTestItem(t *testing.T, o *Order, productID int64, name string, price int64, qty int64) {
	t.Helper()

	err := o.AddItem(
		mustProductID(t, productID),
		mustProductName(t, name),
		mustPrice(t, price),
		mustQuantity(t, qty),
	)
	if err != nil {
		t.Fatalf("AddItem() error = %v", err)
	}
}

func TestNewOrder(t *testing.T) {
	createdAt := time.Date(2026, 4, 22, 12, 0, 0, 0, time.UTC)

	order := NewOrder(
		mustOrderID(t, 1),
		mustCustomerID(t, 2),
		StatusCreated,
		createdAt,
	)

	if order.ID() != 1 {
		t.Fatalf("ID() = %v, want %v", order.ID(), OrderID(1))
	}
	if order.CustomerID() != 2 {
		t.Fatalf("CustomerID() = %v, want %v", order.CustomerID(), CustomerID(2))
	}
	if order.Status() != StatusCreated {
		t.Fatalf("Status() = %v, want %v", order.Status(), StatusCreated)
	}
	if !order.CreatedAt().Equal(createdAt) {
		t.Fatalf("CreatedAt() = %v, want %v", order.CreatedAt(), createdAt)
	}
	if len(order.Items()) != 0 {
		t.Fatalf("len(Items()) = %d, want 0", len(order.Items()))
	}
	if order.NextItemID() != 1 {
		t.Fatalf("NextItemID() = %v, want %v", order.NextItemID(), ItemID(1))
	}
	if order.Version() != 1 {
		t.Fatalf("Version() = %v, want %v", order.Version(), OrderVersion(1))
	}
}

func TestRestoreOrder(t *testing.T) {
	createdAt := time.Date(2026, 4, 22, 12, 0, 0, 0, time.UTC)
	item := NewOrderItem(
		mustItemID(t, 7),
		mustProductID(t, 100),
		mustProductName(t, "apple"),
		mustPrice(t, 150),
		mustQuantity(t, 2),
	)

	order := RestoreOrder(
		mustOrderID(t, 1),
		mustCustomerID(t, 2),
		StatusPaid,
		createdAt,
		mustItemID(t, 8),
		mustVersion(t, 3),
		[]*OrderItem{item},
	)

	if order.ID() != 1 {
		t.Fatalf("ID() = %v, want %v", order.ID(), OrderID(1))
	}
	if order.CustomerID() != 2 {
		t.Fatalf("CustomerID() = %v, want %v", order.CustomerID(), CustomerID(2))
	}
	if order.Status() != StatusPaid {
		t.Fatalf("Status() = %v, want %v", order.Status(), StatusPaid)
	}
	if !order.CreatedAt().Equal(createdAt) {
		t.Fatalf("CreatedAt() = %v, want %v", order.CreatedAt(), createdAt)
	}
	if order.NextItemID() != 8 {
		t.Fatalf("NextItemID() = %v, want %v", order.NextItemID(), ItemID(8))
	}
	if order.Version() != 3 {
		t.Fatalf("Version() = %v, want %v", order.Version(), OrderVersion(3))
	}
	if len(order.Items()) != 1 {
		t.Fatalf("len(Items()) = %d, want 1", len(order.Items()))
	}
}

func TestOrderAddItemSuccess(t *testing.T) {
	order := newTestOrder(t)

	err := order.AddItem(
		mustProductID(t, 100),
		mustProductName(t, "apple"),
		mustPrice(t, 150),
		mustQuantity(t, 2),
	)
	if err != nil {
		t.Fatalf("AddItem() error = %v", err)
	}

	items := order.Items()
	if len(items) != 1 {
		t.Fatalf("len(Items()) = %d, want 1", len(items))
	}
	if items[0].ID() != 1 {
		t.Fatalf("item ID = %v, want %v", items[0].ID(), ItemID(1))
	}
	if order.NextItemID() != 2 {
		t.Fatalf("NextItemID() = %v, want %v", order.NextItemID(), ItemID(2))
	}
}

func TestOrderAddItemFailsForForbiddenStatuses(t *testing.T) {
	tests := []struct {
		name   string
		status OrderStatus
		want   error
	}{
		{"cancelled", StatusCancelled, ErrOrderCancelled},
		{"shipped", StatusShipped, ErrOrderShipped},
		{"paid", StatusPaid, ErrOrderPaid},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			order := NewOrder(
				mustOrderID(t, 1),
				mustCustomerID(t, 2),
				tt.status,
				time.Now(),
			)

			err := order.AddItem(
				mustProductID(t, 100),
				mustProductName(t, "apple"),
				mustPrice(t, 150),
				mustQuantity(t, 2),
			)
			if err != tt.want {
				t.Fatalf("AddItem() error = %v, want %v", err, tt.want)
			}
			if len(order.Items()) != 0 {
				t.Fatalf("len(Items()) = %d, want 0", len(order.Items()))
			}
		})
	}
}

func TestOrderRemoveItemSuccess(t *testing.T) {
	order := newTestOrder(t)
	addTestItem(t, order, 100, "apple", 100, 1)
	addTestItem(t, order, 101, "banana", 200, 2)

	err := order.RemoveItem(1)
	if err != nil {
		t.Fatalf("RemoveItem() error = %v", err)
	}

	items := order.Items()
	if len(items) != 1 {
		t.Fatalf("len(Items()) = %d, want 1", len(items))
	}
	if items[0].ID() != 2 {
		t.Fatalf("remaining item ID = %v, want %v", items[0].ID(), ItemID(2))
	}
}

func TestOrderRemoveItemFailsForForbiddenStatuses(t *testing.T) {
	tests := []struct {
		name   string
		status OrderStatus
		want   error
	}{
		{"cancelled", StatusCancelled, ErrOrderCancelled},
		{"shipped", StatusShipped, ErrOrderShipped},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			order := NewOrder(
				mustOrderID(t, 1),
				mustCustomerID(t, 2),
				tt.status,
				time.Now(),
			)

			err := order.RemoveItem(1)
			if err != tt.want {
				t.Fatalf("RemoveItem() error = %v, want %v", err, tt.want)
			}
		})
	}
}

func TestOrderRemoveItemFailsWhenEmpty(t *testing.T) {
	order := newTestOrder(t)

	err := order.RemoveItem(1)
	if err != ErrOrderItemNotFound {
		t.Fatalf("RemoveItem() error = %v, want %v", err, ErrOrderItemNotFound)
	}
}

func TestOrderRemoveItemFailsWhenNotFound(t *testing.T) {
	order := newTestOrder(t)
	addTestItem(t, order, 100, "apple", 100, 1)

	err := order.RemoveItem(999)
	if err != ErrOrderItemNotFound {
		t.Fatalf("RemoveItem() error = %v, want %v", err, ErrOrderItemNotFound)
	}
}

func TestOrderPay(t *testing.T) {
	t.Run("success from created with items", func(t *testing.T) {
		order := newTestOrder(t)
		addTestItem(t, order, 100, "apple", 100, 1)

		err := order.Pay()
		if err != nil {
			t.Fatalf("Pay() error = %v", err)
		}
		if order.Status() != StatusPaid {
			t.Fatalf("Status() = %v, want %v", order.Status(), StatusPaid)
		}
	})

	t.Run("fails when empty", func(t *testing.T) {
		order := newTestOrder(t)

		err := order.Pay()
		if err != ErrOrderEmpty {
			t.Fatalf("Pay() error = %v, want %v", err, ErrOrderEmpty)
		}
		if order.Status() != StatusCreated {
			t.Fatalf("Status() = %v, want %v", order.Status(), StatusCreated)
		}
	})

	t.Run("fails from non-created status", func(t *testing.T) {
		tests := []OrderStatus{StatusPaid, StatusShipped, StatusCancelled}

		for _, status := range tests {
			t.Run(string(status), func(t *testing.T) {
				order := NewOrder(
					mustOrderID(t, 1),
					mustCustomerID(t, 2),
					StatusCreated,
					time.Now(),
				)
				addTestItem(t, order, 100, "apple", 100, 1)
				order.status = status

				err := order.Pay()
				if err != ErrCannotPay {
					t.Fatalf("Pay() error = %v, want %v", err, ErrCannotPay)
				}
				if order.Status() != status {
					t.Fatalf("Status() = %v, want %v", order.Status(), status)
				}
			})
		}
	})
}

func TestOrderShip(t *testing.T) {
	t.Run("success from paid", func(t *testing.T) {
		order := NewOrder(
			mustOrderID(t, 1),
			mustCustomerID(t, 2),
			StatusPaid,
			time.Now(),
		)

		err := order.Ship()
		if err != nil {
			t.Fatalf("Ship() error = %v", err)
		}
		if order.Status() != StatusShipped {
			t.Fatalf("Status() = %v, want %v", order.Status(), StatusShipped)
		}
	})

	t.Run("fails from other statuses", func(t *testing.T) {
		tests := []OrderStatus{StatusCreated, StatusShipped, StatusCancelled}

		for _, status := range tests {
			t.Run(string(status), func(t *testing.T) {
				order := NewOrder(
					mustOrderID(t, 1),
					mustCustomerID(t, 2),
					status,
					time.Now(),
				)

				err := order.Ship()
				if err != ErrCannotShip {
					t.Fatalf("Ship() error = %v, want %v", err, ErrCannotShip)
				}
				if order.Status() != status {
					t.Fatalf("Status() = %v, want %v", order.Status(), status)
				}
			})
		}
	})
}

func TestOrderCancel(t *testing.T) {
	t.Run("success from created", func(t *testing.T) {
		order := newTestOrder(t)

		err := order.Cancel()
		if err != nil {
			t.Fatalf("Cancel() error = %v", err)
		}
		if order.Status() != StatusCancelled {
			t.Fatalf("Status() = %v, want %v", order.Status(), StatusCancelled)
		}
	})

	t.Run("success from paid", func(t *testing.T) {
		order := NewOrder(
			mustOrderID(t, 1),
			mustCustomerID(t, 2),
			StatusPaid,
			time.Now(),
		)

		err := order.Cancel()
		if err != nil {
			t.Fatalf("Cancel() error = %v", err)
		}
		if order.Status() != StatusCancelled {
			t.Fatalf("Status() = %v, want %v", order.Status(), StatusCancelled)
		}
	})

	t.Run("fails from shipped and cancelled", func(t *testing.T) {
		tests := []OrderStatus{StatusShipped, StatusCancelled}

		for _, status := range tests {
			t.Run(string(status), func(t *testing.T) {
				order := NewOrder(
					mustOrderID(t, 1),
					mustCustomerID(t, 2),
					status,
					time.Now(),
				)

				err := order.Cancel()
				if err != ErrCannotCancel {
					t.Fatalf("Cancel() error = %v, want %v", err, ErrCannotCancel)
				}
				if order.Status() != status {
					t.Fatalf("Status() = %v, want %v", order.Status(), status)
				}
			})
		}
	})
}

func TestOrderTotal(t *testing.T) {
	t.Run("empty order", func(t *testing.T) {
		order := newTestOrder(t)

		got := order.Total()
		if got != 0 {
			t.Fatalf("Total() = %v, want 0", got)
		}
	})

	t.Run("sum of all items", func(t *testing.T) {
		order := newTestOrder(t)
		addTestItem(t, order, 100, "apple", 150, 2)  // 300
		addTestItem(t, order, 101, "banana", 200, 3) // 600

		got := order.Total()
		if got != 900 {
			t.Fatalf("Total() = %v, want %v", got, Price(900))
		}
	})
}

func TestOrderItemsReturnsDeepCopy(t *testing.T) {
	order := newTestOrder(t)
	addTestItem(t, order, 100, "apple", 150, 2)

	items := order.Items()
	items[0] = NewOrderItem(
		mustItemID(t, 999),
		mustProductID(t, 999),
		mustProductName(t, "changed"),
		mustPrice(t, 999),
		mustQuantity(t, 9),
	)

	original := order.Items()
	if len(original) != 1 {
		t.Fatalf("len(Items()) = %d, want 1", len(original))
	}
	if original[0].ID() != 1 {
		t.Fatalf("original item ID = %v, want %v", original[0].ID(), ItemID(1))
	}
}

func TestOrderCloneCreatesDeepCopy(t *testing.T) {
	order := newTestOrder(t)
	addTestItem(t, order, 100, "apple", 150, 2)

	cloned := order.Clone()
	cloned.items[0] = NewOrderItem(
		mustItemID(t, 999),
		mustProductID(t, 999),
		mustProductName(t, "changed"),
		mustPrice(t, 999),
		mustQuantity(t, 9),
	)

	originalItems := order.Items()
	if originalItems[0].ID() != 1 {
		t.Fatalf("original item ID after Clone mutation = %v, want %v", originalItems[0].ID(), ItemID(1))
	}
}

func TestOrderWithNextVersion(t *testing.T) {
	order := RestoreOrder(
		mustOrderID(t, 1),
		mustCustomerID(t, 2),
		StatusCreated,
		time.Now(),
		mustItemID(t, 5),
		mustVersion(t, 7),
		nil,
	)

	next := order.WithNextVersion()

	if next.Version() != 8 {
		t.Fatalf("next.Version() = %v, want %v", next.Version(), OrderVersion(8))
	}
	if order.Version() != 7 {
		t.Fatalf("original Version() = %v, want %v", order.Version(), OrderVersion(7))
	}
}

func TestOrderGetNextItemIDSequenceViaAddItem(t *testing.T) {
	order := newTestOrder(t)

	addTestItem(t, order, 100, "apple", 100, 1)
	addTestItem(t, order, 101, "banana", 200, 1)
	addTestItem(t, order, 102, "pear", 300, 1)

	items := order.Items()
	if len(items) != 3 {
		t.Fatalf("len(Items()) = %d, want 3", len(items))
	}

	if items[0].ID() != 1 || items[1].ID() != 2 || items[2].ID() != 3 {
		t.Fatalf("item IDs = [%v %v %v], want [1 2 3]", items[0].ID(), items[1].ID(), items[2].ID())
	}
	if order.NextItemID() != 4 {
		t.Fatalf("NextItemID() = %v, want %v", order.NextItemID(), ItemID(4))
	}
}
