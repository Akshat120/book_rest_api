package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

type DatabaseConfig struct {
	Type             string `json:"type"`
	Host             string `json:"host"`
	Port             int    `json:"port"`
	User             string `json:"user"`
	Password         string `json:"password"`
	ConnectionString string `json:"connection_string"`
	DbName           string `json:"db_name"`
}

type BasicAuthConfig struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ServerConfig struct {
	Port         int    `json:"port"`
	Addr         string `json:"address"`
	WriteTimeout int    `json:"write_timeout"`
	ReadTimeout  int    `json:"read_timeout"`
	IdleTimeout  int    `json:"idle_timeout"`
}

type Config struct {
	AppName   string          `json:"app_name"`
	Debug     bool            `json:"debug"`
	Server    ServerConfig    `json:"server"`
	Database  DatabaseConfig  `json:"database"`
	BasicAuth BasicAuthConfig `json:"basic_auth"`
}

var appConfig *Config

func InitConfig() (*Config, error) {

	// Check if appConfig is already initialized
	if appConfig != nil {
		return appConfig, nil
	}

	// Load the configuration from the JSON file
	// If the file does not exist, return an error
	var err error
	appConfig, err = loadConfig("config.json")
	if err != nil {
		return nil, fmt.Errorf("error loading config.json file")
	}

	return appConfig, nil
}

func GetConfig() *Config {
	if appConfig == nil {
		panic("Configuration not initialized. Call InitConfig() first.")
	}
	return appConfig
}

func loadConfig(filename string) (*Config, error) {
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
