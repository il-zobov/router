package config

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"sync"
	"time"
)

type ClassOfService struct {
	Paid       uint   `toml:"paid"`
	Location   string `toml:"location"`
	Max_points int    `toml:"max_points"`
}

type NetworkConf struct {
	HeaderName          string        `toml:"headerName"`
	BindAddr            string        `toml:"bind_addr"`
	LogLevel            string        `toml:"log_level"`
	CacheDefaultTimeout time.Duration `toml:"cacheDefaultTimeout"`
}
type DBConf struct {
	DBUser            string        `toml:"db_User"`
	DBPass            string        `toml:"db_Pass"`
	DBAddr            string        `toml:"db_Addr"`
	DBPort            string        `toml:"db_Port"`
	DBName            string        `toml:"db_Name"`
	DBConnMaxLifetime time.Duration `toml:"db_ConnMaxLifetime"`
	DBMaxOpenConns    int           `toml:"db_MaxOpenConns"`
	DBMaxIdleConns    int           `toml:"db_MaxIdleConns"`
}
type TypeOfService struct {
	Location    string `toml:"location"`
	AccountName string
}

type Config struct {
	NetworkConf NetworkConf
	DBConf      DBConf
	Classes     []ClassOfService
	FixPlan     map[string]TypeOfService `toml:"fixPlan"`
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		instance = &Config{}
		if _, err := toml.DecodeFile("config.toml", instance); err != nil {
			fmt.Println("Error: ", err)
			return
		}
	})
	return instance
}
