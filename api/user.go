package api

import (
	"net/http"
	"time"

	db "github.com/SunnyChugh99/banking_management_golang/db/sqlc"
	"github.com/SunnyChugh99/banking_management_golang/util"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type CreateUserRequest struct {
	Username    string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	FullName    string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}


type CreateUserResponse struct {
	Username          string    `json:"username"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

//every gin handler will have gin.contet input, it will help us to read input and write out response
func (server *Server) createUser(ctx *gin.Context){
	var req CreateUserRequest
 	if err := ctx.ShouldBindJSON(&req); err!=nil{
		ctx.JSON(http.StatusBadRequest,  errorResponse(err))
		return 
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err!=nil{
		ctx.JSON(http.StatusInternalServerError,  errorResponse(err))
		return
	}
	args := db.CreateUserParams{
		Username: req.Username,
		HashedPassword: hashedPassword,
		FullName: req.FullName,
		Email: req.Email,
	}

	user, err := server.store.CreateUser(ctx, args)
	if err!=nil{
		if pqErr,ok := err.(*pq.Error); ok{
			switch pqErr.Code.Name(){
			case "unique_violation":
				ctx.JSON(http.StatusForbidden,  errorResponse(err))
				return
			}
		}

		ctx.JSON(http.StatusInternalServerError,  errorResponse(err))
		return
	}
	userResponse := CreateUserResponse{
		Username: user.Username,
		FullName: user.FullName,
		Email: user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
	}
	ctx.JSON(http.StatusOK, userResponse)

}