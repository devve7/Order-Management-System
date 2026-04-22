package product

import "testing"

func mustProductID(t *testing.T, id int64) ProductID {
	t.Helper()
	v, err := NewProductID(id)
	if err != nil {
		t.Fatalf("NewProductID(%d) error = %v", id, err)
	}
	return v
}

func mustProductName(t *testing.T, name string) ProductName {
	t.Helper()
	v, err := NewProductName(name)
	if err != nil {
		t.Fatalf("NewProductName(%q) error = %v", name, err)
	}
	return v
}

func mustPrice(t *testing.T, price int64) Price {
	t.Helper()
	v, err := NewPrice(price)
	if err != nil {
		t.Fatalf("NewPrice(%d) error = %v", price, err)
	}
	return v
}

func mustStock(t *testing.T, stock int64) Stock {
	t.Helper()
	v, err := NewStock(stock)
	if err != nil {
		t.Fatalf("NewStock(%d) error = %v", stock, err)
	}
	return v
}

func TestNewProduct(t *testing.T) {
	p := NewProduct(
		mustProductName(t, "iPhone"),
		mustPrice(t, 100000),
		mustStock(t, 10),
	)

	if p.ID() != 0 {
		t.Fatalf("ID() = %v, want 0", p.ID())
	}
	if p.Name() != "iPhone" {
		t.Fatalf("Name() = %v, want %v", p.Name(), ProductName("iPhone"))
	}
	if p.Price() != 100000 {
		t.Fatalf("Price() = %v, want %v", p.Price(), Price(100000))
	}
	if p.Stock() != 10 {
		t.Fatalf("Stock() = %v, want %v", p.Stock(), Stock(10))
	}
	if !p.Active() {
		t.Fatalf("Active() = %v, want true", p.Active())
	}
	if p.HasID() {
		t.Fatalf("HasID() = %v, want false", p.HasID())
	}
}

func TestRestoreProduct(t *testing.T) {
	p := RestoreProduct(
		mustProductID(t, 7),
		mustProductName(t, "MacBook"),
		mustPrice(t, 200000),
		mustStock(t, 5),
		false,
	)

	if p.ID() != 7 {
		t.Fatalf("ID() = %v, want %v", p.ID(), ProductID(7))
	}
	if p.Name() != "MacBook" {
		t.Fatalf("Name() = %v, want %v", p.Name(), ProductName("MacBook"))
	}
	if p.Price() != 200000 {
		t.Fatalf("Price() = %v, want %v", p.Price(), Price(200000))
	}
	if p.Stock() != 5 {
		t.Fatalf("Stock() = %v, want %v", p.Stock(), Stock(5))
	}
	if p.Active() {
		t.Fatalf("Active() = %v, want false", p.Active())
	}
	if !p.HasID() {
		t.Fatalf("HasID() = %v, want true", p.HasID())
	}
}

func TestProductChangePrice(t *testing.T) {
	p := NewProduct(
		mustProductName(t, "iPhone"),
		mustPrice(t, 100000),
		mustStock(t, 10),
	)

	p.ChangePrice(mustPrice(t, 120000))

	if p.Price() != 120000 {
		t.Fatalf("Price() = %v, want %v", p.Price(), Price(120000))
	}
}

func TestProductAddStock(t *testing.T) {
	p := NewProduct(
		mustProductName(t, "iPhone"),
		mustPrice(t, 100000),
		mustStock(t, 10),
	)

	p.AddStock(mustStock(t, 5))

	if p.Stock() != 15 {
		t.Fatalf("Stock() = %v, want %v", p.Stock(), Stock(15))
	}
}

