package gateway

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"gorm.io/gorm"

	"github.com/rahmanazhar/FoodSupplyChain/pkg/auth"
	"github.com/rahmanazhar/FoodSupplyChain/pkg/models"
)

// validRoles is the set of roles an admin may assign to a user.
var validRoles = map[string]struct{}{
	auth.RoleAdmin:    {},
	auth.RoleManager:  {},
	auth.RoleOperator: {},
	auth.RoleViewer:   {},
}

// handleListUsers returns all users (admin only). The password hash is never
// serialised thanks to the json:"-" tag on the model.
func (a *Auth) handleListUsers(w http.ResponseWriter, r *http.Request) {
	var users []models.User
	if err := a.db.Order("created_at asc").Find(&users).Error; err != nil {
		writeJSON(w, http.StatusInternalServerError, errBody(err.Error()))
		return
	}
	writeJSON(w, http.StatusOK, users)
}

// handleUpdateUserRole changes a user's role (admin only). It validates the
// requested role and 404s when the user does not exist.
func (a *Auth) handleUpdateUserRole(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var body struct {
		Role string `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, errBody("invalid request body"))
		return
	}
	if _, ok := validRoles[body.Role]; !ok {
		writeJSON(w, http.StatusBadRequest, errBody("role must be one of admin, manager, operator, viewer"))
		return
	}

	var user models.User
	if err := a.db.First(&user, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			writeJSON(w, http.StatusNotFound, errBody("user not found"))
			return
		}
		writeJSON(w, http.StatusInternalServerError, errBody(err.Error()))
		return
	}

	user.Role = body.Role
	user.UpdatedAt = time.Now()
	if err := a.db.Save(&user).Error; err != nil {
		writeJSON(w, http.StatusInternalServerError, errBody(err.Error()))
		return
	}
	writeJSON(w, http.StatusOK, user)
}
