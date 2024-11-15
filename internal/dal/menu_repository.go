package dal

import (
	"encoding/json"
	"errors"
	"io"
	"os"

	"hot-coffee/models"
)

type MenuRepository interface {
	ReadJSONMenu() ([]models.MenuItem, error)
	WriteJSONMenu(newMenuItem []models.MenuItem) error
	ReadJSONInventory() ([]models.InventoryItem, error)
}
type jsonMenuRepository struct{}

// Creates and returns a new instance of jsonMenuRepository
func NewJSONMenuRepository() MenuRepository {
	return &jsonMenuRepository{}
}

// Reads and decodes menu items data from the JSON file, returning a slice of menu items
func (r *jsonMenuRepository) ReadJSONMenu() ([]models.MenuItem, error) {
	var newMenu []models.MenuItem

	option := os.O_RDONLY | os.O_RDWR
	content, err := os.OpenFile(Menuitems(), option, 0o644)
	if err != nil {
		return newMenu, err
	}
	defer content.Close()

	reader := json.NewDecoder(content)
	err1 := reader.Decode(&newMenu)
	if err1 != nil {
		return newMenu, err1
	}
	return newMenu, nil
}

// Writes the provided menu items to the JSON file, creating or truncating as necessary
func (r *jsonMenuRepository) WriteJSONMenu(newMenuItem []models.MenuItem) error {
	option := os.O_CREATE | os.O_WRONLY | os.O_TRUNC
	content, err := os.OpenFile(Menuitems(), option, 0o644)
	if err != nil {
		return errors.New("")
	}
	defer content.Close()
	encoder := json.NewEncoder(content)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(newMenuItem)
	if err != nil {
		return err
	}
	newContent, err := os.Open(Menuitems())
	if err != nil {
		return err
	}
	reserveMenu, err := os.Create(ReserveMenu)
	if err != nil {
		return err
	}

	_, err = io.Copy(reserveMenu, newContent)
	if err != nil {
		return err
	}
	return nil
}

// Reads and decodes inventory items data from the JSON file, returning a slice of inventory items
func (r *jsonMenuRepository) ReadJSONInventory() ([]models.InventoryItem, error) {
	var newInv []models.InventoryItem

	option := os.O_RDONLY | os.O_RDWR
	content, err := os.OpenFile(Inventoryitem(), option, 0o644)
	if err != nil {
		return newInv, err
	}
	defer content.Close()

	reader := json.NewDecoder(content)
	err1 := reader.Decode(&newInv)
	if err1 != nil {
		return nil, err1
	}
	return newInv, nil
}
