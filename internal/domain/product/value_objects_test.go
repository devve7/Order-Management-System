package product

import "testing"

func TestNewProductID(t *testing.T) {
	tests := []struct {
		name    string
		input   int64
		want    ProductID
		wantErr error
	}{
		{"positive", 1, 1, nil},
		{"many", 42, 42, nil},
		{"zero", 0, 0, ErrInvalidProductID},
		{"negative", -1, 0, ErrInvalidProductID},
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

func TestNewProductName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    ProductName
		wantErr error
	}{
		{"valid", "iPhone", "iPhone", nil},
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

func TestNewStock(t *testing.T) {
	tests := []struct {
		name    string
		input   int64
		want    Stock
		wantErr error
	}{
		{"positive", 10, 10, nil},
		{"zero", 0, 0, nil},
		{"negative", -1, 0, ErrInvalidStock},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewStock(tt.input)
			if err != tt.wantErr {
				t.Fatalf("NewStock(%d) error = %v, want %v", tt.input, err, tt.wantErr)
			}
			if got != tt.want {
				t.Fatalf("NewStock(%d) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}
