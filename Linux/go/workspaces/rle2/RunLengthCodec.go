package main

import (
	"errors"
)

func Encode(input string) string {
	if input == "" {
		return ""
	}

	var result []byte
	for i := 0; i < len(input); {
		sym := input[i]
		run := 1
		for i+run < len(input) && input[i+run] == sym && run < 255 {
			run++
		}
		result = append(result, byte(run), sym)
		i += run
	}
	return string(result)
}

func Decode(s string) (string, error) {
	if len(s)%2 != 0 {
		return "", errors.New("")
	}

	var result []byte
	for i := 0; i < len(s); i += 2 {
		count := int(s[i])
		if count == 0 {
			return "", errors.New("")
		}
		sym := s[i+1]
		for j := 0; j < count; j++ {
			result = append(result, sym)
		}
	}
	return string(result), nil
}
