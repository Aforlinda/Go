package main

import (
	"fmt"
	"strings"
)

func main() {
	// List of fruits in the shop
	fruitsInShop := []string{"apple", "banana", "orange", "grapes", "mango"}

	//  A variable called shoppingCart is created to take in strings
	var shoppingCart []string

	// Display available fruits
	fmt.Println("Available Fruits in the Shop:")
	for _, fruit := range fruitsInShop {
		fmt.Println("-", fruit)
	}

	// User selects fruits to add to the shopping cart
	for {
		fmt.Print("Enter the name of the fruit to add to the cart (or 'done' to finish): ")
		var userInput string
		fmt.Scanln(&userInput)

		userInput = strings.ToLower(userInput)

		if userInput == "done" {
			break
		}

		// Check if the selected fruit is available
		found := false
		for _, fruit := range fruitsInShop {
			if userInput == fruit {
				// Add the selected fruit to the shopping cart
				shoppingCart = append(shoppingCart, userInput)
				fmt.Printf("%s added to the shopping cart.\n", strings.Title(userInput))
				found = true
				break

			}
		}

		if !found {
			fmt.Printf("Sorry, %s is not available in the shop.\n", strings.Title(userInput))
		}
	}

	// Display the items in the shopping cart
	fmt.Println("\nShopping Cart:")
	for _, item := range shoppingCart {
		fmt.Println("-", item)
	}
}
