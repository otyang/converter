package converter

import (
	"fmt"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestNewCurrencies(t *testing.T) {
	// Valid data
	validRates := []struct {
		ISOCode  string
		BuyRate  decimal.Decimal
		SellRate decimal.Decimal
	}{
		{"USD", decimal.NewFromInt(100), decimal.NewFromInt(101)},
		{"EUR", decimal.NewFromInt(120), decimal.NewFromInt(121)},
	}

	currencies, err := NewCurrencies(validRates)
	assert.NoError(t, err)
	assert.Equal(t, len(currencies.currencies), 2)

	// Invalid data
	invalidRates := []string{"invalid data"}
	_, err = NewCurrencies(invalidRates)
	assert.Error(t, err)

	// Empty source should return error
	_, err = NewCurrencies[string](nil)
	assert.Error(t, err)
	assert.Equal(t, err, ErrEmptyCurrencySource)

	// Valid source should create currencies instance
	source := []struct{ ISOCode string }{{ISOCode: "USD"}}
	currencies, err = NewCurrencies(source)
	assert.NoError(t, err)
	assert.Equal(t, len(currencies.currencies), 1)
}

func TestFindCurrency(t *testing.T) {
	// Setup Currencies
	currencies := Currencies{
		currencies: []Currency{
			{ISOCode: "USD", BuyRate: decimal.NewFromInt(100), SellRate: decimal.NewFromInt(101)},
			{ISOCode: "EUR", BuyRate: decimal.NewFromInt(120), SellRate: decimal.NewFromInt(121)},
		},
	}

	// Valid code
	currency, err := currencies.FindCurrency("USD")
	assert.NoError(t, err)
	assert.Equal(t, currency.ISOCode, "USD")

	// Invalid code
	_, err = currencies.FindCurrency("GBP")
	assert.Error(t, err)
	assert.Equal(t, fmt.Errorf(ErrCurrencyNotFound, "GBP"), err)

	// Empty currency source
	currencies.currencies = nil
	_, err = currencies.FindCurrency("USD")
	assert.Error(t, err)
	assert.Equal(t, ErrEmptyCurrencySource, err)
}

func TestCalculateRate(t *testing.T) {
	testCases := []struct {
		name      string
		base      string
		from      string
		to        string
		expected  decimal.Decimal
		shouldErr bool
	}{
		{
			name:     "Same currency conversion",
			base:     "USD",
			from:     "USD",
			to:       "USD",
			expected: decimal.NewFromInt(1),
		},
		{
			name:     "Base currency to target currency conversion",
			base:     "USD",
			from:     "USD",
			to:       "EUR",
			expected: decimal.NewFromFloat(1.18),
		},
		{
			name:     "Target currency to base currency conversion",
			base:     "USD",
			from:     "EUR",
			to:       "USD",
			expected: decimal.NewFromFloat(1.18),
		},
		{
			name:     "Cross-rate conversion",
			base:     "USD",
			from:     "EUR",
			to:       "JPY",
			expected: decimal.NewFromFloat(145.24),
		},
		{
			name:      "Non-existent base currency",
			base:      "ZZZ",
			from:      "USD",
			to:        "EUR",
			shouldErr: true,
		},
		{
			name:      "Non-existent from currency",
			base:      "USD",
			from:      "ZZZ",
			to:        "EUR",
			shouldErr: true,
		},
		{
			name:      "Non-existent to currency",
			base:      "USD",
			from:      "EUR",
			to:        "ZZZ",
			shouldErr: true,
		},
	}

	// Test currencies
	currencies := Currencies{
		currencies: []Currency{
			{ISOCode: "USD", Precision: 2, BuyRate: decimal.NewFromFloat(1.0), SellRate: decimal.NewFromFloat(1.0)},
			{ISOCode: "EUR", Precision: 2, BuyRate: decimal.NewFromFloat(0.85), SellRate: decimal.NewFromFloat(1.18)},
			{ISOCode: "JPY", Precision: 2, BuyRate: decimal.NewFromFloat(0.0081), SellRate: decimal.NewFromFloat(123.45)},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualRate, err := currencies.CalculateRate(tc.base, tc.from, tc.to)

			if tc.shouldErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			// assert.Equal(t, tc.expected.String(), actualRate.RoundCeil(2).String())
			assert.True(t, tc.expected.Equal(actualRate.RoundCeil(2)))
		})
	}
}
