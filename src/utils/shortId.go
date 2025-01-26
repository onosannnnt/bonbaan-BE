package utils

import (
	"fmt"
	"math/rand"
	"regexp"

	// "regexp"
	"time"
)

// GenerateRandomID generates a random ID with a random character in front in the format "nt000001"
func GenerateRandomID() (string, error) {
	rand.Seed(time.Now().UnixNano())
	randomNumber := rand.Intn(999999) + 1
	randomChar1 := 'A' + rune(rand.Intn(26)) // Generate the first random character between 'a' and 'z'
	randomChar2 := 'A' + rune(rand.Intn(26)) // Generate the second random character between 'a' and 'z'
	id := fmt.Sprintf("%c%c%06d", randomChar1, randomChar2, randomNumber)

	regexPattern := `^[A-Z][A-Z]\d{6}$`
	matched, err := regexp.MatchString(regexPattern, id)
	if err != nil {
		return "", err
	}
	if !matched {
		return "", fmt.Errorf("generated ID does not match the required format")
	}

	return id, nil
}
