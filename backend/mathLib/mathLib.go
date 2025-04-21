package mathLib

import (
	"JovaCentral/database"
	cryptoLibImport "crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	mathLibImport "math/rand"
	"os"
)

type SignedInt interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

type UnsignedInt interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

type Int interface {
	SignedInt | UnsignedInt
}

type Float interface {
	~float32 | ~float64
}

type Number interface {
	Int | Float
}

// Return the smallest of 2 numbers
func Min[T Number](a, b T) T {
	if a < b {
		return a
	}

	return b
}

// Return the largest of 2 numbers
func Max[T Number](a, b T) T {
	if a > b {
		return a
	}

	return b
}

// Clamp a value between min and max
func Clamp[T Number](value, minValue, maxValue T) T {
	return Max(Min(value, maxValue), minValue)
}

// Generate a random int between min and max, inclusive
func randNumRange[T Number](min, max T) T {
	return T(uint64(min) + uint64(mathLibImport.Int63n(int64(max-min+1))))
}

var alphanumericRanges = [][2]int{
	{48, 57},  // 0 - 9
	{65, 90},  // A - Z
	{97, 122}, // a - z
}

// Generates a string with the specified amount of random characters.
// Defaults to 64 characters if no length is provided.
func RandomString(_stringLength ...uint8) string {
	// If alphanumeric is false, certain special characters
	// like [].!, will be allowed in the string.
	var alphanumeric = true
	var stringLength uint8

	if len(_stringLength) == 0 {
		stringLength = 64 // Default value
	} else {
		// Assign the passed value, after clamping it
		stringLength = Clamp(_stringLength[0], 8, 64)
	}

	currentFileContent, _ := os.ReadFile(database.DATA_FILE)
	var newUserMap map[string]database.User

	// Deserialize the json file that contains user IDs as keys.
	// We're doing this so we can validate if a key exists.
	if jsonErr := json.Unmarshal(currentFileContent, &newUserMap); jsonErr != nil {
		return fmt.Sprintf("%v", jsonErr)
	}

	var returnString string
	var tempStringArray []rune

	if alphanumeric {
		for i := 0; i < int(stringLength); i++ {
			// Pick a number between 0-2 and select the set of characters
			// from the 'alphanumericRanges' array.
			characterRange := alphanumericRanges[mathLibImport.Intn(3)]
			// Add the random character to the temporary array
			tempStringArray = append(tempStringArray, rune(randNumRange(characterRange[0], characterRange[1])))
		}
	} else {
		for i := 0; i < int(stringLength); i++ {
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
			tempStringArray = append(tempStringArray, rune(mathLibImport.Intn(94)+33))
		}
	}

	returnString = string(tempStringArray)

	// If the ID already exists, recursively call this function until we don't have one.
	// Assuming we have 64 characters, an ID conflict is... unlikely, to say the least.
	if _, exists := newUserMap[returnString]; exists {
		return RandomString(stringLength)
	}

	return returnString
}

func RandomToken(_tokenLength ...uint8) string {
	var tokenLength uint8

	if len(_tokenLength) == 0 {
		tokenLength = 64
	} else {
		tokenLength = Clamp(_tokenLength[0], 8, 64)
	}

	token := make([]byte, tokenLength)

	if _, err := cryptoLibImport.Read(token); err != nil {
		fmt.Println("Error generating random token:", err)
		return ""
	}

	return base64.RawURLEncoding.EncodeToString(token)
}
