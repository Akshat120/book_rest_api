package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/mcuadros/go-defaults"
)

// DatabaseConfig holds the configuration for the database connection
type DatabaseConfig struct {
	Type             string `json:"type" default:"mysql"`        // e.g., "mysql", "postgres", etc.
	Host             string `json:"host" default:"localhost"`    // Default MySQL host
	Port             int    `json:"port" default:"3306"`         // Default MySQL port
	User             string `json:"user" default:"root"`         // Default MySQL user
	Password         string `json:"password" default:"password"` // Default MySQL password
	DbName           string `json:"db_name" default:"my_books"`  // Default database name
	ConnectionString string `json:"connection_string"`
}

// BasicAuthConfig holds the configuration for basic authentication
type BasicAuthConfig struct {
	Username string `json:"username" default:"admin"`    // Default username for basic auth
	Password string `json:"password" default:"password"` // Default password for basic auth
}

// ServerConfig holds the configuration for the HTTP server
type ServerConfig struct {
	Host         string `json:"host" default:"localhost"`   // Default host for the server
	Port         int    `json:"port" default:"8080"`        // Default port for the server
	Addr         string `json:"address"`                    // Address to bind the server to
	WriteTimeout int    `json:"write_timeout" default:"20"` // Write timeout in seconds
	ReadTimeout  int    `json:"read_timeout" default:"20"`  // Read timeout in seconds
	IdleTimeout  int    `json:"idle_timeout" default:"20"`  // Idle timeout in seconds
}

// RateLimiter configuration
type RateLimiterConfig struct {
	Enabled bool    `json:"enabled" default:"true"` // Enable or disable rate limiting
	Rate    float64 `json:"rate" default:"100"`     // requests per second
	Burst   int     `json:"burst" default:"20"`     // maximum burst size
}

// Config holds the application configuration
type Config struct {
	AppName     string            `json:"app_name" default:"BooksApi"`
	Debug       bool              `json:"debug" default:"false"` // Enable or disable debug mode
	Server      ServerConfig      `json:"server"`
	Database    DatabaseConfig    `json:"database"`
	BasicAuth   BasicAuthConfig   `json:"basic_auth"`
	RateLimiter RateLimiterConfig `json:"rate_limiter"`
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
	defaults.SetDefaults(config)
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
