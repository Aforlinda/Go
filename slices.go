package main

import "fmt"

func main() {
	names := []string{"Afor", "Susan", "Odoma", "Janet", "Linda", "James"}
	fmt.Println(names)

	names = append(names, "Mat")
	fmt.Println(names)

}
