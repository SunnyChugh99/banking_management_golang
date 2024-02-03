package db

import (
	"database/sql"
	"log"
	"testing"

	"github.com/SunnyChugh99/banking_management_golang/util"
	_ "github.com/lib/pq"
)



var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M){
	config, err := util.LoadConfig("../..")
	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err!=nil{
		log.Fatal("Cannot connect to database")
	}

	testQueries = New(testDB)

	m.Run()
}