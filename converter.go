package converter

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/shopspring/decimal"
)

// Currency errors
var (
	ErrCurrencyNotFound     = "currency %s not found"
	ErrEmptyCurrencySource  = errors.New("empty currency source: no rates or currency")
	ErrBaseCurrencyNotFound = errors.New("base currency not found")
)

// Currency structure
type Currency struct {
	ISOCode   string          `json:"isoCode"`
	Precision int             `json:"precision"`
	BuyRate   decimal.Decimal `json:"buyRate"`
	SellRate  decimal.Decimal `json:"sellRate"`
}

// Currencies structure
type Currencies struct {
	mutex      sync.RWMutex
	currencies []Currency // Maps ISO code to currency
}

// NewCurrencies creates a Currencies instance from a source of rates.
func NewCurrencies[T any](sourceRates []T) (*Currencies, error) {
	b, err := json.Marshal(sourceRates)
	if err != nil {
		return nil, err
	}

	var currencies []Currency
	if err := json.Unmarshal(b, &currencies); err != nil {
		return nil, err
	}

	if len(currencies) == 0 {
		return nil, ErrEmptyCurrencySource
	}

	return &Currencies{currencies: currencies}, nil
}

// FindCurrency finds a currency by its ISO code.
func (c *Currencies) FindCurrency(code string) (*Currency, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if len(c.currencies) == 0 {
		return nil, ErrEmptyCurrencySource
	}

	for i := range c.currencies {
		if strings.EqualFold(c.currencies[i].ISOCode, code) {
			return &c.currencies[i], nil
		}
	}

	return nil, fmt.Errorf(ErrCurrencyNotFound, code)
}

// CalculateRate calculates the exchange rate between two currencies.
// Same Currency Conversion:
//   - from 	(you have) 			= base currency
//   - to 		(you want) 			= base currency
//   - Rate: 	1
//
// Base to Target Conversion:
//   - from 	(you have) 			= base currency
//   - to 		(you want/target) 	= another currency
//   - Rate: 	sell-rate of to
//
// Target to Base:
//   - from 	(you have/target) 	= another currency
//   - to 		(you want) 			= base currency
//   - Rate: 	1 / sell-rate of from
//
// Cross rate conversion: [Target to base and then to target]
//   - from 	(you have/source) 	= another currency
//   - to 		(you want/target) 	= another currency
//   - Rate: 	[Target to Base of: from] * [Base to Target of: to]
func (c *Currencies) CalculateRate(baseCurrency, from, to string) (decimal.Decimal, error) {
	baseCurrency = strings.ToUpper(strings.ToUpper(baseCurrency))
	from = strings.ToUpper(strings.ToUpper(from))
	to = strings.ToUpper(strings.ToUpper(to))

	// Same Currency Conversion
	if from == to {
		return decimal.NewFromInt(1), nil
	}

	_, err := c.FindCurrency(baseCurrency)
	if err != nil {
		return decimal.Zero, ErrBaseCurrencyNotFound
	}

	// Base to Target Currency (Sell Rate)
	if from == baseCurrency {
		toCurrency, err := c.FindCurrency(to)
		if err != nil {
			return decimal.Zero, err
		}

		return toCurrency.SellRate, nil
	}

	// Target to Base Currency (Buy Rate)
	if to == baseCurrency {
		fromCurrency, err := c.FindCurrency(from)
		if err != nil {
			return decimal.Zero, err
		}
		return decimal.NewFromInt(1).Div(fromCurrency.BuyRate), nil
	}

	// Cross Rate Conversion
	fromCurrency, err := c.FindCurrency(from)
	if err != nil {
		return decimal.Zero, err
	}
	toCurrency, err := c.FindCurrency(to)
	if err != nil {
		return decimal.Zero, err
	}

	// (target to base) to target
	return (decimal.NewFromInt(1).Div(fromCurrency.BuyRate)).Mul(toCurrency.SellRate), nil
}

// Quote structure
type Quote struct {
	BaseCurrency    string          `json:"baseCurrency"`
	FromCurrency    string          `json:"fromCurrency"`
	FromAmount      decimal.Decimal `json:"fromAmount"`
	Fee             decimal.Decimal `json:"fee"`
	AmountToDeduct decimal.Decimal `json:"amountToDeduct"`
	Rate            decimal.Decimal `json:"rate"`
	ToCurrency      string          `json:"toCurrency"`
	FinalAmount     decimal.Decimal `json:"totalAmount"`
	Date            time.Time       `json:"date"`
}

// NewQuote creates a new quote object.
func NewQuote(
	rateSource *Currencies,
	baseCurrency, fromCurrency, toCurrency string,
	fromAmount,
	fee decimal.Decimal,
) (*Quote, error) {
	if rateSource == nil {
		return nil, errors.New("currency object empty. shouldnt be")
	}

	rate, err := rateSource.CalculateRate(baseCurrency, fromCurrency, toCurrency)
	if err != nil {
		return nil, err
	}

	infoFrom, err := rateSource.FindCurrency(fromCurrency)
	if err != nil {
		return nil, err
	}

	infoTo, err := rateSource.FindCurrency(toCurrency)
	if err != nil {
		return nil, err
	}

	return &Quote{
		BaseCurrency:    baseCurrency,
		FromCurrency:    fromCurrency,
		FromAmount:      fromAmount,
		Fee:             fee,
		AmountToDeduct:  fromAmount.Add(fee).RoundCeil(int32(infoFrom.Precision)),
		Rate:            rate,
		ToCurrency:      toCurrency,
		FinalAmount:     fromAmount.Mul(rate).RoundCeil(int32(infoTo.Precision)),
		Date:            time.Now(),
	}, nil
}
