package controller

import (
	"encoding/json"
	"net/http"
	"rqms/model"
	"rqms/utils/httpReps"

	"strconv"
)

// AddUserData creates a new queue entry
func AddUserData(w http.ResponseWriter, r *http.Request) {
	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		httpReps.RespondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}
	if err := user.CreateUser(); err != nil {
		httpReps.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	httpReps.RespondWithJSON(w, http.StatusCreated, map[string]string{"message": "added to queue"})
}

// GetAllUsers returns all queue entries
func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := model.GetUsers()
	if err != nil {
		httpReps.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if users == nil {
		users = []model.User{}
	}
	httpReps.RespondWithJSON(w, http.StatusOK, users)
}

// DeleteUser removes a queue entry
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		httpReps.RespondWithError(w, http.StatusBadRequest, "Invalid ID")
		return
	}
	if err := model.DeleteUser(id); err != nil {
		httpReps.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	httpReps.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "deleted"})
}

// UpdateStatus changes queue entry status
func UpdateStatus(w http.ResponseWriter, r *http.Request) {
	var body struct {
		ID      int    `json:"id"`
		Status  string `json:"status"`
		TableNo string `json:"table_no"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		httpReps.RespondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}
	if err := model.UpdateQueueStatus(body.ID, body.Status, body.TableNo); err != nil {
		httpReps.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	httpReps.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "status updated"})
}

// GetStats returns queue statistics
func GetStats(w http.ResponseWriter, r *http.Request) {
	stats, err := model.GetQueueStats()
	if err != nil {
		httpReps.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	httpReps.RespondWithJSON(w, http.StatusOK, stats)
}
