package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()
	//router := flow.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodPost, "/v1/receipts/process", app.processReceiptHandler)
	router.HandlerFunc(http.MethodGet, "/v1/receipts", app.getReceiptListHandler)
	router.HandlerFunc(http.MethodGet, "/v1/receipts/:id", app.getReceiptHandler)
	router.HandlerFunc(http.MethodGet, "/v1/receipts/:id/points", app.getReceiptPointsHandler)
	//router.HandleFunc("/v1/healthcheck", app.healthcheckHandler, "GET")
	//router.HandleFunc("/v1/receipts/process", app.processReceiptHandler, "POST")
	//router.HandleFunc("/v1/receipts/{:id}/points", app.getReceiptHandler, "GET")

	return app.recoverPanic(router)
}
