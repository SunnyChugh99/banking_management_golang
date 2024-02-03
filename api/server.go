package api

import (
	db "github.com/SunnyChugh99/banking_management_golang/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

//Server serves http request for banking system
type Server struct{
	store db.Store  // interaction with database
	router *gin.Engine  //this router will help us send each api request to correct handler for processing
} 


//NewServer creates a new HTTP server and sets up routing
func NewServer(store db.Store) *Server{
	server := &Server{store: store}
	router := gin.Default()

	if v,ok :=  binding.Validator.Engine().(*validator.Validate); ok{
		v.RegisterValidation("currency", validCurrency)
	}

	router.POST("/users", server.createUser)

	router.POST("/accounts", server.createAccount)

	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccount)
	router.POST("/transfers", server.createTransferAccount)


	server.router = router
	return server


}

//server runs the HTTP server on a specific address
func (server *Server) Start(address string) error{
	return server.router.Run(address)
}


func errorResponse(err error) gin.H{
	return gin.H{"error": err.Error()}
}