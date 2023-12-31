package utils

import (
	"encoding/json"
	"fmt"
	"os"
)

func MustDumpJsonToFile(obj any, file string) {
	data, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		panic(err)
	}

	os.WriteFile(file, data, os.ModePerm)
}

func MustPrintJson(obj any) {
	data, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(data))
}
