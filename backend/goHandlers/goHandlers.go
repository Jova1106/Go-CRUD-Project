package goHandlers

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
)

const DATA_FILE = "../database/data.json"
const DATA_FILE_SEQUENTIAL = "../database/dataSequential.json"

func erase(users []User, userToDelete User) []User {
	newUserList := []User{}

	for _, user := range users {
		if user.Id != userToDelete.Id {
			newUserList = append(newUserList, user)
		}
	}

	return newUserList
}

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

type User struct {
	Id    string `json:id`
	Name  string `json:name`
	Email string `json:email`
}

var alphanumericRanges = [][2]int{
	{48, 57},  // 0 - 9
	{65, 90},  // A - Z
	{97, 122}, // a - z
}

func randomString() string {
	currentFileContent, _ := os.ReadFile(DATA_FILE)
	var newUserMap map[string]User

	if err := json.Unmarshal(currentFileContent, &newUserMap); err != nil {
		return fmt.Sprintf("%v", err)
	}

	var alphanumeric = true
	var length = 64
	var ret string
	var arr []rune

	if alphanumeric {
		for i := 0; i < length; i++ {
			randRange := alphanumericRanges[rand.Intn(3)]
			arr = append(arr, rune(rand.Intn(randRange[1]-randRange[0]+1)+randRange[0]))
		}
	} else {
		for i := 0; i < length; i++ {
			arr = append(arr, rune(rand.Intn(94)+33))
		}
	}

	ret = string(arr)

	if _, exists := newUserMap[ret]; exists {
		return randomString()
	}

	return ret
}

func PostDataHandler(userToPost User) error {
	currentFileContent, _ := os.ReadFile(DATA_FILE)
	currentFileContentSequential, _ := os.ReadFile(DATA_FILE_SEQUENTIAL)
	var newUserMap map[string]User
	var newUsers []User
	var newUserId = randomString()
	userToPost.Id = newUserId

	if err := json.Unmarshal(currentFileContent, &newUserMap); err != nil {
		return err
	}

	newUserMap[newUserId] = userToPost
	fileContent, err := json.MarshalIndent(newUserMap, "", "  ")

	if err != nil {
		return err
	}

	os.Remove(DATA_FILE)

	if err = os.WriteFile(DATA_FILE, fileContent, 0644); err != nil {
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

	os.Remove(DATA_FILE_SEQUENTIAL)

	if err = os.WriteFile(DATA_FILE_SEQUENTIAL, fileContentSequential, 0644); err != nil {
		return err
	}

	return nil
}

func PutDataHandler(userToPut User) error {
	currentFileContent, err := os.ReadFile(DATA_FILE)

	if err != nil {
		return fmt.Errorf("failed to read data file: %w", err)
	}

	currentFileContentSequential, err := os.ReadFile(DATA_FILE_SEQUENTIAL)

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

	if err := os.WriteFile(DATA_FILE, fileContent, 0644); err != nil {
		return fmt.Errorf("failed to write updated user data: %w", err)
	}

	newUsers = update(newUsers, userToPut)
	fileContentSequential, err := json.MarshalIndent(newUsers, "", "  ")

	if err != nil {
		return fmt.Errorf("failed to serialize updated sequential data: %w", err)
	}

	if err := os.WriteFile(DATA_FILE_SEQUENTIAL, fileContentSequential, 0644); err != nil {
		return fmt.Errorf("failed to write updated sequential user data: %w", err)
	}

	fmt.Println("User successfully updated:", userToPut)

	return nil
}

func GetDataHandler() ([]User, error) {
	fileContent, err := os.ReadFile(DATA_FILE_SEQUENTIAL)

	if err != nil {
		return nil, err
	}

	var users []User

	if err := json.Unmarshal(fileContent, &users); err != nil {
		return nil, err
	}

	return users, nil
}

func DeleteDataHandler(userToDelete User) error {
	currentFileContent, err := os.ReadFile(DATA_FILE)

	if err != nil {
		return fmt.Errorf("failed to read %s: %v", DATA_FILE, err)
	}

	currentFileContentSequential, err := os.ReadFile(DATA_FILE_SEQUENTIAL)

	if err != nil {
		return fmt.Errorf("failed to read %s: %v", DATA_FILE_SEQUENTIAL, err)
	}

	var newUserMap map[string]User
	var newUsers []User

	if err := json.Unmarshal(currentFileContent, &newUserMap); err != nil {
		return fmt.Errorf("failed to unmarshal %s: %v", DATA_FILE, err)
	}

	delete(newUserMap, userToDelete.Id)
	fileContent, err := json.MarshalIndent(newUserMap, "", "  ")

	if err != nil {
		return fmt.Errorf("failed to marshal updated user map: %v", err)
	}

	if err := os.WriteFile(DATA_FILE, fileContent, 0644); err != nil {
		return fmt.Errorf("failed to write to %s: %v", DATA_FILE, err)
	}

	if err := json.Unmarshal(currentFileContentSequential, &newUsers); err != nil {
		return fmt.Errorf("failed to unmarshal %s: %v", DATA_FILE_SEQUENTIAL, err)
	}

	newUsers = erase(newUsers, userToDelete)
	fileContentSequential, err := json.MarshalIndent(newUsers, "", "  ")

	if err != nil {
		return fmt.Errorf("failed to marshal updated sequential user list: %v", err)
	}

	if err := os.WriteFile(DATA_FILE_SEQUENTIAL, fileContentSequential, 0644); err != nil {
		return fmt.Errorf("failed to write to %s: %v", DATA_FILE_SEQUENTIAL, err)
	}

	return nil
}
