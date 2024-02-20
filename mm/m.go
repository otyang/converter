package main

import (
	"fmt"

	"github.com/otyang/converter"
	"github.com/shopspring/decimal"
)

func main() {
	// Sample currencies
	currencies := []converter.Currency{
		{ISOCode: "USD", Precision: 2, BuyRate: decimal.NewFromFloat(1), SellRate: decimal.NewFromFloat(1)},
		{ISOCode: "EUR", Precision: 2, BuyRate: decimal.NewFromFloat(0.9), SellRate: decimal.NewFromFloat(0.95)},
		{ISOCode: "NGN", Precision: 2, BuyRate: decimal.NewFromFloat(450), SellRate: decimal.NewFromFloat(460)},
	}

	// Calculate rate from USD to EUR
	rate, err := converter.CalculateRate(currencies, "USD", "USD", "EUR")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("USD to EUR exchange rate:", rate)

	// Create a new quote
	quote, err := converter.NewQuote(currencies, "USD", "USD", "NGN", decimal.NewFromFloat(100), decimal.NewFromFloat(5))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Quote:", quote)
}
