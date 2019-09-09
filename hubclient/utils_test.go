package hubclient

import (
	"fmt"
	"math/rand"
	"unicode"
)

// rot13(alphabets) + rot5(numeric)
func rot13rot5(input string) string {

	var result []rune
	rot5map := map[rune]rune{}
	for i := 0; i <= 9; i++ {
		r := rune('0' + i)
		rot5map[r] = rune('0' + (i+5)%10)
	}

	for _, i := range input {
		switch {
		case !unicode.IsLetter(i) && !unicode.IsNumber(i):
			result = append(result, i)
		case i >= 'A' && i <= 'Z':
			result = append(result, 'A'+(i-'A'+13)%26)
		case i >= 'a' && i <= 'z':
			result = append(result, 'a'+(i-'a'+13)%26)
		case i >= '0' && i <= '9':
			result = append(result, rot5map[i])
		case unicode.IsSpace(i):
			result = append(result, ' ')
		}
	}

	return fmt.Sprintf(string(result[:]))
}

func randomString(n int) string {
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"

	b := make([]byte, n)
	for i := 0; i < n; i++ {
		b[i] = chars[rand.Intn(len(chars))]
	}

	return string(b)
}
