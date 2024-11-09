package data

import (
	"github.com/shopspring/decimal"
	"regexp"
	"strings"
	"time"
)

const (
	RegexValue                     = `[a-zA-Z0-9]`
	MatchValue                     = -1
	DecimalValue                   = 1
	RoundValue                     = 2
	RoundDollarPointValue          = 50
	ZeroValue                      = 0
	CentValue                      = 100
	QuarterMultipleValue           = 25
	PairValue                      = 2
	ItemPairValue                  = 5
	ItemDescriptionTrimValue       = 3
	ItemDescriptionMultiplierValue = 0.2
	OddModuloValue                 = 2
	OddValue                       = 1
	OddDayValue                    = 6
	AfterTimeValue                 = 14
	BeforeTimValue                 = 16
	TimeRangeValue                 = 10
)

type Calculator struct {
	Points int32
}

func New() *Calculator {
	return &Calculator{
		Points: 0,
	}
}

func (c *Calculator) AddPoints(points int32) {
	c.Points += points
}

func (c *Calculator) TotalPoints() int32 {
	return c.Points
}

func RetailerNamePoints(retailer string) int32 {
	alphaNumeric := regexp.MustCompile(RegexValue)
	match := alphaNumeric.FindAllString(retailer, MatchValue)

	return int32(len(match))
}

func RoundDollarPoints(total Price) int32 {
	decimalTotal, _ := decimal.NewFromString(total.String())
	roundedTotal := decimalTotal.Round(int32(RoundValue))
	modulo := decimalTotal.Mod(decimal.NewFromInt(int64(DecimalValue)))

	isRoundDollar := roundedTotal.Equal(decimalTotal) && modulo.Equal(decimal.Zero)
	if isRoundDollar {
		return int32(RoundDollarPointValue)
	}

	return int32(ZeroValue)
}

func QuarterMultiplePoints(total Price) int32 {
	decimalTotal, _ := decimal.NewFromString(total.String())
	centsTotal := decimalTotal.Mul(decimal.NewFromInt(int64(CentValue)))
	modulo := centsTotal.Mod(decimal.NewFromInt(int64(QuarterMultipleValue)))

	isQuarterMultiple := modulo.Equal(decimal.Zero)
	if isQuarterMultiple {
		return int32(QuarterMultipleValue)
	}

	return int32(ZeroValue)
}

func ItemPairPoints(items []Item) int32 {
	itemPair := int32(len(items) / PairValue)
	points := itemPair * int32(ItemPairValue)

	return points
}

func ItemDescriptionPoints(items []Item) int32 {
	var points int32

	for _, item := range items {
		trimmed := len(strings.TrimSpace(item.ShortDescription))
		if trimmed%ItemDescriptionTrimValue == 0 {
			itemPrice, _ := decimal.NewFromString(item.Price.String())
			multiplier := decimal.NewFromFloat(ItemDescriptionMultiplierValue)
			points += int32(itemPrice.Mul(multiplier).Ceil().IntPart())
		}
	}

	return points
}

func OddDayPoints(purchaseDate, dateLayout string) int32 {
	if purchaseDate == "" {
		return ZeroValue
	}

	parsedDate, err := time.Parse(dateLayout, purchaseDate)
	if err != nil {
		return ZeroValue
	}

	if parsedDate.Day()%OddModuloValue == OddValue {
		return OddDayValue
	}

	return ZeroValue
}

func PurchaseTimeRangePoints(purchaseTime, timeLayout string) int32 {
	if purchaseTime == "" {
		return ZeroValue
	}

	parsedTime, err := time.Parse(timeLayout, purchaseTime)
	if err != nil {
		return ZeroValue
	}

	if parsedTime.Hour() >= AfterTimeValue && parsedTime.Hour() < BeforeTimValue {
		return TimeRangeValue
	}

	return ZeroValue
}
