package main

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
