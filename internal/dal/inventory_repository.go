package dal

import (
	"encoding/json"
	"errors"
	"io"
	"os"

	"hot-coffee/models"
)

// InventoryRepository defines the methods for reading and writing inventory data.
type InventoryRepository interface {
	ReadJSONInv() ([]models.InventoryItem, error)   // Reads the inventory data from a JSON file.
	WriteJSONInv(body []models.InventoryItem) error // Writes the updated inventory data to a JSON file.
}

// jsonInvRepository implements the InventoryRepository interface using JSON file storage.
type jsonInvRepository struct{}

// NewJSONInvRepository creates and returns a new instance of jsonInvRepository.
func NewJSONInvRepository() InventoryRepository {
	return &jsonInvRepository{}
}

// ReadJSONInv reads the inventory data from a JSON file and returns it as a slice of InventoryItem objects.
// If there is an error opening or decoding the file, it returns an empty slice and the error.
func (r *jsonInvRepository) ReadJSONInv() ([]models.InventoryItem, error) {
	var newInv []models.InventoryItem

	option := os.O_RDONLY | os.O_RDWR                           // Set file open options to read-only and read-write.
	content, err := os.OpenFile(Inventoryitem(), option, 0o644) // Open the JSON file with read access.
	if err != nil {
		return newInv, err
	}
	defer content.Close()

	reader := json.NewDecoder(content) // Create a JSON decoder for reading the inventory file.
	err1 := reader.Decode(&newInv)     // Decode the JSON data into the newInv slice.
	if err1 != nil {
		return newInv, err1 // Return an empty slice and error if decoding fails.
	}
	return newInv, nil // Return the populated slice of inventory items.
}

// WriteJSONInv writes the updated inventory data to the JSON file.
// It also creates a backup of the inventory file by copying it to ReserveInventory.
// Returns an error if writing to the file or creating the backup fails.
func (r *jsonInvRepository) WriteJSONInv(newInventory []models.InventoryItem) error {
	option := os.O_CREATE | os.O_WRONLY | os.O_TRUNC            // Set options to create, write-only, and truncate the file.
	content, err := os.OpenFile(Inventoryitem(), option, 0o644) // Open or create the inventory file with write access.
	if err != nil {
		return errors.New("")
	}
	defer content.Close()

	encoder := json.NewEncoder(content) // Create a JSON encoder for writing the updated inventory.
	encoder.SetIndent("", "  ")         // Set indentation for formatted JSON output.
	err = encoder.Encode(newInventory)  // Write the new inventory data to the file.
	if err != nil {
		return err // Return error if encoding fails.
	}

	// Create a backup of the updated inventory file.
	newContent, err := os.Open(Inventoryitem()) // Re-open the inventory file for reading.
	if err != nil {
		return err // Return error if re-opening fails.
	}

	reserveInv, err := os.Create(ReserveInventory) // Create the backup file.
	if err != nil {
		return err // Return error if backup creation fails.
	}

	_, err = io.Copy(reserveInv, newContent) // Copy the contents of the inventory file to the backup.
	if err != nil {
		return err // Return error if copying fails.
	}
	return nil // Return nil if writing and backup succeed.
}
