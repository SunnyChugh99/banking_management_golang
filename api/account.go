package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	db "github.com/SunnyChugh99/banking_management_golang/db/sqlc"
	"github.com/SunnyChugh99/banking_management_golang/token"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type CreateAccountRequest struct {
	Currency string `json:"currency" binding:"required,currency"`
}


//every gin handler will have gin.contet input, it will help us to read input and write out response
func (server *Server) createAccount(ctx *gin.Context){
	var req CreateAccountRequest
 	if err := ctx.ShouldBindJSON(&req); err!=nil{
		ctx.JSON(http.StatusBadRequest,  errorResponse(err))
		return 
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	args := db.CreateAccountParams{
		Owner: authPayload.Username,
		Balance: 0,
		Currency: req.Currency,
	}

	account, err := server.store.CreateAccount(ctx, args)
	if err!=nil{
		if pqErr,ok := err.(*pq.Error); ok{
			switch pqErr.Code.Name(){
			case "foreign_key_violation", "unique_violation":
				ctx.JSON(http.StatusForbidden,  errorResponse(err))
				return
			}
		}

		ctx.JSON(http.StatusInternalServerError,  errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)

}

type GetAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1" validate:"min=1"` // Use the validate tag
}

func (server *Server) getAccount(ctx *gin.Context){
	var req GetAccountRequest

 	if err := ctx.ShouldBindUri(&req); err!=nil{
		ctx.JSON(http.StatusBadRequest,  errorResponse(err))
		return 
	}

	account, err := server.store.GetAccount(ctx, req.ID)
	if err!=nil{
		if err == sql.ErrNoRows{
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return 
		}
		ctx.JSON(http.StatusInternalServerError,  errorResponse(err))
		return 
	}
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if account.Owner != authPayload.Username{
		err := errors.New("account doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized,  errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)

}




type ListAccountRequest struct {
	PageID    int32    `form:"page_id"  binding:"required,min=1"`
	PageSize     int32    `form:"page_size"  binding:"required,min=5,max=10"`
}
func (server *Server) listAccount(ctx *gin.Context){
	var req ListAccountRequest
	if err := ctx.Bind(&req); err!=nil{
		ctx.JSON(http.StatusBadRequest,  errorResponse(err))
		return 
	}

	limitValue := (req.PageSize)
	fmt.Println(limitValue)	
	offsetValue := (req.PageID * req.PageSize - req.PageSize)
	fmt.Println(offsetValue)


	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)


	args :=  db.ListAccountsParams{
		Owner: authPayload.Username,
		Limit: limitValue,
		Offset: offsetValue,
	}

	accountList, err := server.store.ListAccounts(ctx, args)
	if err !=nil {
		return
	}
	fmt.Println(accountList)
	
	ctx.JSON(http.StatusOK, accountList)

}