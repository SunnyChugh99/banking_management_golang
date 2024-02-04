package api

import (
	"database/sql"
	"fmt"
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


type userResponse struct {
	Username          string    `json:"username"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

func newUserResponse(user db.User) userResponse{
	userResponse1 := userResponse{
		Username: user.Username,
		FullName: user.FullName,
		Email: user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
	}
	return userResponse1
}

//every gin handler will have gin.contet input, it will help us to read input and write out response
func (server *Server) createUser(ctx *gin.Context){
	var req CreateUserRequest
	fmt.Println("here-1")
 	if err := ctx.ShouldBindJSON(&req); err!=nil{
		ctx.JSON(http.StatusBadRequest,  errorResponse(err))
		return 
	}

	fmt.Println("here-2")

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
	userResponse := newUserResponse(user)
	ctx.JSON(http.StatusOK, userResponse)

}

type loginUserRequest struct {
	Username    string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginUserResponse struct{
	AccessToken string `json:"access_token"`
	User userResponse `json:"user"` 
}

func (server *Server) loginUser(ctx *gin.Context){
	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err!=nil{
		ctx.JSON(http.StatusBadRequest,  errorResponse(err))
		return 
	}
	user, err := server.store.GetUser(ctx, req.Username)
	if err !=nil{
		if err == sql.ErrNoRows{
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return 
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return 
	}	

	err = util.CheckPassword(req.Password, user.HashedPassword)
		
	if err!=nil{
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}


	accessToken, err := server.tokenMaker.CreateToken(
		user.Username,
		server.config.AccessTokenDuration,
	)

	if err!=nil{
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return 
	}
	
	rsp := loginUserResponse{
		AccessToken: accessToken,
		User: newUserResponse(user),
	}
	ctx.JSON(http.StatusOK, rsp)
}