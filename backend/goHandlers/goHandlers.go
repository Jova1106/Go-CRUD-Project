package goHandlers

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
)

const DATA_FILE = "../database/data.json"
const DATA_FILE_SEQUENTIAL = "../database/dataSequential.json"

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

// The default user data structure
type User struct {
	Id    string `json:id`
	Name  string `json:name`
	Email string `json:email`
}

// Return the smallest of 2 ints
func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}

// Return the largest of 2 ints
func max(a, b int) int {
	if a > b {
		return a
	}

	return b
}

// Clamp a value between min and max
func clamp(value int, minValue int, maxValue int) int {
	return max(min(value, maxValue), minValue)
}

// Generate a random int between min and max, inclusive
func randInt64Range(min int, max int) uint64 {
	return uint64(min) + uint64(rand.Int63n(int64(max-min+1)))
}

var alphanumericRanges = [][2]int{
	{48, 57},  // 0 - 9
	{65, 90},  // A - Z
	{97, 122}, // a - z
}

// This function repeats '_stringLength' times and generates
// a random alpnanumeric character that many times.
//
// Since Go doesn't allow optional parameters,
// I chose to use variadic args, and set a default
// if there are no parameters provided.
// Only the first parameter is evaluated, so
// randomString(7, 8, 9) will return a random string
// with 7 characters.
func randomString(_stringLength ...int) string {
	// If alphanumeric is false, certain special characters
	// like [].!, will be allowed in the string.
	var alphanumeric = true
	var stringLength int
	if len(_stringLength) == 0 {
		stringLength = 64 // Default value
	} else {
		// Assign the passed value, after clamping it
		stringLength = clamp(_stringLength[0], 0, 64)
	}
	currentFileContent, _ := os.ReadFile(DATA_FILE)
	var newUserMap map[string]User

	// Deserialize the json file that contains user IDs as keys.
	// We're doing this so we can validate if a key exists.
	if err := json.Unmarshal(currentFileContent, &newUserMap); err != nil {
		return fmt.Sprintf("%v", err)
	}

	var returnString string
	var tempStringArray []rune

	if alphanumeric {
		for i := 0; i < stringLength; i++ {
			// Pick a number between 0-2 and select the set of characters
			// from the 'alphanumericRanges' array.
			characterRange := alphanumericRanges[rand.Intn(3)]
			// Add the random character to the temporary array
			tempStringArray = append(tempStringArray, rune(randInt64Range(characterRange[0], characterRange[1])))
		}
	} else {
		for i := 0; i < stringLength; i++ {
			// Set the range to be between 32 and 126 in the ASCII table,
			// Allowing special ASCII characters.
			//
			// 32 - 47	 : Space and basic punctuation characters, 	EXAMPLE: ! " # *
			// 58 - 64	 : Additional punctuation and symbols, 		EXAMPLE: : ; < =
			// 91 - 96	 : More punctuation and special symbols, 	EXAMPLE: [ ] \ ^
			// 123 - 126 : Brackets, pipe, and tilde, 				EXAMPLE: { } | ~
			//
			// Since rand.Intn(94) returns 93 (n - 1),
			// we need to add 33 to get 126.
			tempStringArray = append(tempStringArray, rune(rand.Intn(94)+33))
		}
	}

	returnString = string(tempStringArray)

	// If the ID already exists, recursively call this function until we don't have one.
	// Assuming we have 64 characters, an ID conflict is... unlikely, to say the least.
	if _, exists := newUserMap[returnString]; exists {
		return randomString(stringLength)
	}

	return returnString
}

// Reads the JSON database files and gives the new User instance a unique ID.
// Adds the User to the user list and writes the data back to the JSON database files.
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

// Reads the JSON database files and updates the existing User instance.
// Writes the data back to the JSON database files.
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

// Returns the sequential JSON database file.
// The sequential database file is used to display
// the users in the order that they were created.
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

// Reads the JSON database files and deletes the User instance.
// Writes the data back to the JSON database files.
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
