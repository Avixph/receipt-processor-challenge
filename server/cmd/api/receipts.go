package main

import (
	"errors"
	"fmt"
	"github.com/Avixph/receipt-processor-challenge/server/internal/data"
	"github.com/Avixph/receipt-processor-challenge/server/internal/validator"
	"net/http"
)

// ProcessReceiptHandler for the 'Post /v1/receipts/process' endpoint.
func (app *application) processReceiptHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Retailer     string `json:"retailer"`
		PurchaseDate string `json:"purchaseDate"`
		PurchaseTime string `json:"purchaseTime"`
		Items        []struct {
			ShortDescription string     `json:"shortDescription"`
			Price            data.Price `json:"price"`
		} `json:"items"`
		Total data.Price `json:"total"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	items := make([]data.Item, len(input.Items))
	for i, item := range input.Items {
		items[i] = data.Item{
			ShortDescription: item.ShortDescription,
			Price:            item.Price,
		}
	}
	receipt := &data.Receipt{
		Retailer:     input.Retailer,
		PurchaseDate: input.PurchaseDate,
		PurchaseTime: input.PurchaseTime,
		Items:        items,
		Total:        input.Total,
	}

	v := validator.New()
	if data.ValidateReceipt(v, receipt); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.store.Receipts.Insert(receipt)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/receipts/%d", receipt.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"points": receipt}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// GetReceiptHandler for the 'Get /v1/receipts' endpoint.
func (app *application) getReceiptListHandler(w http.ResponseWriter, r *http.Request) {
	receipts, err := app.store.Receipts.GetAll()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"receipts": receipts}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// GetReceiptHandler for the 'Get /v1/receipts/:id' endpoint.
func (app *application) getReceiptHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.realIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	receipt, err := app.store.Receipts.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"receipt": receipt}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// GetReceiptPointsHandler for the 'Get /v1/receipts/:id/points' endpoint.
func (app *application) getReceiptPointsHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.realIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	receipt, err := app.store.Receipts.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"points": receipt.Points}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
