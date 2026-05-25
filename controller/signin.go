package controller

import (
	"encoding/json"
	"net/http"
	"rqms/model"
	"rqms/utils/httpReps"
)

// Adduser handles customer registration
func Adduser(w http.ResponseWriter, r *http.Request) {
	var s model.Sigin
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		httpReps.RespondWithError(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}
	defer r.Body.Close()

	if err := s.Adduser(); err != nil {
		httpReps.RespondWithError(w, http.StatusInternalServerError, "Email already registered")
		return
	}
	httpReps.RespondWithJSON(w, http.StatusCreated, map[string]string{"status": "created"})
}
