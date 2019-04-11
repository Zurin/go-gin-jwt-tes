package config

import (
	"gin-jwt/structs"
	"io/ioutil"

	"github.com/buger/jsonparser"
	"github.com/jinzhu/gorm"
)

// DBInit create connection to database
func DBInit() *gorm.DB {
	json, err := ioutil.ReadFile("./config/config.json")
	dbengine, err := jsonparser.GetString(json, "databases", "[0]", "engine")
	dbconnstr, err := jsonparser.GetString(json, "databases", "[0]", "connstr")
	dbmaxconnopen, err := jsonparser.GetInt(json, "databases", "[0]", "maxconnopen")

	db, err := gorm.Open(dbengine, dbconnstr)
	if err != nil {
		panic("failed to connect to database")
	}

	db.DB().SetMaxIdleConns(int(dbmaxconnopen / 4))
	db.DB().SetMaxOpenConns(int(dbmaxconnopen))
	db.AutoMigrate(structs.Person{})
	return db
}
