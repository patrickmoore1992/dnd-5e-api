package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

// album represents data about a record album.
type Race struct {
	ID          int    `json:"id"`
	RaceName    string `json:"race_name"`
	Description string `json:"description"`
	InsertTS    string `json:"insert_ts"`
}

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

var dbConf DatabaseConfig

// Main driver method.
func main() {
	dbConf = parseDatabaseConf()
	router := gin.Default()
	router.GET("/races", getRaces)
	router.GET("/races/:name", getRacesByName)
	router.Run("localhost:8080")
}

// getRaces responds with the list of all races & descriptions as JSON.
func getRaces(c *gin.Context) {
	var records []Race
	var results = readFromMySQL(dbConf, "SELECT * FROM races")
	for results.Next() {
		var race Race
		// for each row, scan the result into our tag composite object
		var err = results.Scan(&race.ID, &race.RaceName, &race.Description, &race.InsertTS)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		records = append(records, race)
	}
	c.JSON(http.StatusOK, records)
}

// Handler for returning a race given a 'name' parameter.
func getRacesByName(c *gin.Context) {
	var records []Race
	var name = c.Param("name")
	var results = readFromMySQL(dbConf, fmt.Sprintf("SELECT * FROM races WHERE race_name = '%s'", name))

	// Probably a better way to do this.
	for results.Next() {
		var race Race
		// for each row, scan the result into our tag composite object
		var err = results.Scan(&race.ID, &race.RaceName, &race.Description, &race.InsertTS)
		if err != nil {
			fmt.Println(err.Error())
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		records = append(records, race)
	}

	if len(records) == 0 {
		c.JSON(http.StatusBadRequest, "Please choose a valid 5th edition race: [Dwarf, Human, Tiefling, Half-elf, Half-orc, Gnome, Halfling, Elf, Orc]")
	} else {
		c.JSON(http.StatusOK, records[0])
	}
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

// Helper function for opening connection to database.
func readFromMySQL(dbconf DatabaseConfig, query string) *sql.Rows {
	var HOST = dbconf.Host
	var PORT = dbconf.Port
	var USER = dbconf.User
	var PASS = dbconf.Password
	var DB = dbconf.Schema

	// Open up our database connection.
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", USER, PASS, HOST, PORT, DB))

	// if there is an error opening the connection, handle it
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	var results, readErr = db.Query(query)
	if readErr != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	return results
}
