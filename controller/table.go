package controller

import (
	"encoding/json"
	"net/http"
	"rqms/model"
	"rqms/utils/httpReps"
)

// GetTables returns all tables
func GetTables(w http.ResponseWriter, r *http.Request) {
	tables, err := model.GetTables()
	if err != nil {
		httpReps.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if tables == nil {
		tables = []model.Table{}
	}
	httpReps.RespondWithJSON(w, http.StatusOK, tables)
}

// UpdateTableStatus changes a table status
func UpdateTableStatus(w http.ResponseWriter, r *http.Request) {
	var body struct {
		TableNo string `json:"table_no"`
		Status  string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		httpReps.RespondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}
	if err := model.UpdateTableStatus(body.TableNo, body.Status); err != nil {
		httpReps.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	httpReps.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "table updated"})
}

// AddTable inserts a new table
func AddTable(w http.ResponseWriter, r *http.Request) {
	var body struct {
		TableNo  string `json:"table_no"`
		Capacity int    `json:"capacity"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		httpReps.RespondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}
	if body.TableNo == "" {
		httpReps.RespondWithError(w, http.StatusBadRequest, "table_no required")
		return
	}
	if body.Capacity <= 0 {
		body.Capacity = 4
	}
	if err := model.AddTable(body.TableNo, body.Capacity); err != nil {
		httpReps.RespondWithError(w, http.StatusInternalServerError, "Table number already exists")
		return
	}
	httpReps.RespondWithJSON(w, http.StatusCreated, map[string]string{"message": "table added"})
}

// DeleteTable removes a table
func DeleteTable(w http.ResponseWriter, r *http.Request) {
	var body struct {
		TableNo string `json:"table_no"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		httpReps.RespondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}
	if err := model.DeleteTable(body.TableNo); err != nil {
		httpReps.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	httpReps.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "table deleted"})
}
