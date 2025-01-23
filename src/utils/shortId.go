package utils

import (
	"fmt"
	"math/rand"
	"time"
)

func GenerateRandomID() (string, error) {
	rand.Seed(time.Now().UnixNano())
	randomNumber := rand.Intn(999999) + 1
	randomChar1 := 'a' + rune(rand.Intn(26))
	randomChar2 := 'a' + rune(rand.Intn(26))
	id := fmt.Sprintf("%c%c%06d", randomChar1, randomChar2, randomNumber)

	return id, nil
}
