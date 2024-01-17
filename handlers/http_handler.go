package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/kenny-mwendwa/go-restapi-crud/db"
	"github.com/kenny-mwendwa/go-restapi-crud/models"
)

// CREATE USER
func HttpCreateUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	db, err := db.ConnectDB()

	if err != nil {
		log.Fatal(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Parse form data
	err = r.ParseForm()
	if err != nil {
		log.Println("Error parsing form data:", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Get form values
	name := r.FormValue("name")
	email := r.FormValue("email")
	ageStr := r.FormValue("age")

	// Convert age to unit
	age, err := strconv.ParseUint(ageStr, 10, 32)
	if err != nil {
		log.Println("Error converting age to unit:", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	newUser := models.User{
		Name:  name,
		Email: email,
		Age:   uint(age),
	}

	// Create user
	result := db.Create(&newUser)
	if result.Error != nil {
		log.Println("Error creating user:", result.Error)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Return a success response
	w.WriteHeader(http.StatusCreated)
}

// GET ALL USERS
func HttpGetUsers(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	db, err := db.ConnectDB()
	if err != nil {
		log.Fatal(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var users []models.User

	// Query the database for all users
	result := db.Find(&users)
	if result.Error != nil {
		log.Println("Error fetching users from the database:", result.Error)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	usersJSON, err := json.Marshal(users)

	if err != nil {
		log.Println("Error marshaling users to JSON:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Set Content-Type header
	w.Header().Set("Content-Type", "application/json")

	// Write JSON response to the client
	w.WriteHeader(http.StatusOK)
	w.Write(usersJSON)
}

// GET ONE USER
func HttpGetUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	db, err := db.ConnectDB()
	if err != nil {
		log.Fatal(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Extract user ID from request URL parameters
	userIDStr := ps.ByName("id")

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		log.Println("Error converting userId to unit32:", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	var user models.User

	// Query the DB for the user with the specified ID
	result := db.First(&user, userID)
	if result.Error != nil {
		log.Println("Error fetching user from the database:", result.Error)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	userJSON, err := json.Marshal(user)

	if err != nil {
		log.Println("Error marshaling user to JSON:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(userJSON)
}

// UPDATE USER
func HttpUpdateUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	db, err := db.ConnectDB()
	if err != nil {
		log.Fatal(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Extract user ID from request URL parameters
	userIDStr := ps.ByName("id")

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		log.Println("Error converting userId to unit32:", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	var existingUser models.User
	result := db.First(&existingUser, userID)
	if result.Error != nil {
		log.Println("Error fetching user from the database:", result.Error)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Get form values
	name := r.FormValue("name")
	email := r.FormValue("email")
	ageStr := r.FormValue("age")

	// Update user fields if provided
	if name != "" {
		existingUser.Name = name
	}

	if email != "" {
		existingUser.Email = email
	}

	if ageStr != "" {
		age, err := strconv.ParseUint(ageStr, 10, 32)
		if err != nil {
			log.Println("Error converting age to unit32:", err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		existingUser.Age = uint(age)
	}

	// Save the updated user to the database
	db.Save(&existingUser)

	w.WriteHeader(http.StatusOK)
}

// DELETE USER
func HttpDeleteUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	db, err := db.ConnectDB()
	if err != nil {
		log.Fatal(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Extract user ID from request URL parameters
	userIDStr := ps.ByName("id")

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		log.Println("Error converting userID to unit32:", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	var existingUser models.User
	result := db.First(&existingUser, userID)
	if result.Error != nil {
		log.Println("Error fetching user from the database:", result.Error)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	db.Delete(&existingUser)

	w.WriteHeader(http.StatusNoContent)
}
