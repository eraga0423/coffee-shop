package service

import (
	"errors"
	"sort"

	"hot-coffee/internal/dal"
	"hot-coffee/models"
)

type AggregationsService interface {
	ServiceTotalSales() (float64, error)
	ServicePopularItems() (error, []models.Popular)
}

type aggregationsService struct {
	aggregationsRepo dal.AggregationsRepository
}

// Initializes and returns a new instance of aggregationsService with the provided repository
func NewAggregationsService(aggregationsRepo dal.AggregationsRepository) AggregationsService {
	return &aggregationsService{aggregationsRepo: aggregationsRepo}
}

// Calculates the total sales from all orders by summing the cost of each order
func (s *aggregationsService) ServiceTotalSales() (float64, error) {
	err, allorders := s.AllOrders()
	if err != nil {
		return 0, err
	}

	err, total := s.TotalMenu(allorders)
	if err != nil {
		return 0, err
	}

	return total, nil
}

// Aggregates quantities of each product ordered across all orders and returns a map of product quantities
func (s *aggregationsService) AllOrders() (error, map[string]int) {
	orders, err := s.aggregationsRepo.ReadJSONOrder()
	if err != nil {
		return err, nil
	}
	allorders := make(map[string]int)
	for _, oneOrder := range orders {
		for _, orderItem := range oneOrder.Items {
			allorders[orderItem.ProductID] += orderItem.Quantity
		}
	}
	return nil, allorders
}

// Calculates the total revenue by multiplying product quantities with their prices from the menu
func (s *aggregationsService) TotalMenu(allOrders map[string]int) (error, float64) {
	menu, err := s.aggregationsRepo.ReadJSONMenu()
	if err != nil {
		return err, 0
	}
	total := 0.0
	for _, item := range menu {
		quantity, exists := allOrders[item.ID]
		if exists {
			total += float64(quantity) * item.Price
		}
	}
	return nil, total
}

// Finds and returns a sorted list of popular items based on quantities ordered
func (s *aggregationsService) ServicePopularItems() (error, []models.Popular) {
	orders, err := s.aggregationsRepo.ReadJSONOrder()
	if err != nil {
		return err, nil
	}
	result := []models.Popular{}
	tempMap := make(map[string]int)
	for _, oneOrder := range orders {
		for _, item := range oneOrder.Items {
			tempMap[item.ProductID] += item.Quantity
		}
	}
	for id, quantity := range tempMap {
		result = append(result, models.Popular{PopularSales: id, Quantity: quantity})
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Quantity > result[j].Quantity
	})
	if len(result) == 0 {
		return errors.New("No popular items found"), nil
	}

	return nil, result
}
