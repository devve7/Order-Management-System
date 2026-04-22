package order

import "testing"

func TestNewOrderItemAndGetters(t *testing.T) {
	item := NewOrderItem(
		1,
		10,
		"apple",
		150,
		2,
	)

	if item.ID() != 1 {
		t.Fatalf("ID() = %v, want %v", item.ID(), ItemID(1))
	}
	if item.ProductID() != 10 {
		t.Fatalf("ProductID() = %v, want %v", item.ProductID(), ProductID(10))
	}
	if item.Name() != "apple" {
		t.Fatalf("Name() = %v, want %v", item.Name(), ProductName("apple"))
	}
	if item.Price() != 150 {
		t.Fatalf("Price() = %v, want %v", item.Price(), Price(150))
	}
	if item.Quantity() != 2 {
		t.Fatalf("Quantity() = %v, want %v", item.Quantity(), Quantity(2))
	}
}

func TestOrderItemHasID(t *testing.T) {
	tests := []struct {
		name string
		item OrderItem
		want bool
	}{
		{
			name: "has positive id",
			item: OrderItem{id: 1},
			want: true,
		},
		{
			name: "zero id",
			item: OrderItem{id: 0},
			want: false,
		},
		{
			name: "negative id",
			item: OrderItem{id: -1},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.item.HasID()
			if got != tt.want {
				t.Fatalf("HasID() = %v, want %v", got, tt.want)
			}
		})
	}
}
