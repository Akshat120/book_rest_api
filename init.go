package main

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
)

var config *Config

func init() {
	var err error
	config, err = LoadConfig("config.json")
	if err != nil {
		log.Fatal("Error loading config.json file")
		return
	}
}

func LoadConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	config := &Config{}
	err = decoder.Decode(config)
	if err != nil {
		return config, err
	}
	setDefaults(config)
	return config, nil
}

func setDefaults(config *Config) {
	config.Server.Addr = getAddress(config.Server.Port)
	config.Database.ConnectionString = getDatabaseConnectString(
		config.Database.User,
		config.Database.Password,
		config.Database.Host,
		config.Database.Port,
		config.Database.DbName,
	)
}

func getAddress(port int) string {
	return ":" + strconv.Itoa(port)
}

func getDatabaseConnectString(user string, password string, host string, port int, dbName string) string {
	return (user + ":" + password + "@tcp(" + host + ":" + strconv.Itoa(port) + ")/" + dbName)
}
