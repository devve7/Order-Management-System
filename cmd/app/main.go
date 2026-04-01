package main

import (
	application_order "Order-Management-System/internal/application/order"
	domain_order "Order-Management-System/internal/domain/order"
	inmemory_catalog "Order-Management-System/internal/infractructure/catalog/inmemory"
	inmemory_storage "Order-Management-System/internal/infractructure/storage/inmemory"
	transport_http "Order-Management-System/internal/transport/http"
)

func fieldCatalog(catalog *inmemory_catalog.InMemoryOrderCatalog) {
	catalog.AddProduct(domain_order.Product{
		Name:  "Iphone 13",
		Price: 1000,
		ID:    1,
	})
	catalog.AddProduct(domain_order.Product{
		Name:  "Iphone 14",
		Price: 1200,
		ID:    2,
	})
	catalog.AddProduct(domain_order.Product{
		Name:  "Iphone 15",
		Price: 1600,
		ID:    3,
	})
}

func main() {
	catalog := inmemory_catalog.NewInMemoryOrderCatalog()
	fieldCatalog(catalog)
	factory := domain_order.NewOrderItemFactory(catalog)
	repo := inmemory_storage.NewInmemoryOrderRepository()
	usecase := application_order.NewUseCase(factory, repo)
	handler := transport_http.NewHandler(usecase)
	router := transport_http.NewRouter(handler)
	server := transport_http.NewHTTPServer(router)
	server.Start(":9091")
}
