package handlers

import (
	"Product/data"
	"context"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

// Products is a http.Handler
type Products struct {
	logger *log.Logger
}

// NewProducts creates a products handler with the given logger
func NewProducts(logger *log.Logger) *Products {
	return &Products{logger}
}

// getProducts returns the products from the data store
func (p *Products) GetProducts(writer http.ResponseWriter, request *http.Request) {
	p.logger.Println("Handle GET Products")

	// fetch the products from the datastore
	products := data.GetProducts()

	// serialize the list to JSON
	err := products.ToJSON(writer)
	if err != nil {
		http.Error(writer, "Unable to marshal json", http.StatusInternalServerError)
		return
	}
}

func (p *Products) AddProduct(writer http.ResponseWriter, request *http.Request) {
	p.logger.Println("Handle POST Product")

	product := request.Context().Value(KeyProduct{}).(*data.Product)
	data.AddProduct(product)
}

func (p *Products) UpdateProduct(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(writer, "Unable to parse ID to int", http.StatusBadRequest)
		return
	}

	p.logger.Println("Handle PUT Product", id)
	product := request.Context().Value(KeyProduct{}).(*data.Product)

	err = data.UpdateProduct(id, product)
	if err == data.ErrProductNotFound {
		http.Error(writer, "Product not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(writer, "Product not found", http.StatusInternalServerError)
		return
	}
}

type KeyProduct struct{}

func (p *Products) MiddlewareProductValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		product := &data.Product{}

		err := product.FromJSON(request.Body)
		if err != nil {
			http.Error(writer, "Unable to unmarshal json", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(request.Context(), KeyProduct{}, product)
		newRequest := request.WithContext(ctx)
		next.ServeHTTP(writer, newRequest)
	})
}
