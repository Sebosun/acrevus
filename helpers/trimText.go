// Package helpers stores general helper functions
package helpers

func TrimText(input string) string {
	if len(input) == 0 {
		return ""
	}

	acc := ""
	// gets first 50 letters
	for i := 0; i < 50 || i < len(input); i++ {
		acc += string(input[i])
	}
	return acc
}
