package goControllers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"goCrud/goHandlers"
	"net/http"
)

func SetupUserRoutes(router *mux.Router) {
	router.HandleFunc("/PostData", PostData).Methods("POST")
	router.HandleFunc("/PutData", PutData).Methods("PUT")
	router.HandleFunc("/GetData", GetData).Methods("GET")
	router.HandleFunc("/DeleteData", DeleteData).Methods("DELETE")
}

func PostData(w http.ResponseWriter, r *http.Request) {
	var userToPost goHandlers.User
	err := json.NewDecoder(r.Body).Decode(&userToPost)

	if err != nil {
		http.Error(w, "Invalid user data", http.StatusBadRequest)
		return
	}

	err = goHandlers.PostDataHandler(userToPost)

	if err != nil {
		http.Error(w, "Issue creating user data", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func PutData(w http.ResponseWriter, r *http.Request) {
	fmt.Println("PutData handler called")

	var userToUpdate goHandlers.User
	err := json.NewDecoder(r.Body).Decode(&userToUpdate)

	if err != nil {
		fmt.Println("Error decoding request body:", err)
		http.Error(w, "Invalid user data", http.StatusBadRequest)
		return
	}

	fmt.Printf("Received user data for update: %+v\n", userToUpdate)

	err = goHandlers.PutDataHandler(userToUpdate)

	if err != nil {
		fmt.Println("Error updating user data:", err)
		http.Error(w, "Issue updating user data", http.StatusInternalServerError)
		return
	}

	response := map[string]string{"message": "User updated successfully"}
	jsonResponse, _ := json.Marshal(response)

	fmt.Println("User update successful!")

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func GetData(w http.ResponseWriter, r *http.Request) {
	users, err := goHandlers.GetDataHandler()

	if err != nil {
		http.Error(w, "Error loading users", http.StatusInternalServerError)
		return
	}

	if err = json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, "Error encoding users to JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func DeleteData(w http.ResponseWriter, r *http.Request) {
	var userToDelete goHandlers.User

	err := json.NewDecoder(r.Body).Decode(&userToDelete)

	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = goHandlers.DeleteDataHandler(userToDelete)
	if err != nil {
		http.Error(w, "Error deleting user", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(`{"message": "User deleted successfully"}`))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
