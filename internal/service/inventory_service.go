package service

import (
	"errors"
	"strings"

	"hot-coffee/internal/dal"
	"hot-coffee/models"
)

// InventoryService defines methods for handling inventory operations.
type InventoryService interface {
	ServiceGetInvItem() ([]models.InventoryItem, error)                                                                  // Retrieves all inventory items.
	ServicePostInv(content []models.InventoryItem) error                                                                 // Adds new inventory items.
	ServiceGetInvID(id string) (models.InventoryItem, error)                                                             // Retrieves a single inventory item by ID.
	ServicePutInvID(id string, newEdit models.InventoryItem) error                                                       // Updates an existing inventory item by ID.
	EditInvStructure(EditableStructure models.InventoryItem, newEdit models.InventoryItem) (models.InventoryItem, error) // Edits specific fields of an inventory item.
	ServiceInvDelete(id string) error                                                                                    // Deletes an inventory item by ID.
}

// invService implements the InventoryService interface using InventoryRepository.
type invService struct {
	invRepo dal.InventoryRepository
}

// NewInvService creates and returns a new instance of invService.
func NewInvService(invRepo dal.InventoryRepository) InventoryService {
	return &invService{invRepo: invRepo}
}

// ServicePostInv adds new inventory items to the inventory if they pass validation and don't already exist.
func (s *invService) ServicePostInv(content []models.InventoryItem) error {
	result := []models.InventoryItem{}
	invItems, err := s.invRepo.ReadJSONInv() // Reads existing inventory items.
	result = invItems
	if err != nil {
		return err
	}
	for _, oneInvItem := range content {
		if check, err := s.CheckInvPost(oneInvItem); !check {
			return err // Return error if the new item fails validation.
		}

		if !CheckIsNew(oneInvItem, invItems) {
			return errors.New("Such ID already exists") // Ensure each new item has a unique ID.
		}
		result = append(result, oneInvItem) // Append new valid items to the result list.
	}

	err = s.invRepo.WriteJSONInv(result) // Save the updated inventory.
	if err != nil {
		return err
	}
	return nil
}

// ServiceGetInvItem retrieves all inventory items from storage.
func (s *invService) ServiceGetInvItem() ([]models.InventoryItem, error) {
	return s.invRepo.ReadJSONInv()
}

// ServiceGetInvID retrieves a specific inventory item by its ID.
// Returns an error if the ID is not found.
func (s *invService) ServiceGetInvID(id string) (models.InventoryItem, error) {
	checker := false
	newGetInvID := models.InventoryItem{}
	jsonfileinv, err := s.invRepo.ReadJSONInv()
	if err != nil {
		return newGetInvID, err
	}
	for _, value := range jsonfileinv {
		if value.IngredientID == id {
			checker = true
			newGetInvID = value
		}
	}
	if !checker {
		return newGetInvID, errors.New("ID not found")
	}
	return newGetInvID, nil
}

// ServicePutInvID updates an existing inventory item identified by ID with new data.
func (s *invService) ServicePutInvID(id string, newEdit models.InventoryItem) error {
	checker := false
	jsonfileinv, err := s.invRepo.ReadJSONInv()
	if err != nil {
		return err
	}
	for i, oneStructure := range jsonfileinv {
		if oneStructure.IngredientID == id {
			checker = true
			newEditedStructure, err := s.EditInvStructure(oneStructure, newEdit)
			if err != nil {
				return err
			}
			jsonfileinv[i] = newEditedStructure
		}
	}
	err = s.invRepo.WriteJSONInv(jsonfileinv)
	if err != nil {
		return err
	}
	if !checker {
		return errors.New("ID not found")
	}
	return nil
}

// EditInvStructure applies non-empty fields from newEdit to EditableStructure.
// Returns the edited inventory item if validation succeeds.
func (s *invService) EditInvStructure(EditableStructure models.InventoryItem, newEdit models.InventoryItem) (models.InventoryItem, error) {
	var newEditedStructure models.InventoryItem
	newEditedStructure = EditableStructure

	listinv := CheckPutInv(newEdit) // Get list of fields to edit.
	for _, oneObject := range listinv {
		switch oneObject {
		case "name":
			newEditedStructure.Name = newEdit.Name
		case "quantity":
			newEditedStructure.Quantity = newEdit.Quantity
		case "unit":
			newEditedStructure.Unit = newEdit.Unit
		}
	}
	if check, err := s.CheckInvPost(newEditedStructure); !check && err != nil {
		return newEditedStructure, err // Return error if the edited structure is invalid.
	}

	return newEditedStructure, nil
}

// CheckPutInv checks for non-empty fields in newEdit to determine which fields to update.
func CheckPutInv(newinv models.InventoryItem) []string {
	listInventory := []string{}
	newInvName := strings.TrimSpace(newinv.Name)
	if newInvName != "" {
		listInventory = append(listInventory, "name")
	}
	if newinv.Quantity != 0 {
		listInventory = append(listInventory, "quantity")
	}
	newInvUnit := strings.TrimSpace(newinv.Unit)
	if newInvUnit != "" {
		listInventory = append(listInventory, "unit")
	}
	return listInventory
}

// ServiceInvDelete deletes an inventory item by ID.
// Returns an error if the item with the specified ID is not found.
func (s *invService) ServiceInvDelete(id string) error {
	check := false
	newInv, err := s.invRepo.ReadJSONInv()
	if err != nil {
		return err
	}
	index := 0
	for i, oneObjectMenu := range newInv {
		if oneObjectMenu.IngredientID == id {
			check = true
			index = i
		}
	}
	newInv = append(newInv[:index], newInv[index+1:]...) // Remove item by index.
	if !check {
		return errors.New("Item not found with this ID")
	}
	err = s.invRepo.WriteJSONInv(newInv)
	if err != nil {
		return err
	}
	return nil
}

// CheckInvPost validates the fields of a new inventory item.
// Returns false with an error if any required field is missing or invalid.
func (r *invService) CheckInvPost(newinv models.InventoryItem) (bool, error) {
	newInvID := strings.TrimSpace(newinv.IngredientID)
	if newInvID == "" {
		return false, errors.New("Missing Ingredient ID")
	}
	newInvName := strings.TrimSpace(newinv.Name)
	if newInvName == "" {
		return false, errors.New("Missing name")
	}
	if newinv.Quantity < 0 {
		return false, errors.New("Quantity cannot be negative")
	}
	newInvUnit := strings.TrimSpace(newinv.Unit)
	if newInvUnit == "" {
		return false, errors.New("Missing Unit")
	}

	return true, nil
}

// CheckIsNew checks if a given inventory item ID already exists in the current inventory.
// Returns false if the ID exists, true otherwise.
func CheckIsNew(newjson models.InventoryItem, JSONFile []models.InventoryItem) bool {
	for _, value := range JSONFile {
		if value.IngredientID == newjson.IngredientID {
			return false
		}
	}
	return true
}
