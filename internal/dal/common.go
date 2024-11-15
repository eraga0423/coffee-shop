package dal

import (
	"fmt"
	"os"
	"time"
)

var Directory string

const (
	ReserveInventory = "../reserve_copy/inventory.json"
	ReserveMenu      = "../reserve_copy/menu_items.json"

	InventoryitemFile = "inventory.json"
	MenuItemFile      = "menu_items.json"
	OrdersFile        = "orders.json"
)

// Sets the global directory path
func NewDirectory(dir string) {
	Directory = dir
}

// Generates the filename for the reserve orders file with the current date
func ReserveOrder() string {
	nowTime := time.Now()
	timeString := nowTime.Format("2006-01-02")
	return fmt.Sprintf("../reserve_copy/%s_orders.json", timeString)
}

// Returns the full path to the orders file in the specified directory
func Orders() string {
	return fmt.Sprintf("../%s/orders.json", Directory)
}

// Returns the full path to the inventory file in the specified directory
func Inventoryitem() string {
	return fmt.Sprintf("../%s/inventory.json", Directory)
}

// Returns the full path to the menu items file in the specified directory
func Menuitems() string {
	return fmt.Sprintf("../%s/menu_items.json", Directory)
}

// Checks if a file exists at the specified path and returns true if it does
func FileExistsInDirectory(path string) (bool, error) {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	check := info.IsDir()
	return !check, nil
}
