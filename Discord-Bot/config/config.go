package config

import (
	"encoding/json"
	"fmt"
	"os"
)

var (
	Token     string
	BotPrefix string
	config    *ConfigStruct
)

type ConfigStruct struct {
	Token     string `json:"Token"`
	BotPrefix string `json:"BotPrefix"`
}

func ReadConfig() error {
	fmt.Println("Reading config file...")

	file, err := os.ReadFile("./config.json")

	if err != nil {
		fmt.Printf("Error reading config file: %s\n", err.Error())
		return err
	}

	fmt.Println(string(file))

	config = &ConfigStruct{}
	err = json.Unmarshal(file, config)

	fmt.Println(config)
	fmt.Println("Token:", config.Token)
	fmt.Println("BotPrefix:", config.BotPrefix)

	if err != nil {
		fmt.Printf("Error unmarshalling config file: %s\n", err.Error())
		return err
	}

	return nil
}
