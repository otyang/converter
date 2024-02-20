# converter

This package provides functionalities for currency conversion and generating quotes. It facilitates handling exchange rates and precision for different currencies.

## Features

* **Currency Definition:** Creates well-defined currency structures with ISO codes, precision, buy rates, and sell rates.
* **Rate Calculation:** Accurately calculates exchange rates between currencies, offering these calculation modes:
    * Same currency conversion
    * Base currency to target currency conversion
    * Target currency to base currency conversion
    * Cross-rate conversion (when neither currency is the base)
* **Quote Generation:** Creates quote objects containing essential information for currency conversion transactions, including:
     * Base currency
     * Source currency ("from")
     * Target currency ("to")
     *  Original amount 
     * Fee
     * Total amount to deduct
     * Exchange rate
     * Final converted amount
     * Date of the quote

## Usage Example

```go
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
```

## Error Handling

The package defines custom errors for handling issues like:

* `ErrCurrencyNotFound`: Indicates a currency could not be found in the data source.
* `ErrEmptyCurrencySource`: Signals that the provided currency data source is empty.
* `ErrBaseCurrencyNotFound`: Indicates the base currency is missing.

 
## License

This package is licensed under the MIT.