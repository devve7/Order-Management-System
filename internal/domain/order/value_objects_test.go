package order

import "testing"

func TestNewOrderID(t *testing.T) {
	tests := []struct {
		name    string
		input   int64
		want    OrderID
		wantErr error
	}{
		{"positive", 10, 10, nil},
		{"zero", 0, 0, nil},
		{"negative", -1, 0, ErrInvalidID},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewOrderID(tt.input)
			if err != tt.wantErr {
				t.Fatalf("NewOrderID(%d) error = %v, want %v", tt.input, err, tt.wantErr)
			}
			if got != tt.want {
				t.Fatalf("NewOrderID(%d) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestNewCustomerID(t *testing.T) {
	tests := []struct {
		name    string
		input   int64
		want    CustomerID
		wantErr error
	}{
		{"positive", 10, 10, nil},
		{"zero", 0, 0, nil},
		{"negative", -1, 0, ErrInvalidID},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewCustomerID(tt.input)
			if err != tt.wantErr {
				t.Fatalf("NewCustomerID(%d) error = %v, want %v", tt.input, err, tt.wantErr)
			}
			if got != tt.want {
				t.Fatalf("NewCustomerID(%d) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestNewItemID(t *testing.T) {
	tests := []struct {
		name    string
		input   int64
		want    ItemID
		wantErr error
	}{
		{"positive", 10, 10, nil},
		{"zero", 0, 0, nil},
		{"negative", -1, 0, ErrInvalidID},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewItemID(tt.input)
			if err != tt.wantErr {
				t.Fatalf("NewItemID(%d) error = %v, want %v", tt.input, err, tt.wantErr)
			}
			if got != tt.want {
				t.Fatalf("NewItemID(%d) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestNewProductID(t *testing.T) {
	tests := []struct {
		name    string
		input   int64
		want    ProductID
		wantErr error
	}{
		{"positive", 10, 10, nil},
		{"zero", 0, 0, nil},
		{"negative", -1, 0, ErrInvalidID},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewProductID(tt.input)
			if err != tt.wantErr {
				t.Fatalf("NewProductID(%d) error = %v, want %v", tt.input, err, tt.wantErr)
			}
			if got != tt.want {
				t.Fatalf("NewProductID(%d) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestNewPrice(t *testing.T) {
	tests := []struct {
		name    string
		input   int64
		want    Price
		wantErr error
	}{
		{"positive", 100, 100, nil},
		{"zero", 0, 0, nil},
		{"negative", -1, 0, ErrInvalidPrice},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewPrice(tt.input)
			if err != tt.wantErr {
				t.Fatalf("NewPrice(%d) error = %v, want %v", tt.input, err, tt.wantErr)
			}
			if got != tt.want {
				t.Fatalf("NewPrice(%d) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestNewQuantity(t *testing.T) {
	tests := []struct {
		name    string
		input   int64
		want    Quantity
		wantErr error
	}{
		{"positive", 1, 1, nil},
		{"many", 10, 10, nil},
		{"zero", 0, 0, ErrInvalidQuantity},
		{"negative", -1, 0, ErrInvalidQuantity},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewQuantity(tt.input)
			if err != tt.wantErr {
				t.Fatalf("NewQuantity(%d) error = %v, want %v", tt.input, err, tt.wantErr)
			}
			if got != tt.want {
				t.Fatalf("NewQuantity(%d) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestNewProductName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    ProductName
		wantErr error
	}{
		{"valid", "apple", "apple", nil},
		{"empty", "", "", ErrInvalidProductName},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewProductName(tt.input)
			if err != tt.wantErr {
				t.Fatalf("NewProductName(%q) error = %v, want %v", tt.input, err, tt.wantErr)
			}
			if got != tt.want {
				t.Fatalf("NewProductName(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestNewOrderVersion(t *testing.T) {
	tests := []struct {
		name    string
		input   int64
		want    OrderVersion
		wantErr error
	}{
		{"positive", 1, 1, nil},
		{"many", 10, 10, nil},
		{"zero", 0, 0, ErrInvalidOrderVersion},
		{"negative", -1, 0, ErrInvalidOrderVersion},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewOrderVersion(tt.input)
			if err != tt.wantErr {
				t.Fatalf("NewOrderVersion(%d) error = %v, want %v", tt.input, err, tt.wantErr)
			}
			if got != tt.want {
				t.Fatalf("NewOrderVersion(%d) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}
