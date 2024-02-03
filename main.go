package main

import (
	"database/sql"
	"log"

	"github.com/SunnyChugh99/banking_management_golang/api"
	db "github.com/SunnyChugh99/banking_management_golang/db/sqlc"
	"github.com/SunnyChugh99/banking_management_golang/util"
	_ "github.com/lib/pq"
)




func main(){

	config, err := util.LoadConfig(".")
	if err!=nil{
		log.Fatal("Cannot load config ", err)
	}




	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err!=nil{
		log.Fatal("Cannot connect to database")
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)
	err = server.Start(config.ServerAddress)
	if err!=nil{
		log.Fatal("Cannot start server", err)
	}
}

