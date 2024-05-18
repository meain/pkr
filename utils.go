package main

import "fmt"

func getCompletions(input string) []string {
	outs := []string{}

	if len(input) == 0 {
		return outs
	}

	for i := 1; i <= 10; i++ {
		outs = append(outs, fmt.Sprintf("%s %d", input, i))
	}

	return outs
}
