package goControllers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"goCrud/goHandlers"
	"net/http"
)

// Define constant route variables
const (
	PostDataRoute   = "/PostData"
	PutDataRoute    = "/PutData"
	GetDataRoute    = "/GetData"
	DeleteDataRoute = "/DeleteData"
)

// Data structure for containing the route handler function and the request method
type Route struct {
	handlerFunction http.HandlerFunc
	method          string
}

// Define routes to be setup
var routeMap = map[string]Route{
	PostDataRoute:   {PostData, "POST"},
	PutDataRoute:    {PutData, "PUT"},
	GetDataRoute:    {GetData, "GET"},
	DeleteDataRoute: {DeleteData, "DELETE"},
}

func SetupUserRoutes(router *mux.Router) {
	// Loop through the routes map and register the routes
	for path, route := range routeMap {
		router.HandleFunc(path, route.handlerFunction).Methods(route.method)
	}
}

// Decodes a new JSON-encoded User instance and
// calls goHandlers.PostDataHandler to insert
// the new User into the JSON database.
func PostData(w http.ResponseWriter, r *http.Request) {
	var userToPost goHandlers.User

	if err := json.NewDecoder(r.Body).Decode(&userToPost); err != nil {
		http.Error(w, "Invalid user data", http.StatusBadRequest)
		return
	}

	if err := goHandlers.PostDataHandler(userToPost); err != nil {
		http.Error(w, "Issue creating user data", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

// Decodes an existing JSON-encoded User instance and
// calls goHandlers.PutDataHandler to update
// the User in the JSON database.
func PutData(w http.ResponseWriter, r *http.Request) {
	fmt.Println("PutData handler called")

	var userToUpdate goHandlers.User

	if err := json.NewDecoder(r.Body).Decode(&userToUpdate); err != nil {
		fmt.Println("Error decoding request body:", err)
		http.Error(w, "Invalid user data", http.StatusBadRequest)
		return
	}

	fmt.Printf("Received user data for update: %+v\n", userToUpdate)

	if err := goHandlers.PutDataHandler(userToUpdate); err != nil {
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

// Calls goHandlers.GetDataHandler to get the data
// from the JSON database, and sends it to the client
func GetData(w http.ResponseWriter, r *http.Request) {
	users, err := goHandlers.GetDataHandler()

	if err != nil {
		http.Error(w, "Error loading users", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err = json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, "Error encoding users to JSON", http.StatusInternalServerError)
		return
	}
}

// Decodes an existing JSON-encoded User instance and
// calls goHandlers.DeleteDataHandler to delete
// the User from the JSON database.
func DeleteData(w http.ResponseWriter, r *http.Request) {
	var userToDelete goHandlers.User

	if err := json.NewDecoder(r.Body).Decode(&userToDelete); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := goHandlers.DeleteDataHandler(userToDelete); err != nil {
		http.Error(w, "Error deleting user", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(`{"message": "User deleted successfully"}`))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