func TestProductRemoveStock(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		p := NewProduct(
			mustProductName(t, "iPhone"),
			mustPrice(t, 100000),
			mustStock(t, 10),
		)

		err := p.RemoveStock(mustStock(t, 4))
		if err != nil {
			t.Fatalf("RemoveStock() error = %v", err)
		}
		if p.Stock() != 6 {
			t.Fatalf("Stock() = %v, want %v", p.Stock(), Stock(6))
		}
	})

	t.Run("exact amount", func(t *testing.T) {
		p := NewProduct(
			mustProductName(t, "iPhone"),
			mustPrice(t, 100000),
			mustStock(t, 10),
		)

		err := p.RemoveStock(mustStock(t, 10))
		if err != nil {
			t.Fatalf("RemoveStock() error = %v", err)
		}
		if p.Stock() != 0 {
			t.Fatalf("Stock() = %v, want %v", p.Stock(), Stock(0))
		}
	})

	t.Run("insufficient stock", func(t *testing.T) {
		p := NewProduct(
			mustProductName(t, "iPhone"),
			mustPrice(t, 100000),
			mustStock(t, 10),
		)

		err := p.RemoveStock(mustStock(t, 11))
		if err != ErrInsufficientStock {
			t.Fatalf("RemoveStock() error = %v, want %v", err, ErrInsufficientStock)
		}
		if p.Stock() != 10 {
			t.Fatalf("Stock() = %v, want %v", p.Stock(), Stock(10))
		}
	})
}

func TestProductDeactivate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		p := NewProduct(
			mustProductName(t, "iPhone"),
			mustPrice(t, 100000),
			mustStock(t, 10),
		)

		err := p.Deactivate()
		if err != nil {
			t.Fatalf("Deactivate() error = %v", err)
		}
		if p.Active() {
			t.Fatalf("Active() = %v, want false", p.Active())
		}
	})

	t.Run("already inactive", func(t *testing.T) {
		p := RestoreProduct(
			mustProductID(t, 1),
			mustProductName(t, "iPhone"),
			mustPrice(t, 100000),
			mustStock(t, 10),
			false,
		)

		err := p.Deactivate()
		if err != ErrInactiveProduct {
			t.Fatalf("Deactivate() error = %v, want %v", err, ErrInactiveProduct)
		}
		if p.Active() {
			t.Fatalf("Active() = %v, want false", p.Active())
		}
	})
}

func TestProductActivate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		p := RestoreProduct(
			mustProductID(t, 1),
			mustProductName(t, "iPhone"),
			mustPrice(t, 100000),
			mustStock(t, 10),
			false,
		)

		err := p.Activate()
		if err != nil {
			t.Fatalf("Activate() error = %v", err)
		}
		if !p.Active() {
			t.Fatalf("Active() = %v, want true", p.Active())
		}
	})

	t.Run("already active", func(t *testing.T) {
		p := NewProduct(
			mustProductName(t, "iPhone"),
			mustPrice(t, 100000),
			mustStock(t, 10),
		)

		err := p.Activate()
		if err != ErrProductAlreadyActive {
			t.Fatalf("Activate() error = %v, want %v", err, ErrProductAlreadyActive)
		}
		if !p.Active() {
			t.Fatalf("Active() = %v, want true", p.Active())
		}
	})
}

func TestProductEnsureAvailable(t *testing.T) {
	t.Run("active and enough stock", func(t *testing.T) {
		p := NewProduct(
			mustProductName(t, "iPhone"),
			mustPrice(t, 100000),
			mustStock(t, 10),
		)

		err := p.EnsureAvailable(mustStock(t, 5))
		if err != nil {
			t.Fatalf("EnsureAvailable() error = %v", err)
		}
	})

	t.Run("active and exact stock", func(t *testing.T) {
		p := NewProduct(
			mustProductName(t, "iPhone"),
			mustPrice(t, 100000),
			mustStock(t, 10),
		)

		err := p.EnsureAvailable(mustStock(t, 10))
		if err != nil {
			t.Fatalf("EnsureAvailable() error = %v", err)
		}
	})

	t.Run("inactive product", func(t *testing.T) {
		p := RestoreProduct(
			mustProductID(t, 1),
			mustProductName(t, "iPhone"),
			mustPrice(t, 100000),
			mustStock(t, 10),
			false,
		)

		err := p.EnsureAvailable(mustStock(t, 1))
		if err != ErrInactiveProduct {
			t.Fatalf("EnsureAvailable() error = %v, want %v", err, ErrInactiveProduct)
		}
	})

	t.Run("insufficient stock", func(t *testing.T) {
		p := NewProduct(
			mustProductName(t, "iPhone"),
			mustPrice(t, 100000),
			mustStock(t, 10),
		)

		err := p.EnsureAvailable(mustStock(t, 11))
		if err != ErrInsufficientStock {
			t.Fatalf("EnsureAvailable() error = %v, want %v", err, ErrInsufficientStock)
		}
	})
}
