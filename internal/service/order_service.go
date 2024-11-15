package service

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"hot-coffee/internal/dal"
	"hot-coffee/models"
)

type OrderService interface {
	ServicePostOrders(body models.Order) error
	ServicePutOrderID(id string, newEdit models.Order) error
	CloseOrder(id string) error
	ServiceDeleteOrdersID(id string) error
	GetOrdersService() ([]models.Order, error)
	GetIDOrdersService(id string) (models.Order, error)
	IsItOnTheMenu(body models.Order) error
}

type orderService struct {
	orderRepo dal.OrderRepository
}

var Id int

// Initializes and returns a new instance of orderService with the provided repository
func NewOrderService(orderRepo dal.OrderRepository) OrderService {
	return &orderService{orderRepo: orderRepo}
}

// Creates a new order, validates the order details, and ensures no open orders exist
func (s orderService) ServicePostOrders(body models.Order) error {
	if err := checkBodyOrder(body); err != nil {
		return err
	}
	if err := s.IsItOnTheMenu(body); err != nil {
		return err
	}
	listOrder, err := s.orderRepo.ReadJSONOrder()
	if err != nil {
		return err
	}
	body.ID = orderNumberCreator()
	for _, oneOrder := range listOrder {
		if oneOrder.Status == "open" {
			return errors.New("You already have an open order")
		}
		if oneOrder.ID == body.ID {
			body.ID = orderNumberCreator()
		}
	}
	body.Status = "open"
	nowTime := time.Now()
	timeString := nowTime.Format("2006-01-02 15:04:05")
	body.CreatedAt = timeString
	listOrder = append(listOrder, body)

	if err := s.orderRepo.WriteJSONNewOrder(listOrder); err != nil {
		return err
	}

	return nil
}

// Validates the fields of an order to ensure all required information is present
func checkBodyOrder(body models.Order) error {
	newbodyCustomer := strings.Trim(body.CustomerName, " ")
	if newbodyCustomer == "" {
		return errors.New("Missing customer name")
	}
	if body.Items == nil {
		return errors.New("Missing items in menu")
	}
	for _, item := range body.Items {
		newbodyItemProductID := strings.Trim(item.ProductID, " ")
		if newbodyItemProductID == "" {
			return errors.New("Missing product id")
		}
		if item.Quantity < 1 {
			return errors.New("Quantity cannot be negative")
		}

	}
	return nil
}

// Generates a unique order ID by incrementing a counter
func orderNumberCreator() string {
	var id string
	Id++
	id = fmt.Sprintf("%d", Id)
	return id
}

// Updates an existing order by ID, ensuring it is still open and validating the new data
func (s *orderService) ServicePutOrderID(id string, body models.Order) error {
	if err := s.IsItOnTheMenu(body); err != nil {
		return err
	}
	checker := false
	jsonfilemenu, err := s.orderRepo.ReadJSONOrder()
	if err != nil {
		return err
	}
	for i, oneStructure := range jsonfilemenu {
		if oneStructure.ID == id {

			checker = true
			err, newEditedStructure := s.EditStructureOrder(body, oneStructure)
			if err != nil {
				return err
			}
			if err := checkBodyOrder(newEditedStructure); err != nil {
				return err
			}
			if newEditedStructure.Status != "open" {
				return errors.New("Order with this ID closed")
			}
			jsonfilemenu[i] = newEditedStructure
		}
	}

	if !checker {
		return errors.New("ID not found")
	}
	if err := s.orderRepo.WriteJSONNewOrder(jsonfilemenu); err != nil {
		return err
	}
	return nil
}

// Merges a new order with an existing one, applying updates and validating the result
func (s *orderService) EditStructureOrder(newOrder models.Order, oldOrder models.Order) (error, models.Order) {
	var newEditedStructure models.Order
	newOrderName := strings.Trim(newOrder.CustomerName, " ")
	newEditedStructure = oldOrder
	if newOrderName != "" {
		newEditedStructure.CustomerName = newOrder.CustomerName
	}
	newEditedStructure.Items = newOrder.Items
	err := checkBodyOrder(newEditedStructure)
	if err != nil {
		return err, newEditedStructure
	}

	return nil, newEditedStructure
}

// Closes an open order by ID, updates inventory quantities, and writes changes
func (s *orderService) CloseOrder(id string) error {
	orders, err := s.orderRepo.ReadJSONOrder()
	if err != nil {
		return err
	}
	closeOrder := models.Order{}
	checker := false
	for _, oneOrder := range orders {
		if oneOrder.ID == id {
			if oneOrder.Status == "open" {
				checker = true
				closeOrder = oneOrder

			} else {
				return errors.New("Status not open")
			}
		}
	}
	if !checker {
		return errors.New("ID not found")
	}

	orderQuantity, menuItems, err := s.orderRepo.PresentInTheMenu(closeOrder)
	if err != nil {
		return err
	}
	invItems, err := s.orderRepo.PresentInTheInventory(orderQuantity, menuItems)
	if err != nil {
		return err
	}
	if err := s.orderRepo.WriteJSONEditIngredients(invItems); err != nil {
		return err
	}

	for i, order := range orders {
		if id == order.ID {
			orders[i].Status = "closed"
		}
	}
	err = s.orderRepo.WriteJSONNewOrder(orders)
	if err != nil {
		return err
	}
	return nil
}

// Deletes an order by ID, returning an error if the ID is not found
func (s *orderService) ServiceDeleteOrdersID(id string) error {
	orders, err := s.orderRepo.ReadJSONOrder()
	if err != nil {
		return err
	}
	index := 0
	checker := false
	for i, oneOrder := range orders {
		if oneOrder.ID == id {
			checker = true
			index = i
		}
	}
	if !checker {
		return errors.New("Such ID not found")
	}
	orders = append(orders[:index], orders[index+1:]...)
	if err := s.orderRepo.WriteJSONNewOrder(orders); err != nil {
		return err
	}
	return nil
}

// Retrieves all orders from the repository
func (s *orderService) GetOrdersService() ([]models.Order, error) {
	return s.orderRepo.ReadJSONOrder()
}

// Retrieves a specific order by ID, returning an error if the ID is not found
func (s *orderService) GetIDOrdersService(id string) (models.Order, error) {
	getOrder := models.Order{}
	orders, err := s.orderRepo.ReadJSONOrder()
	if err != nil {
		return getOrder, err
	}
	checker := false
	for _, oneOrder := range orders {
		if id == oneOrder.ID {
			checker = true
			getOrder = oneOrder
		}
	}
	if !checker {
		return getOrder, errors.New("ID not found")
	}
	return getOrder, nil
}

// Checks if all items in an order exist on the menu, returning an error if any item is missing
func (s *orderService) IsItOnTheMenu(body models.Order) error {
	menu, err := s.orderRepo.ReadJSONMenu()
	if err != nil {
		return err
	}
	menuItems := make(map[string]int)
	for _, onemenuItem := range menu {
		menuItems[onemenuItem.ID] = int(onemenuItem.Price)
	}
	list := []string{}
	for _, oneOrderItem := range body.Items {
		_, exists := menuItems[oneOrderItem.ProductID]
		if !exists {
			list = append(list, oneOrderItem.ProductID)
		}
	}
	if len(list) > 0 {
		newList := strings.Join(list, ", ")
		str := fmt.Sprintf("This item is not on the menu: %s", newList)
		return errors.New(str)
	}
	return nil
}
