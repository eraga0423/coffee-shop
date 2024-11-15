package service

import (
	"errors"
	"fmt"
	"strings"

	"hot-coffee/internal/dal"
	"hot-coffee/models"
)

type MenuService interface {
	ServiceGetMenuItem() ([]models.MenuItem, error)
	ServicePostMenu(content []models.MenuItem) error
	ServiceGetMenuID(id string) (models.MenuItem, error)
	ServicePutMenuID(id string, newEdit models.MenuItem) error
	EditStructureMenu(EditableStructure models.MenuItem, newEdit models.MenuItem) (models.MenuItem, error)
	ServiceDelete(id string) error
	checkIngredients(IngredientId string) error
}

type menuService struct {
	menuRepo dal.MenuRepository
}

// Initializes and returns a new instance of menuService with the provided repository
func NewMenuService(menuRepo dal.MenuRepository) MenuService {
	return &menuService{menuRepo: menuRepo}
}

// Adds new menu items to the menu, checking for duplicates and validating data
func (s *menuService) ServicePostMenu(content []models.MenuItem) error {
	result := []models.MenuItem{}
	menuItems, err := s.menuRepo.ReadJSONMenu()
	result = menuItems
	if err != nil {
		return err
	}
	for _, oneMenuItem := range content {
		for _, menuItemIngr := range oneMenuItem.Ingredients {
			if err := s.checkIngredients(menuItemIngr.IngredientID); err != nil {
				return err
			}
		}
		if check, err := s.CheckMenu(oneMenuItem); !check {
			return err
		}

		if !checkForClone(oneMenuItem, menuItems) {
			return errors.New("Such ID already exists")
		}
		result = append(result, oneMenuItem)

	}

	err = s.menuRepo.WriteJSONMenu(result)
	if err != nil {
		return err
	}
	return nil
}

// Retrieves all menu items from the repository
func (s *menuService) ServiceGetMenuItem() ([]models.MenuItem, error) {
	return s.menuRepo.ReadJSONMenu()
}

// Retrieves a specific menu item by ID, returning an error if not found
func (s *menuService) ServiceGetMenuID(id string) (models.MenuItem, error) {
	checker := false
	newGetMenuID := models.MenuItem{}
	jsonfilemenu, err := s.menuRepo.ReadJSONMenu()
	if err != nil {
		return newGetMenuID, err
	}
	for _, value := range jsonfilemenu {
		if value.ID == id {
			checker = true
			newGetMenuID = value

		}
	}
	if !checker {
		return newGetMenuID, errors.New("ID not found")
	}
	return newGetMenuID, nil
}

// Updates a specific menu item by ID with new data provided, validating changes
func (s *menuService) ServicePutMenuID(id string, newEdit models.MenuItem) error {
	checker := false
	jsonfilemenu, err := s.menuRepo.ReadJSONMenu()
	if err != nil {
		return err
	}
	for i, oneStructure := range jsonfilemenu {
		if oneStructure.ID == id {
			checker = true
			newEditedStructure, err := s.EditStructureMenu(oneStructure, newEdit)
			if err != nil {
				return err
			}
			jsonfilemenu[i] = newEditedStructure
		}
	}
	err = s.menuRepo.WriteJSONMenu(jsonfilemenu)
	if err != nil {
		return err
	}
	if !checker {
		return errors.New("ID not found")
	}
	return nil
}

// Modifies a menu item structure with new fields where specified, validating the updated item
func (s *menuService) EditStructureMenu(EditableStructure models.MenuItem, newEdit models.MenuItem) (models.MenuItem, error) {
	var newEditedStructure models.MenuItem
	newEditedStructure = EditableStructure
	newEditedStructure.Price = newEdit.Price
	listmenu, err := s.CheckPutMenu(newEdit)
	if err != nil {
		return newEditedStructure, err
	}
	for _, oneObject := range listmenu {
		switch oneObject {
		case "name":
			newEditedStructure.Name = newEdit.Name
		case "description":
			newEditedStructure.Description = newEdit.Description

		case "ingredients":
			newEditedStructure.Ingredients = newEdit.Ingredients
		}
	}
	if check, err := s.CheckMenu(newEditedStructure); !check {
		return newEditedStructure, err
	}
	return newEditedStructure, nil
}

// Checks which fields in the new menu item have non-empty values and returns them
func (s *menuService) CheckPutMenu(newmenu models.MenuItem) ([]string, error) {
	listMenu := []string{}
	newmenuName := strings.Trim(newmenu.Name, " ")
	if newmenuName != "" {
		listMenu = append(listMenu, "name")
	}
	newmenuDescription := strings.Trim(newmenu.Description, " ")
	if newmenuDescription != "" {
		listMenu = append(listMenu, "description")
	}

	for _, msq := range newmenu.Ingredients {
		if err := s.checkIngredients(msq.IngredientID); err != nil {
			return nil, err
		}
		msqIngredients := strings.Trim(msq.IngredientID, " ")
		if msqIngredients != "" {
			listMenu = append(listMenu, "ingredients")
		}
	}
	return listMenu, nil
}

// Deletes a menu item by ID, returning an error if the ID is not found
func (s *menuService) ServiceDelete(id string) error {
	check := false
	newMenu, err := s.menuRepo.ReadJSONMenu()
	if err != nil {
		return err
	}
	index := 0
	for i, oneObjectMenu := range newMenu {
		if oneObjectMenu.ID == id {
			check = true
			index = i
		}
	}
	newMenu = append(newMenu[:index], newMenu[index+1:]...)
	if !check {
		return errors.New("Such ID not found")
	}
	err = s.menuRepo.WriteJSONMenu(newMenu)
	if err != nil {
		return err
	}
	return nil
}

// Checks if a menu item with the same ID already exists in the menu
func checkForClone(newjson models.MenuItem, JSONFile []models.MenuItem) bool {
	for _, value := range JSONFile {
		if value.ID == newjson.ID {
			return false
		}
	}
	return true
}

// Validates the fields of a new menu item to ensure all required fields are filled correctly
func (s *menuService) CheckMenu(newmenu models.MenuItem) (bool, error) {
	newmenuID := strings.Trim(newmenu.ID, " ")
	if newmenuID == "" {
		return false, errors.New("Missing ID")
	}
	newmenuName := strings.Trim(newmenu.Name, " ")
	if newmenuName == "" {
		return false, errors.New("Missing name")
	}
	newmenuDescription := strings.Trim(newmenu.Description, " ")
	if newmenuDescription == "" {
		return false, errors.New("Missing Description")
	}
	if newmenu.Price < 0.0 {
		return false, errors.New("Price cannot be negative")
	}
	for _, msq := range newmenu.Ingredients {
		if msq.Quantity <= 0 {
			return false, errors.New("Ingredients quantity cannot be 0 or negative")
		}
		msqIngredients := strings.Trim(msq.IngredientID, " ")
		if msqIngredients == "" {
			return false, errors.New("Missing ingredients ID")
		}
	}

	return true, nil
}

// Checks if the given ingredient ID exists in the inventory, returning an error if it is missing
func (s *menuService) checkIngredients(IngredientId string) error {
	fileInv, err := s.menuRepo.ReadJSONInventory()
	if err != nil {
		return err
	}
	IngredientMap := make(map[string]int)

	missList := ""
	for _, oneInvItem := range fileInv {
		IngredientMap[oneInvItem.IngredientID] = int(oneInvItem.Quantity)
	}
	_, exists := IngredientMap[IngredientId]
	if !exists {
		missList = IngredientId
		return errors.New(fmt.Sprintf("This ingredient is not in the inventory: %s", missList))
	}
	return nil
}
