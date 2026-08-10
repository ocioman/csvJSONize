package main

import (
	"csvJSONize/converter"
	"fmt"
	"os"
)

// only for testing
func main() {
	jsonFile, err := os.Create("res.json")
	if err != nil {
		return
	}

	if err := converter.Serialize(os.Stdin, jsonFile); err != nil {
		jsonFile.Close()
		return
	}
	jsonFile.Close()

	jsonInput, err := os.Open("res.json")
	if err != nil {
		return
	}
	defer jsonInput.Close()

	csvFile, err := os.Create("res.csv")
	if err != nil {
		return
	}
	defer csvFile.Close()

	if err := converter.Deserialize(jsonInput, csvFile); err != nil {
		fmt.Println(err)
	}
}
