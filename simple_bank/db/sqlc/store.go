package db

import (
	"context"
	"database/sql"
	"fmt"
)

//Provides all functions to execute db queries and transactions
// Store struct holds info about the Queries which in itself encapsulates DBTX interface(containing methods for db operations)
// db variable seems to be a reference to db connection or db pool
type Store struct{
	*Queries 
	db *sql.DB
}

//NewStore creates a new store
func NewStore(db *sql.DB) *Store{
	return &Store{
		db: db,
		Queries:  New(db),
	}
}


func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error{
	tx, err := store.db.BeginTx(ctx, nil)

	if err!=nil{
		return err
	}


	q := New(tx)
	err = fn(q)
	if err != nil{
		if rbErr := tx.Rollback(); rbErr!=nil{
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	
	}
	return tx.Commit()
}


type TransferTxParams struct{
	FromAccountId int64 `json:"from_account_id"`
	ToAccountId int64 `json:"to_account_id"`
	Amount int64 `json:"amount"`
}


type TransferTxResult struct{
	Transfer Transfer `json:"transfer"`
	FromAccount Account `json:"from_account"`
	ToAccount Account `json:"to_account"`
	FromEntry Entry `json:"from_entry"`
	ToEntry Entry `json:"to_entry"`
}

//TransferTx performs transfer from one account to another
// It creates a transfer record, add account entries, and update account's balance within a single transaction

func (store *Store) TransferTx(ctx context.Context, args TransferTxParams) (TransferTxResult,error){
	var result TransferTxResult
	//var err error
	err := store.execTx(ctx, func(q *Queries) error{
		var err error
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
				FromAccountID: args.FromAccountId,
				ToAccountID: args.ToAccountId,
				Amount: args.Amount,
			})
		if err !=nil{
			return err
		}


		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: args.FromAccountId,
			Amount: -args.Amount,
		})
		if err !=nil{
			return err
		}

		//entry of to_account_id
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: args.ToAccountId,
			Amount: args.Amount,
		})

		if err !=nil{
			return err
		}


		//TODO: update accounts balance

		if args.FromAccountId < args.ToAccountId{
			result.FromAccount, result.ToAccount, err =  addMoney(ctx, q, args.FromAccountId, -args.Amount, args.ToAccountId, args.Amount)

		}else{
			result.ToAccount, result.FromAccount, err = addMoney(ctx, q, args.ToAccountId, args.Amount, args.FromAccountId, -args.Amount)

		}

		if err!= nil{
			return err
		}

		return nil
	})
	return result, err
}



func addMoney(
	ctx context.Context,
	q * Queries,
	accountID1 int64,
	amount1 int64,
	accountID2 int64,
	amount2 int64,
)(account1 Account, account2 Account, err error){

	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID: accountID1,
		Amount: amount1,
	})

	if err!=nil{
		return
	}


	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID: accountID2,
		Amount: amount2,
	})

	return 	
}