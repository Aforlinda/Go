package main

import (
	"fmt"
	"strings"
)

func main() {
	scores := map[string]int{
		"afor":   88,
		"austin": 77,
		"julius": 60,
		"linda":  97,
		"lisa":   44,
		"xena":   66,
	}

	fmt.Print("Enter the name you want to check the score: ")
	var inputName string
	fmt.Scanln(&inputName)

	// Convert input name to lowercase
	inputName = strings.ToLower(inputName)

	if score, ok := scores[inputName]; ok {
		fmt.Printf("%s's score is %d\n", strings.Title(inputName), score)
	} else {
		fmt.Printf("No information is found for %s\n", strings.Title(inputName))
	}
}
