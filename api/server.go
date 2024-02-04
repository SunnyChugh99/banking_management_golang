package api

import (
	"fmt"

	db "github.com/SunnyChugh99/banking_management_golang/db/sqlc"
	"github.com/SunnyChugh99/banking_management_golang/token"
	"github.com/SunnyChugh99/banking_management_golang/util"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

//Server serves http request for banking system
type Server struct{
	config util.Config
	store db.Store  // interaction with database
	tokenMaker token.Maker
	router *gin.Engine  //this router will help us send each api request to correct handler for processing
} 


//NewServer creates a new HTTP server and sets up routing
func NewServer(config util.Config,store db.Store) (*Server, error){

	fmt.Println("new server")
	fmt.Println(config.TokenSymmetricKey)
	fmt.Println(len(config.TokenSymmetricKey))
	fmt.Println("new server-2")
	fmt.Println(config.DBDriver)
	fmt.Println(config.AccessTokenDuration)

	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil{
		return nil, fmt.Errorf("cannot create token master: %w", err)
	}

	server := &Server{config: config, store: store, tokenMaker: tokenMaker,}

	if v,ok :=  binding.Validator.Engine().(*validator.Validate); ok{
		v.RegisterValidation("currency", validCurrency)
	}

	server.setUpRouter()

	return server, nil


}

func (server *Server) setUpRouter(){
	router := gin.Default()

	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	authRoutes.POST("/accounts", server.createAccount)

	authRoutes.GET("/accounts/:id", server.getAccount)
	authRoutes.GET("/accounts", server.listAccount)
	authRoutes.POST("/transfers", server.createTransferAccount)


	server.router = router

}

//server runs the HTTP server on a specific address
func (server *Server) Start(address string) error{
	return server.router.Run(address)
}


func errorResponse(err error) gin.H{
	return gin.H{"error": err.Error()}
}