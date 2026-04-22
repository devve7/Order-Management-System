package order

import (
	domain_order "Order-Management-System/internal/domain/order"
	"testing"
	"time"
)

func TestToOrderDTO(t *testing.T) {
	createdAt := time.Date(2026, 4, 22, 12, 0, 0, 0, time.UTC)
	order := domain_order.NewOrder(1, 10, domain_order.StatusCreated, createdAt)

	err := order.AddItem(100, "apple", 150, 2) // 300
	if err != nil {
		t.Fatalf("order.AddItem() error = %v", err)
	}
	err = order.AddItem(101, "banana", 200, 1) // 200
	if err != nil {
		t.Fatalf("order.AddItem() error = %v", err)
	}

	dto := ToOrderDTO(order)

	if dto.ID != 1 {
		t.Fatalf("dto.ID = %d, want 1", dto.ID)
	}
	if dto.CustomerID != 10 {
		t.Fatalf("dto.CustomerID = %d, want 10", dto.CustomerID)
	}
	if dto.Status != "created" {
		t.Fatalf("dto.Status = %q, want %q", dto.Status, "created")
	}
	if !dto.CreatedAt.Equal(createdAt) {
		t.Fatalf("dto.CreatedAt = %v, want %v", dto.CreatedAt, createdAt)
	}
	if len(dto.Items) != 2 {
		t.Fatalf("len(dto.Items) = %d, want 2", len(dto.Items))
	}
	if dto.Items[0].ID != 1 {
		t.Fatalf("dto.Items[0].ID = %d, want 1", dto.Items[0].ID)
	}
	if dto.Items[0].ProductID != 100 {
		t.Fatalf("dto.Items[0].ProductID = %d, want 100", dto.Items[0].ProductID)
	}
	if dto.Items[0].Name != "apple" {
		t.Fatalf("dto.Items[0].Name = %q, want %q", dto.Items[0].Name, "apple")
	}
	if dto.Items[0].Price != 150 {
		t.Fatalf("dto.Items[0].Price = %d, want 150", dto.Items[0].Price)
	}
	if dto.Items[0].Quantity != 2 {
		t.Fatalf("dto.Items[0].Quantity = %d, want 2", dto.Items[0].Quantity)
	}
	if dto.Total != 500 {
		t.Fatalf("dto.Total = %d, want 500", dto.Total)
	}
}

func TestToOrderDTOEmptyItems(t *testing.T) {
	createdAt := time.Date(2026, 4, 22, 12, 0, 0, 0, time.UTC)
	order := domain_order.NewOrder(1, 10, domain_order.StatusCreated, createdAt)

	dto := ToOrderDTO(order)

	if len(dto.Items) != 0 {
		t.Fatalf("len(dto.Items) = %d, want 0", len(dto.Items))
	}
	if dto.Total != 0 {
		t.Fatalf("dto.Total = %d, want 0", dto.Total)
	}
}
