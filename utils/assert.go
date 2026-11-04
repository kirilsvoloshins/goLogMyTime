package utils

import (
	"fmt"
	"os"
)

func Assert(condition bool, message string) {
	if !condition {
		fmt.Printf("Assertion error: %v \n", message)
		os.Exit(2)
	}
}
