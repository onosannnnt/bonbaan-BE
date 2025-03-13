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
	randomNumber := rand.Intn(99999999) + 1
	randomChar1 := 'A' + rune(rand.Intn(26)) // Generate the first random character between 'a' and 'z'
	randomChar2 := 'A' + rune(rand.Intn(26)) // Generate the second random character between 'a' and 'z'
	id := fmt.Sprintf("%c%c%06d", randomChar1, randomChar2, randomNumber)

	regexPattern := `^[A-Z][A-Z]\d{8}$`
	matched, err := regexp.MatchString(regexPattern, id)
	if err != nil {
		return "", err
	}
	if !matched {
		return "", fmt.Errorf("generated ID does not match the required format")
	}

	return id, nil
}

func GenerateOTP(length int) (string, error) {
	const digits = "0123456789"
	rand.Seed(time.Now().UnixNano())
	var result string
	for i := 0; i < length; i++ {
		result += string(digits[rand.Intn(len(digits))])
	}

	regexPattern := fmt.Sprintf(`^\d{%d}$`, length)
	matched, err := regexp.MatchString(regexPattern, result)
	if err != nil {
		return "", err
	}
	if !matched {
		return "", fmt.Errorf("generated ID does not match the required format")
	}

	return result, nil
}
