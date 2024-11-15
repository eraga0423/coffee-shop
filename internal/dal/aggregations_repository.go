package dal

import (
	"encoding/json"
	"os"

	"hot-coffee/models"
)

// AggregationsRepository defines the interface for reading JSON data for orders and menu items
type AggregationsRepository interface {
	ReadJSONOrder() ([]models.Order, error)
	ReadJSONMenu() ([]models.MenuItem, error)
}

type aggregationsRepository struct{}

// NewAggregationsRepository creates and returns a new instance of aggregationsRepository
func NewAggregationsRepository() AggregationsRepository {
	return &aggregationsRepository{}
}

// ReadJSONOrder reads and decodes order data from the JSON file, returning a slice of orders
func (r *aggregationsRepository) ReadJSONOrder() ([]models.Order, error) {
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

// ReadJSONMenu reads and decodes menu item data from the JSON file, returning a slice of menu items
func (r *aggregationsRepository) ReadJSONMenu() ([]models.MenuItem, error) {
	var menu []models.MenuItem

	option := os.O_RDONLY | os.O_RDWR
	content, err := os.OpenFile(Menuitems(), option, 0o644)
	if err != nil {
		return nil, err
	}
	defer content.Close()
	fileInfo, err := content.Stat()
	if err != nil {
		return nil, err
	}
	if fileInfo.Size() == 0 {
		return menu, nil
	}
	reader := json.NewDecoder(content)

	err1 := reader.Decode(&menu)
	if err1 != nil {
		return nil, err1
	}
	return menu, nil
}
