package main

import (
	"csvJSONize/converter"
	"os"
)

// only for testing
func main() {
	file, fileErr := os.Create("res.json")

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			return
		}
	}(file)

	if fileErr != nil {
		return
	}

	convErr := converter.Serialize(os.Stdin, file)

	if convErr != nil {
		return
	}
}
