package main

import (
	"fmt"
	"log"
	"net/http"

	"hot-coffee/internal/dal"
	"hot-coffee/internal/handler"
	"hot-coffee/internal/service"
)

// StartServer initializes and starts the HTTP server on the specified port
func StartServer(port string) error {
	// Ensure the data directory exists
	dal.NewDirectory(*dir)
	// Set up Aggregations: repository, service, and handler
	aggregationsRepo := dal.NewAggregationsRepository()
	aggregationsService := service.NewAggregationsService(aggregationsRepo)
	aggregationsHandler := handler.NewAggregationsHandler(aggregationsService)
	http.HandleFunc("GET /reports/total-sales", aggregationsHandler.TotalSales)
	http.HandleFunc("GET /reports/popular-items", aggregationsHandler.PopularItems)

	// Set up Orders: repository, service, and handler
	orderRepo := dal.NewJSONOrderRepository()
	orderService := service.NewOrderService(orderRepo)
	orderHandler := handler.NewOrderHandler(orderService)
	http.HandleFunc("POST /orders", orderHandler.PostOrders)
	http.HandleFunc("GET /orders", orderHandler.GetOrders)
	http.HandleFunc("GET /orders/{id}", orderHandler.GetOrdersID)
	http.HandleFunc("PUT /orders/{id}", orderHandler.PutOrdersID)
	http.HandleFunc("DELETE /orders/{id}", orderHandler.DeleteOrdersID)
	http.HandleFunc("POST /orders/{id}/close", orderHandler.PostOrdersIDClose)

	// Set up Menu: repository, service, and handler
	menuRepo := dal.NewJSONMenuRepository()
	menuService := service.NewMenuService(menuRepo)
	menuHandler := handler.NewMenuHandler(menuService)
	http.HandleFunc("POST /menu", menuHandler.PostMenu)
	http.HandleFunc("GET /menu", menuHandler.GetMenu)
	http.HandleFunc("GET /menu/{id}", menuHandler.GetMenuID)
	http.HandleFunc("PUT /menu/{id}", menuHandler.PutMenuID)
	http.HandleFunc("DELETE /menu/{id}", menuHandler.DeleteMenuID)

	// Set up Inventory: repository, service, and handler
	invRepo := dal.NewJSONInvRepository()
	invService := service.NewInvService(invRepo)
	invHandler := handler.NewInvHandler(invService)
	http.HandleFunc("POST /inventory", invHandler.PostInv)
	http.HandleFunc("GET /inventory", invHandler.GetInv)
	http.HandleFunc("GET /inventory/{id}", invHandler.GetInvID)
	http.HandleFunc("PUT /inventory/{id}", invHandler.PutInvID)
	http.HandleFunc("DELETE /inventory/{id}", invHandler.DeleteInvID)

	// Set up server port and log the server start
	port = fmt.Sprintf(":%s", port)
	log.Println("Server started on port:", port)
	// Ensure required JSON files are created
	CreatedJSONfile()
	// Start the HTTP server
	return http.ListenAndServe(port, nil)
}
