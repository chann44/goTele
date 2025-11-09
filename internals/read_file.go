package internals

import (
	"fmt"
	"os"
)

func ReadFile(path string) []byte {
	templateData, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("Error reading template file: %v", err)
		os.Exit(1)
	}
	return templateData
}
