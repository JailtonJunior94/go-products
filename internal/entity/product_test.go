package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProduct(t *testing.T) {
	testCases := []struct {
		name     string
		input    Product
		expected func(p *Product, err error)
	}{
		{
			name:  "Should created product with success",
			input: Product{Name: "Product 1", Price: 10.00},
			expected: func(p *Product, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, p)
				assert.NotEmpty(t, p.ID)
				assert.Equal(t, "Product 1", p.Name)
				assert.Equal(t, 10.00, p.Price)
				assert.Nil(t, p.Validate())
			},
		},
		{
			name:  "Should validated when name is required",
			input: Product{Name: "", Price: 10.00},
			expected: func(p *Product, err error) {
				assert.Nil(t, p)
				assert.Equal(t, ErrNameIsRequired, err)
			},
		},
		{
			name:  "Should validated when price is required",
			input: Product{Name: "Product", Price: 0},
			expected: func(p *Product, err error) {
				assert.Nil(t, p)
				assert.Equal(t, ErrPriceIsRequired, err)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			product, err := NewProduct(tc.input.Name, tc.input.Price)
			tc.expected(product, err)
		})
	}
}
