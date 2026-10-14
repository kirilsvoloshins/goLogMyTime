package assert

import (
	"fmt"
	"os"
)

func Assert(condition bool, message string) {
	if !condition {
		fmt.Printf("Error: %v", message)
		os.Exit(2)
	}
}
