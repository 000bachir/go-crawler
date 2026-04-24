package main

import "fmt"

func main() {
	test := "bachir"
	secondTest := &test
	thirdTest := *secondTest
	fmt.Printf("the memory adress is : %p\n", secondTest)
	fmt.Println("the value is : \n", thirdTest)
}
