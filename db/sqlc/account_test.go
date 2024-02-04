package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/SunnyChugh99/banking_management_golang/util"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account{
	user := createRandomUser(t)
	args := CreateAccountParams{
		Owner: user.Username,
		Balance:util.RandomMoney(),
		Currency:util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(),args)
	
	//require will checl if there is error or not, if error found it will fail the test case itself.
	require.NoError(t,err)
	
	require.NotEmpty(t,account)
	require.Equal(t, args.Owner, account.Owner)
	require.Equal(t, args.Balance, account.Balance)
	require.Equal(t, args.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)
	return account
}
func TestCreateAccount(t *testing.T){
	createRandomAccount(t)

}



func TestGetAccount(t *testing.T){
	account := createRandomAccount(t)
	retrieveAccount, err := testQueries.GetAccount(context.Background(), account.ID)
	require.NoError(t,err)
	require.NotEmpty(t,retrieveAccount)
	require.Equal(t, retrieveAccount.ID, account.ID)
	require.Equal(t, retrieveAccount.Currency, account.Currency)

}



func TestListAccounts(t *testing.T){
	var lastAccount Account
	for i:=0; i<10; i++{
		lastAccount = createRandomAccount(t)
	}
	args := ListAccountsParams{
		Owner : lastAccount.Owner,
		Limit: 5,
		Offset: 0,
	}
	accounts, err := testQueries.ListAccounts(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, accounts)
	for _,account := range accounts{
		require.NotEmpty(t, account)
		require.Equal(t, lastAccount.Owner, account.Owner)
	}
}

func TestUpdateAccount(t *testing.T){
	account1 := createRandomAccount(t)

	args := UpdateAccountParams{
			ID: account1.ID,
			Balance:util.RandomMoney(),
		}
	
	account2, err := testQueries.UpdateAccount(context.Background(),args)

	require.NoError(t, err)

	require.NotEmpty(t, account2)
	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account2.Balance, args.Balance)

}

func TestDeleteAccount(t *testing.T){
	account1 := createRandomAccount(t)	
	err := testQueries.DeleteAccount(context.Background(),account1.ID)
	require.NoError(t, err)


	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, account2)

}