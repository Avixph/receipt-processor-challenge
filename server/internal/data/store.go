package data

import (
	"sync"
)

type Stores struct {
	Receipts ReceiptModel
}

func NewStores() Stores {
	return Stores{
		Receipts: ReceiptModel{
			Store: make(map[string]Receipt),
			mu:    &sync.RWMutex{},
		},
	}
}

//type Store struct {
//	mu       sync.RWMutex
//	Receipts map[string]Receipt
//}
//
//func NewStore() *Store {
//	return &Store{
//		Receipts: make(map[string]Receipt),
//	}
//}
