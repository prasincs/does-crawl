package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type DBConfig struct {
	Type string `json:"type"`
}

type Configuration struct {
	Db DBConfig `json:"db"`
}

// From http://stackoverflow.com/questions/16465705/how-to-handle-configuration-in-go
func GetConfiguration(filePath string) Configuration {
	// pwd, err := os.Getwd()
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }
	// filePath := path.Join(pwd, "conf", "conf.json")
	//fmt.Println(filePath)
	file, _ := os.Open(filePath)
	decoder := json.NewDecoder(file)

	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("error:", err)
	}

	return configuration
}
