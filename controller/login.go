package controller

import (
	"encoding/json"
	"net/http"
	"rqms/model"
	"rqms/utils/httpReps"
)

// Login handles customer login
func Login(w http.ResponseWriter, r *http.Request) {
	var l model.Login
	if err := json.NewDecoder(r.Body).Decode(&l); err != nil {
		httpReps.RespondWithError(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}
	defer r.Body.Close()

	user, err := l.Authenticate()
	if err != nil {
		httpReps.RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	httpReps.RespondWithJSON(w, http.StatusOK, map[string]string{
		"status":     "success",
		"message":    "Welcome back!",
		"role":       "customer",
		"email":      user.Email,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
	})
}
