package helpers

import (
	"encoding/json"
	"fmt"
	"os"
)

func CreateModel(modelName string, u interface{}) error{
	bytes, err := json.Marshal(u)
	if err != nil{
		return err
	}

	bytes = append(bytes, '\n')

	file, err := os.OpenFile(
		fmt.Sprintf("./datastore/%s.txt", modelName),
		os.O_WRONLY|os.O_CREATE|os.O_APPEND,
		0600,
	)
	if err!=nil{
		return err
	}

	defer file.Close()

	_, err = file.Write(bytes)
	if err!=nil {
		return err
	}

	return nil
}
