package router

import (
	"log"
	"net/http"
	"rqms/controller"

	"github.com/gorilla/mux"
)

func Router() {
	router := mux.NewRouter()

	// ---- CUSTOMER AUTH ----
	router.HandleFunc("/signin", controller.Adduser).Methods("POST")
	router.HandleFunc("/login", controller.Login).Methods("POST")

	// ---- STAFF AUTH ----
	router.HandleFunc("/staff/register", controller.StaffRegister).Methods("POST")
	router.HandleFunc("/staff/login", controller.StaffLogin).Methods("POST")
	router.HandleFunc("/staff/list", controller.GetStaffList).Methods("GET")

	// ---- QUEUE ----
	router.HandleFunc("/user", controller.AddUserData).Methods("POST")
	router.HandleFunc("/users", controller.GetAllUsers).Methods("GET")
	router.HandleFunc("/user", controller.DeleteUser).Methods("DELETE")
	router.HandleFunc("/queue/status", controller.UpdateStatus).Methods("PUT")
	router.HandleFunc("/queue/stats", controller.GetStats).Methods("GET")

	// ---- TABLES ----
	router.HandleFunc("/tables", controller.GetTables).Methods("GET")
	router.HandleFunc("/tables/status", controller.UpdateTableStatus).Methods("PUT")
	router.HandleFunc("/tables/add", controller.AddTable).Methods("POST")
	router.HandleFunc("/tables/delete", controller.DeleteTable).Methods("DELETE")

	// ---- STATIC FILES ----
	fs := http.FileServer(http.Dir("./view"))
	router.PathPrefix("/").Handler(fs)

	log.Println("RQMS running on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
