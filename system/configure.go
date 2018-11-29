package system

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
)

type Config struct {
	Postgres struct {
		Host     string `json:"host"`
		Port     int    `json:"port"`
		Password string `json:"password"`
		User     string `json:"user"`
		DBName   string `json:"db_name"`
		SSLMode  string `json:"ssl_mode"`
	} `json:"postgres"`
	Web struct {
		Port int `json:"port"`
	} `json:"web"`
}

func (c *Config) load() {
	file, err := os.Open("./config.json")
	if err != nil {
		log.Fatalf("[OPEN CONFIG FILE] %v", err)
	}
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(c); err != nil {
		log.Fatalf("[READ CONFIG FILE] %v", err)
	}
}

func (c *Config) ConnectPostgresString() string {
	return fmt.Sprintf("dbname=%v user=%v password=%v host=%v port=%v sslmode=%v",
		c.Postgres.DBName,
		c.Postgres.User,
		c.Postgres.Password,
		c.Postgres.Host,
		c.Postgres.Port,
		c.Postgres.SSLMode)
}

func (c *Config) Port() string {
	return ":" + strconv.Itoa(c.Web.Port)
}

var (
	_once   sync.Once
	_config *Config
)

func GetConfig() *Config {
	_once.Do(func() {
		_config = new(Config)
		_config.load()
	})
	return _config
}
