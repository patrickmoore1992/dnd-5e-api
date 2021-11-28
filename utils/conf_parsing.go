package utils

import (
	"encoding/json"
	"io/ioutil"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

// API Conf JSON Structure.
type Conf struct {
	DatabaseConfigs []DatabaseConfig `json:"mysqldb"`
}

// Mysql Connection JSON Structure.
type DatabaseConfig struct {
	Host     string `json:"DB_HOST"`
	Port     string `json:"DB_PORT"`
	User     string `json:"DB_USER"`
	Password string `json:"DB_PASSWORD"`
	Schema   string `json:"DB_SCHEMA"`
}

// Function for parsing the API conf.json.
func parseDatabaseConf() DatabaseConfig {
	jsonFile, err := os.Open("../../.conf.json")

	// if we os.Open returns an error then handle it
	if err != nil {
		panic(err.Error())
	}

	// Parse DatabaseConfig object out of conf.json
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var conf Conf
	json.Unmarshal(byteValue, &conf)
	return conf.DatabaseConfigs[0]
}
