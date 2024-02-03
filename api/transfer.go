package api

import (
	"database/sql"
	"fmt"
	"net/http"

	db "github.com/SunnyChugh99/banking_management_golang/db/sqlc"
	"github.com/gin-gonic/gin"
)

// type CreateTransferParams struct {
// 	FromAccountID int64 `json:"from_account_id"`
// 	ToAccountID   int64 `json:"to_account_id"`
// 	Amount        int64 `json:"amount"`
// }

// func (q *Queries) CreateTransfer(ctx context.Context, arg CreateTransferParams) (Transfer, error) {


type TransferAccountRequest struct {
	FromAccountID int64 `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64 `json:"to_account_id" binding:"required,min=1"`
	Amount        int64 `json:"amount" binding:"required,gt=0"`
	Currency 	  string `json:"currency" binding:"required,currency"`
}




func (server *Server) createTransferAccount(ctx *gin.Context){
	var req TransferAccountRequest
 	if err := ctx.ShouldBindJSON(&req); err!=nil{
		ctx.JSON(http.StatusBadRequest,  errorResponse(err))
		return 
	}


	if !server.validAccount(ctx, req.FromAccountID, req.Currency){
		return
	}

	if !server.validAccount(ctx, req.ToAccountID, req.Currency){
		return
	}
	args := db.TransferTxParams{
		FromAccountId: req.FromAccountID,
		ToAccountId: req.ToAccountID,
		Amount: req.Amount,
	}

	result, err := server.store.TransferTx(ctx, args)
	if err!=nil{
		ctx.JSON(http.StatusInternalServerError,  errorResponse(err))
	}

	ctx.JSON(http.StatusOK, result)

}


func (server *Server) validAccount(ctx *gin.Context, accountID int64, currency string) bool{
	account, err := server.store.GetAccount(ctx, accountID)	
	if err!=nil{
		if err == sql.ErrNoRows{
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return false
		}
		ctx.JSON(http.StatusInternalServerError,  errorResponse(err))
		return false
	}

	if account.Currency != currency{
		err := fmt.Errorf("account [%d] currency mismatch: %s vs %s", accountID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return false
	}	

	return true
}