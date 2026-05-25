package controller

import (
	"encoding/json"
	"net/http"
	"rqms/model"
	"rqms/utils/httpReps"
)

// StaffRegister handles staff registration
func StaffRegister(w http.ResponseWriter, r *http.Request) {
	var s model.Staff
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		httpReps.RespondWithError(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}
	defer r.Body.Close()

	if s.Role == "" {
		s.Role = "staff"
	}

	if err := s.RegisterStaff(); err != nil {
		httpReps.RespondWithError(w, http.StatusInternalServerError, "Email already registered")
		return
	}
	httpReps.RespondWithJSON(w, http.StatusCreated, map[string]string{"status": "registered"})
}

// StaffLogin handles staff login
func StaffLogin(w http.ResponseWriter, r *http.Request) {
	var s model.Staff
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		httpReps.RespondWithError(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}
	defer r.Body.Close()

	role, err := s.AuthenticateStaff()
	if err != nil {
		httpReps.RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	httpReps.RespondWithJSON(w, http.StatusOK, map[string]string{
		"status":  "success",
		"message": "Welcome, " + s.Email,
		"role":    role,
		"email":   s.Email,
	})
}

// GetStaffList returns all staff
func GetStaffList(w http.ResponseWriter, r *http.Request) {
	staffList, err := model.GetAllStaff()
	if err != nil {
		httpReps.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	httpReps.RespondWithJSON(w, http.StatusOK, staffList)
}
