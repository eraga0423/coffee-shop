package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"hot-coffee/internal/service"
	"hot-coffee/models"
)

// InventoryHandler interface defines HTTP handler methods for inventory operations.
type InventoryHandler interface {
	PostInv(w http.ResponseWriter, r *http.Request)     // Handles adding new inventory items.
	GetInv(w http.ResponseWriter, r *http.Request)      // Retrieves all inventory items.
	GetInvID(w http.ResponseWriter, r *http.Request)    // Retrieves a single inventory item by ID.
	PutInvID(w http.ResponseWriter, r *http.Request)    // Updates an inventory item by ID.
	DeleteInvID(w http.ResponseWriter, r *http.Request) // Deletes an inventory item by ID.
}

// InvHandler struct handles requests related to inventory operations.
type InvHandler struct {
	invService service.InventoryService
}

// NewInvHandler creates and returns a new inventory handler with the provided service.
func NewInvHandler(invService service.InventoryService) InventoryHandler {
	return &InvHandler{invService: invService}
}

// PostInv handles the addition of new inventory items by decoding JSON from the request body.
func (h *InvHandler) PostInv(w http.ResponseWriter, r *http.Request) {
	if err := CheckContentType(r); err != nil {
		SendError(w, http.StatusBadRequest, err)
		return
	}
	newInventory := []models.InventoryItem{}
	err := json.NewDecoder(r.Body).Decode(&newInventory)
	if err != nil {
		SendError(w, http.StatusBadRequest, err)
		return
	}
	err = h.invService.ServicePostInv(newInventory)
	if err != nil {
		SendError(w, http.StatusBadRequest, err)
		return
	}

	SendSucces(w, http.StatusCreated, "New inventory item added")
}

// DeleteInvID handles the deletion of an inventory item by its ID, parsed from the URL path.
func (h *InvHandler) DeleteInvID(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	path = strings.Trim(path, "/")
	parts := strings.SplitN(path, "/", 2)
	if len(parts) != 2 {
		err := errors.New("URL length")
		SendError(w, http.StatusBadRequest, err)
		return
	}
	err := h.invService.ServiceInvDelete(parts[1])
	if err != nil {
		SendError(w, http.StatusBadRequest, err)
		return
	}
	SendSucces(w, http.StatusNoContent, "Inventory item deleted")
}

// GetInv retrieves and sends all inventory items as JSON.
func (h *InvHandler) GetInv(w http.ResponseWriter, r *http.Request) {
	content, err := h.invService.ServiceGetInvItem()
	if err != nil {
		SendError(w, http.StatusBadRequest, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(content)
	if err != nil {
		SendError(w, http.StatusInternalServerError, err)
		return
	}
}

// GetInvID retrieves a single inventory item by ID, parsed from the URL path, and sends it as JSON.
func (h *InvHandler) GetInvID(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	path = strings.Trim(path, "/")
	parts := strings.SplitN(path, "/", 2)
	if len(parts) != 2 {
		err := errors.New("URL length")
		SendError(w, http.StatusBadRequest, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	newGetInvID, err := h.invService.ServiceGetInvID(parts[1])
	if err != nil {
		SendError(w, http.StatusBadRequest, err)
		return
	}
	err = json.NewEncoder(w).Encode(newGetInvID)
	if err != nil {
		SendError(w, http.StatusInternalServerError, err)
		return
	}
}

// PutInvID updates an existing inventory item by its ID, parsed from the URL path, using data from the request body.
func (h *InvHandler) PutInvID(w http.ResponseWriter, r *http.Request) {
	if err := CheckContentType(r); err != nil {
		SendError(w, http.StatusBadRequest, err)
		return
	}
	var newEdit models.InventoryItem
	file, err := io.ReadAll(r.Body)
	if err != nil {
		SendError(w, http.StatusBadRequest, err)
		return
	}
	err = json.Unmarshal(file, &newEdit)
	if err != nil {
		SendError(w, http.StatusBadRequest, err)
		return
	}
	path := r.URL.Path
	path = strings.Trim(path, "/")
	parts := strings.SplitN(path, "/", 2)
	if len(parts) != 2 {
		err := errors.New("URL length")
		SendError(w, http.StatusBadRequest, err)
		return
	}
	err = h.invService.ServicePutInvID(parts[1], newEdit)
	if err != nil {
		SendError(w, http.StatusBadRequest, err)
		return
	}
	SendSucces(w, http.StatusOK, "Inventory item updated")
}
