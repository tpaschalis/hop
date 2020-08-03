package main

import "fmt"

func clearPlusUsage() {
	fmt.Printf("\033[2J")
	fmt.Printf("\033[1;1H")
	fmt.Printf(usage)
	fmt.Printf(`› `)
}
func clear() {
	fmt.Printf("\033[2J")
	fmt.Printf("\033[1;1H")
	fmt.Printf(`› `)
}
