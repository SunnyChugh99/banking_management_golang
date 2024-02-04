package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	db "github.com/SunnyChugh99/banking_management_golang/db/sqlc"
	"github.com/SunnyChugh99/banking_management_golang/token"
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

	fromAccount, valid := server.validAccount(ctx, req.FromAccountID, req.Currency)


	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	if fromAccount.Owner != authPayload.Username{
		err := errors.New("from account doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized,  errorResponse(err))
		return
	}

	if !valid{
		return
	}

	_, valid = server.validAccount(ctx, req.ToAccountID, req.Currency)
	if !valid{
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


func (server *Server) validAccount(ctx *gin.Context, accountID int64, currency string) (db.Account,bool){
	account, err := server.store.GetAccount(ctx, accountID)	
	if err!=nil{
		if err == sql.ErrNoRows{
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return account, false
		}
		ctx.JSON(http.StatusInternalServerError,  errorResponse(err))
		return account, false
	}

	if account.Currency != currency{
		err := fmt.Errorf("account [%d] currency mismatch: %s vs %s", accountID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return account, false
	}	

	return account, true
}