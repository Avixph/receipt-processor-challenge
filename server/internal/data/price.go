package data

import (
	"errors"
	"github.com/shopspring/decimal"
	"strconv"
)

var ErrInvalidPriceFormat = errors.New("invalid price format")

type Price struct {
	decimal.Decimal
}

//// NewPrice creates a Price from a float64
//func NewPrice(value float64) Price {
//	return Price{decimal.NewFromFloat(value)}
//}

// UnmarshalJSON is a method on the Price type that should return the JSON-decoded
// value for the receipt item price and total.
func (p *Price) UnmarshalJSON(jsonValue []byte) error {
	unquotedJSONValue, err := strconv.Unquote(string(jsonValue))

	if err != nil {
		return ErrInvalidPriceFormat
	}

	i, err := decimal.NewFromString(unquotedJSONValue)
	if err != nil {
		return ErrInvalidPriceFormat
	}

	p.Decimal = i

	return nil
}

// MarshalJSON is a method on the Price type that should return the JSON-encoded
// value for the receipt item price and total.
func (p *Price) MarshalJSON() ([]byte, error) {
	jsnValue := p.Decimal.StringFixed(2)
	quotedJsnValue := strconv.Quote(jsnValue)

	return []byte(quotedJsnValue), nil
}
