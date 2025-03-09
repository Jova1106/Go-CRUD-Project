package goHandlers

import (
	"JovaCentral/database"
	"JovaCentral/mathLib"
	"encoding/json"
	"fmt"
	"os"
)

type User = database.User

// Loop through the sequential data array and
// add all users that do not match the ID.
func erase(users []User, userToDelete User) []User {
	newUserList := []User{}

	for _, user := range users {
		if user.Id != userToDelete.Id {
			newUserList = append(newUserList, user)
		}
	}

	return newUserList
}

// Loop through the sequential data array
// and update the user if a matching ID is found
func update(users []User, userToPut User) []User {
	for i, user := range users {
		if user.Id == userToPut.Id {
			users[i] = userToPut
			return users
		}
	}

	users = append(users, userToPut)

	return users
}

// Reads the JSON database files and gives the new User instance a unique ID.
// Adds the User to the user list and writes the data back to the JSON database files.
func PostDataHandler(userToPost User) error {
	currentFileContent, _ := os.ReadFile(database.DATA_FILE)
	currentFileContentSequential, _ := os.ReadFile(database.DATA_FILE_SEQUENTIAL)
	var newUserMap map[string]User
	var newUsers []User
	var newUserId = mathLib.RandomString()
	userToPost.Id = newUserId

	if err := json.Unmarshal(currentFileContent, &newUserMap); err != nil {
		return err
	}

	newUserMap[newUserId] = userToPost
	fileContent, err := json.MarshalIndent(newUserMap, "", "  ")

	if err != nil {
		return err
	}

	os.Remove(database.DATA_FILE)

	if err = os.WriteFile(database.DATA_FILE, fileContent, 0644); err != nil {
		return err
	}

	if err := json.Unmarshal(currentFileContentSequential, &newUsers); err != nil {
		return err
	}

	newUsers = append(newUsers, userToPost)
	fileContentSequential, err := json.MarshalIndent(newUsers, "", "  ")

	if err != nil {
		return err
	}

	os.Remove(database.DATA_FILE_SEQUENTIAL)

	if err = os.WriteFile(database.DATA_FILE_SEQUENTIAL, fileContentSequential, 0644); err != nil {
		return err
	}

	return nil
}

// Reads the JSON database files and updates the existing User instance.
// Writes the data back to the JSON database files.
func PutDataHandler(userToPut User) error {
	currentFileContent, err := os.ReadFile(database.DATA_FILE)

	if err != nil {
		return fmt.Errorf("failed to read data file: %w", err)
	}

	currentFileContentSequential, err := os.ReadFile(database.DATA_FILE_SEQUENTIAL)

	if err != nil {
		return fmt.Errorf("failed to read sequential data file: %w", err)
	}

	var newUserMap map[string]User
	var newUsers []User

	if err := json.Unmarshal(currentFileContent, &newUserMap); err != nil {
		return fmt.Errorf("failed to parse JSON data file: %w", err)
	}

	if err := json.Unmarshal(currentFileContentSequential, &newUsers); err != nil {
		return fmt.Errorf("failed to parse JSON sequential file: %w", err)
	}

	if userToPut.Id == "" {
		return fmt.Errorf("user ID is missing")
	}

	newUserMap[userToPut.Id] = userToPut
	fileContent, err := json.MarshalIndent(newUserMap, "", "  ")

	if err != nil {
		return fmt.Errorf("failed to serialize updated user data: %w", err)
	}

	if err := os.WriteFile(database.DATA_FILE, fileContent, 0644); err != nil {
		return fmt.Errorf("failed to write updated user data: %w", err)
	}

	newUsers = update(newUsers, userToPut)
	fileContentSequential, err := json.MarshalIndent(newUsers, "", "  ")

	if err != nil {
		return fmt.Errorf("failed to serialize updated sequential data: %w", err)
	}

	if err := os.WriteFile(database.DATA_FILE_SEQUENTIAL, fileContentSequential, 0644); err != nil {
		return fmt.Errorf("failed to write updated sequential user data: %w", err)
	}

	fmt.Println("User successfully updated:", userToPut)

	return nil
}

// Returns the sequential JSON database file.
// The sequential database file is used to display
// the users in the order that they were created.
func GetDataHandler() ([]User, error) {
	fileContent, err := os.ReadFile(database.DATA_FILE_SEQUENTIAL)

	if err != nil {
		return nil, err
	}

	var users []User

	if err := json.Unmarshal(fileContent, &users); err != nil {
		return nil, err
	}

	return users, nil
}

// Reads the JSON database files and deletes the User instance.
// Writes the data back to the JSON database files.
func DeleteDataHandler(userToDelete User) error {
	currentFileContent, err := os.ReadFile(database.DATA_FILE)

	if err != nil {
		return fmt.Errorf("failed to read %s: %v", database.DATA_FILE, err)
	}

	currentFileContentSequential, err := os.ReadFile(database.DATA_FILE_SEQUENTIAL)

	if err != nil {
		return fmt.Errorf("failed to read %s: %v", database.DATA_FILE_SEQUENTIAL, err)
	}

	var newUserMap map[string]User
	var newUsers []User

	if err := json.Unmarshal(currentFileContent, &newUserMap); err != nil {
		return fmt.Errorf("failed to unmarshal %s: %v", database.DATA_FILE, err)
	}

	delete(newUserMap, userToDelete.Id)
	fileContent, err := json.MarshalIndent(newUserMap, "", "  ")

	if err != nil {
		return fmt.Errorf("failed to marshal updated user map: %v", err)
	}

	if err := os.WriteFile(database.DATA_FILE, fileContent, 0644); err != nil {
		return fmt.Errorf("failed to write to %s: %v", database.DATA_FILE, err)
	}

	if err := json.Unmarshal(currentFileContentSequential, &newUsers); err != nil {
		return fmt.Errorf("failed to unmarshal %s: %v", database.DATA_FILE_SEQUENTIAL, err)
	}

	newUsers = erase(newUsers, userToDelete)
	fileContentSequential, err := json.MarshalIndent(newUsers, "", "  ")

	if err != nil {
		return fmt.Errorf("failed to marshal updated sequential user list: %v", err)
	}

	if err := os.WriteFile(database.DATA_FILE_SEQUENTIAL, fileContentSequential, 0644); err != nil {
		return fmt.Errorf("failed to write to %s: %v", database.DATA_FILE_SEQUENTIAL, err)
	}

	return nil
}
