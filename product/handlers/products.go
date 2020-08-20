package handlers

import (
	"Product/data"
	"log"
	"net/http"
	"regexp"
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

// ServeHTTP is the main entry point for the handler and staisfies the http.Handler
// interface
func (p *Products) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	// handle the request for a list of products
	if request.Method == http.MethodGet {
		p.getProducts(writer, request)
		return
	}

	// create a product
	if request.Method == http.MethodPost {
		p.addProduct(writer, request)
		return
	}

	// update a product
	if request.Method == http.MethodPut {
		// expect the id in the URI
		regex := regexp.MustCompile(`/([0-9]+)`)
		group := regex.FindAllStringSubmatch(request.URL.Path, -1)

		if len(group) != 1 {
			p.logger.Println("Invalid URI more than one id")
			http.Error(writer, "Invalid URI", http.StatusBadRequest)
			return
		}

		if len(group[0]) != 2 {
			p.logger.Println("Invalid URI more than two capture group")
			http.Error(writer, "Invalid URI", http.StatusBadRequest)
			return
		}

		idString := group[0][1]
		id, err := strconv.Atoi(idString)
		if err != nil {
			p.logger.Println("Unable to convert to number")
			http.Error(writer, "Invalid URI", http.StatusBadRequest)
			return
		}

		p.updateProduct(id, writer, request)
		return
	}

	// catch-all
	// if no method is satisfied return an error
	writer.WriteHeader(http.StatusMethodNotAllowed)
}

// getProducts returns the products from the data store
func (p *Products) getProducts(writer http.ResponseWriter, request *http.Request) {
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

func (p *Products) addProduct(writer http.ResponseWriter, request *http.Request) {
	p.logger.Println("Handle POST Product")

	product := &data.Product{}

	err := product.FromJSON(request.Body)
	if err != nil {
		http.Error(writer, "Unable to unmarshal json", http.StatusBadRequest)
		return
	}

	data.AddProduct(product)
}

func (p *Products) updateProduct(id int, writer http.ResponseWriter, request *http.Request) {
	p.logger.Println("Handle PUT Product")

	product := &data.Product{}

	err := product.FromJSON(request.Body)
	if err != nil {
		http.Error(writer, "Unable to unmarshal json", http.StatusBadRequest)
		return
	}

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
