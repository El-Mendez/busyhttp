package main

import (
	"log"
	"os"
	"strconv"
	"strings"
)

func validateIntFlag(key string) int {
	value := os.Getenv(key)
	if value == "" {
		return 0
	}

	numberValue, err := strconv.Atoi(value)
	if err != nil {
		log.Fatalf("Invalid value of %s. Received %v.\n", key, value)
	}

	return numberValue
}

func SplitAtLast(s, sep string) string {
	index := strings.LastIndex(s, sep)
	if index == -1 {
		return s
	}
	return s[index-1+len(sep):]
}
