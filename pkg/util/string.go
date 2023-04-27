package util

import "strings"

func ToTitleCase(str string) string {
	words := strings.Split(str, " ")
	for i := range words {
		words[i] = strings.Title(words[i])
	}
	return strings.Join(words, "")
}

func ToSnakeCase(str string) string {
	words := strings.Split(str, " ")
	for i := range words {
		words[i] = strings.ToLower(words[i])
	}
	return strings.Join(words, "_")
}
