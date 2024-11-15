package dal

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"hot-coffee/models"
)

type OrderRepository interface {
	WriteJSONNewOrder(body []models.Order) error
	ReadJSONOrder() ([]models.Order, error)
	PresentInTheInventory(orderQuantity map[string]int, neworderMenu []models.MenuItem) ([]models.InventoryItem, error)
	PresentInTheMenu(neworder models.Order) (map[string]int, []models.MenuItem, error)
	WriteJSONEditIngredients(body []models.InventoryItem) error
	ReadJSONMenu() ([]models.MenuItem, error)
}

type jsonOrderRepository struct{}

// Creates and returns a new instance of jsonOrderRepository
func NewJSONOrderRepository() OrderRepository {
	return &jsonOrderRepository{}
}

// Checks if ordered items exist in the menu and returns the order quantities and menu items
func (r jsonOrderRepository) PresentInTheMenu(neworder models.Order) (map[string]int, []models.MenuItem, error) {
	var newMenuItems []models.MenuItem
	var returnedMenuItems []models.MenuItem
	orderQuantity := make(map[string]int)
	option := os.O_RDONLY
	content, err := os.OpenFile(Menuitems(), option, 0o644)
	defer content.Close()
	if err != nil {
		return nil, nil, err
	}
	reader := json.NewDecoder(content)
	err = reader.Decode(&newMenuItems)
	if err != nil {
		return nil, nil, err
	}
	missMenu := []string{}
	mapMenuItems := make(map[string]models.MenuItem, len(newMenuItems))
	for _, item := range newMenuItems {
		mapMenuItems[item.ID] = item
	}
	for _, orderItem := range neworder.Items {
		if menuitem, exists := mapMenuItems[orderItem.ProductID]; exists {
			orderQuantity[orderItem.ProductID] = orderItem.Quantity
			returnedMenuItems = append(returnedMenuItems, menuitem)
		} else {
			missMenu = append(missMenu, orderItem.ProductID)
		}
	}

	if len(missMenu) > 0 {
		str := fmt.Sprintf("These items are not on the menu: %s", strings.Join(missMenu, ", "))
		return nil, nil, errors.New(str)
	}
	return orderQuantity, returnedMenuItems, nil
}

// Checks if required inventory items are available for an order and updates inventory quantities
func (r jsonOrderRepository) PresentInTheInventory(orderQuantity map[string]int, neworderMenu []models.MenuItem) ([]models.InventoryItem, error) {
	var newInvItems []models.InventoryItem
	var returnedInvItems []models.InventoryItem

	option := os.O_RDONLY
	content, err := os.OpenFile(Inventoryitem(), option, 0o644)
	defer content.Close()
	if err != nil {
		return nil, err
	}
	reader := json.NewDecoder(content)
	err = reader.Decode(&newInvItems)
	if err != nil {
		return nil, err
	}
	missOrder := []string{}

	mapInventoryItems := make(map[string]*models.InventoryItem, len(newInvItems))

	for i := range newInvItems {
		item := &newInvItems[i]
		mapInventoryItems[item.IngredientID] = item
	}
	for _, oneStructOrder := range neworderMenu {
		for _, oneItemOrderIngredient := range oneStructOrder.Ingredients {
			mapstruct, exists := mapInventoryItems[oneItemOrderIngredient.IngredientID]
			quantityOrder, exists1 := orderQuantity[oneStructOrder.ID]

			if exists && exists1 {
				if mapstruct.Quantity-(oneItemOrderIngredient.Quantity*float64(quantityOrder)) < 0 {
					return nil, errors.New("Not enough ingredients")
				}
				mapstruct.Quantity -= (oneItemOrderIngredient.Quantity * float64(quantityOrder))

			} else {
				missOrder = append(missOrder, oneItemOrderIngredient.IngredientID)
			}
		}
	}
	for _, v := range mapInventoryItems {
		returnedInvItems = append(returnedInvItems, *v)
	}
	if len(missOrder) > 0 {
		str := fmt.Sprintf("These items are not in the inventory: %s", strings.Join(missOrder, ", "))
		return nil, errors.New(str)
	}
	return returnedInvItems, nil
}

// Writes a new order to the JSON file and creates a backup in the reserve copy
func (r *jsonOrderRepository) WriteJSONNewOrder(body []models.Order) error {
	option := os.O_CREATE | os.O_WRONLY | os.O_TRUNC
	content, err := os.OpenFile(Orders(), option, 0o644)
	if err != nil {
		return errors.New("")
	}
	defer content.Close()
	encoder := json.NewEncoder(content)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(body)
	if err != nil {
		return err
	}

	newContent, err := os.Open(Orders())
	if err != nil {
		return err
	}
	reserveOrders, err := os.Create(ReserveOrder())
	if err != nil {
		return err
	}

	_, err = io.Copy(reserveOrders, newContent)
	if err != nil {
		return err
	}
	return nil
}

// Reads and decodes order data from the JSON file, returning a slice of orders
func (r *jsonOrderRepository) ReadJSONOrder() ([]models.Order, error) {
	var newOrder []models.Order

	option := os.O_RDONLY | os.O_RDWR
	content, err := os.OpenFile(Orders(), option, 0o644)
	if err != nil {
		return nil, err
	}
	defer content.Close()
	fileInfo, err := content.Stat()
	if err != nil {
		return nil, err
	}
	if fileInfo.Size() == 0 {
		return newOrder, nil
	}
	reader := json.NewDecoder(content)
	err1 := reader.Decode(&newOrder)
	if err1 != nil {
		return nil, err1
	}
	return newOrder, nil
}

// Writes updated inventory items to the JSON file, replacing the current content
func (r *jsonOrderRepository) WriteJSONEditIngredients(body []models.InventoryItem) error {
	option := os.O_CREATE | os.O_WRONLY | os.O_TRUNC
	content, err := os.OpenFile(Inventoryitem(), option, 0o644)
	if err != nil {
		return errors.New("")
	}
	defer content.Close()
	encoder := json.NewEncoder(content)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(body)
	if err != nil {
		return err
	}
	return nil
}

// Reads and decodes menu item data from the JSON file, returning a slice of menu items
func (r *jsonOrderRepository) ReadJSONMenu() ([]models.MenuItem, error) {
	var newMenu []models.MenuItem

	option := os.O_RDONLY | os.O_RDWR
	content, err := os.OpenFile(Menuitems(), option, 0o644)
	if err != nil {
		return nil, err
	}
	defer content.Close()

	reader := json.NewDecoder(content)
	err1 := reader.Decode(&newMenu)
	if err1 != nil {
		return nil, err1
	}
	return newMenu, nil
}
