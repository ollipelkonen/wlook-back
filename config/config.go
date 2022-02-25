package config

/*
	this package loads configurations from given JSON file
*/

import (
	"encoding/json"
	"fmt"
	"os"
)

// struct to hold configurations
type Config struct {
	Port     string `json:"port"`
	Token    string `json:"token"`
	Database struct {
		Host     string `json:"host"`
		Port     string `json:"port"`
		Database string `json:"database"`
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"database"`
}

// function to load configurations, returns Config
func LoadConfig(filename string) Config {
	configFile, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Missing config file '%s'\n", filename)
		os.Exit(1)
	}
	if err != nil {
		fmt.Println(err.Error())
	}
	config := Config{}
	json.Unmarshal(configFile, &config)
	return config
}
