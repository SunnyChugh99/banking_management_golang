package gapi

import (
	"fmt"

	db "github.com/SunnyChugh99/banking_management_golang/db/sqlc"
	"github.com/SunnyChugh99/banking_management_golang/pb"
	"github.com/SunnyChugh99/banking_management_golang/token"
	"github.com/SunnyChugh99/banking_management_golang/util"
)

//Server serves http request for banking system
type Server struct{
	pb.UnimplementedSimpleBankServer
	config util.Config
	store db.Store  // interaction with database
	tokenMaker token.Maker

} 


//NewServer creates a new HTTP server and sets up routing
func NewServer(config util.Config,store db.Store) (*Server, error){

	fmt.Println("new server")
	fmt.Println(config.TokenSymmetricKey)
	fmt.Println(len(config.TokenSymmetricKey))
	fmt.Println("new server-2")
	fmt.Println(config.DBDriver)
	fmt.Println(config.AccessTokenDuration)

	fmt.Println("DB SOURCE")
	fmt.Println(config.DBSource)

	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil{
		return nil, fmt.Errorf("cannot create token master: %w", err)
	}

	server := &Server{config: config, store: store, tokenMaker: tokenMaker,}

	return server, nil


}

