package utils

import (
	"encoding/json"
	"os"
)

func MustDumpJsonToFile(obj any, file string) {
	data, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		panic(err)
	}

	os.WriteFile(file, data, os.ModePerm)
}
