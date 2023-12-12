package main

import "fmt"

func main() {
	// A variable called daily sales is created to take in floats for 5 days
	var dailysales [5]float64

	// Get input of daily sales from the user using the for loop
	fmt.Println("Enter Daily sales:")
	for i := 0; i < len(dailysales); i++ {
		fmt.Printf("Enter sales for Day %d: ", i+1)
		var sales float64
		fmt.Scanln(&sales)
		dailysales[i] = sales
	}

	// Display the daily sales for 5 days. Again the for loop is used
	fmt.Println("\nDaily Sales:")
	for i, sales := range dailysales {
		fmt.Printf("Day %d: %.2f Pounds\n", i+1, sales)
	}
}
