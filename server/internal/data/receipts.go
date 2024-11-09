package data

import (
	"github.com/Avixph/receipt-processor-challenge/server/internal/validator"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"html"
	"sync"
	"time"
)

type Item struct {
	ShortDescription string `json:"shortDescription"`
	Price            Price  `json:"price"`
}

type Receipt struct {
	ID           uuid.UUID `json:"id,string"`
	CreatedAt    time.Time `json:"-"`
	Retailer     string    `json:"retailer"`
	PurchaseDate string    `json:"purchaseDate"`
	PurchaseTime string    `json:"purchaseTime"`
	Items        []Item    `json:"items"`
	Total        Price     `json:"total"`
	Points       int32     `json:"points"`
	Version      int32     `json:"version"`
}

func ValidateReceipt(v *validator.Validator, receipt *Receipt) {
	v.Check(receipt.Retailer != "", "retailer", "must be provided")
	v.Check(len(receipt.Retailer) <= 500, "retailer", "must not be more than 500 bytes long")
	v.Check(validator.TimeFormat(receipt.PurchaseDate, "2006-01-02"), "purchaseDate", "must be in the format YYYY-MM-DD")
	v.Check(validator.TimeFormat(receipt.PurchaseTime, "15:04"), "purchaseTime", "must be in the format HH:MM")
	v.Check(receipt.Items != nil, "items", "must be provided")
	for _, item := range receipt.Items {
		v.Check(item.ShortDescription != "", "shortDescription", "must be provided")
		v.Check(!item.Price.Equal(decimal.Zero), "item price", "must be provided")
		v.Check(item.Price.GreaterThan(decimal.Zero) && item.Price.IsPositive(), "price", "must be positive")
	}
	v.Check(!receipt.Total.Equal(decimal.Zero), "total", "must be provided")
	v.Check(receipt.Total.GreaterThan(decimal.Zero) && receipt.Total.IsPositive(), "total", "must be positive")
}

func CalculatePoints(c *Calculator, receipt *Receipt) int32 {
	c.AddPoints(RetailerNamePoints(receipt.Retailer))
	c.AddPoints(RoundDollarPoints(receipt.Total))
	c.AddPoints(QuarterMultiplePoints(receipt.Total))
	c.AddPoints(ItemPairPoints(receipt.Items))
	c.AddPoints(ItemDescriptionPoints(receipt.Items))
	c.AddPoints(OddDayPoints(receipt.PurchaseDate, "2006-01-02"))
	c.AddPoints(PurchaseTimeRangePoints(receipt.PurchaseTime, "15:04"))

	return c.TotalPoints()
}

type ReceiptModel struct {
	Store map[string]Receipt
	mu    *sync.RWMutex
}

func (m ReceiptModel) Insert(receipt *Receipt) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	c := New()

	receipt.ID = uuid.New()
	receipt.Retailer = html.UnescapeString(receipt.Retailer)
	receipt.CreatedAt = time.Now()
	receipt.Points = CalculatePoints(c, receipt)
	receipt.Version += 1

	m.Store[receipt.ID.String()] = *receipt
	return nil
}

func (m ReceiptModel) GetAll() ([]*Receipt, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	receipts := make([]*Receipt, 0, len(m.Store))
	for _, receipt := range m.Store {
		receipts = append(receipts, &receipt)
	}

	return receipts, nil
}

func (m ReceiptModel) Get(id uuid.UUID) (*Receipt, error) {
	if id == uuid.Nil {
		return nil, ErrRecordNotFound
	}
	m.mu.RLock()
	defer m.mu.RUnlock()

	receipt, exists := m.Store[id.String()]
	if !exists {
		return nil, ErrRecordNotFound
	}

	return &receipt, nil
}
