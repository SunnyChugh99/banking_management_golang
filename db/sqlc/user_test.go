package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/SunnyChugh99/banking_management_golang/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User{

	hashedPassword, err:= util.HashPassword(util.RandomString(6))
	require.NoError(t,err)
	require.NotEmpty(t, hashedPassword) 
	args := CreateUserParams{
		Username: util.RandomOwner(),
		HashedPassword:hashedPassword,
		FullName:util.RandomOwner(),
		Email: util.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(),args)
	
	//require will checl if there is error or not, if error found it will fail the test case itself.
	require.NoError(t,err)
	
	require.NotEmpty(t,user)
	require.Equal(t, args.Username, user.Username)
	require.Equal(t, args.HashedPassword, user.HashedPassword)
	require.Equal(t, args.FullName, user.FullName)

	require.True(t, user.PasswordChangedAt.IsZero())

	require.NotZero(t, user.CreatedAt)
	return user
}
func TestCreateUser(t *testing.T){
	createRandomUser(t)

}



func TestGetUser(t *testing.T){
	user1 := createRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user1.Username)
	require.NoError(t,err)
	require.NotEmpty(t,user2)
	require.Equal(t, user2.Username, user1.Username)
	require.Equal(t, user2.HashedPassword, user1.HashedPassword)
	require.Equal(t, user2.FullName, user1.FullName)
	require.Equal(t, user2.Email, user1.Email)
	require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)


}

func TestUpdateUser(t *testing.T){
	oldUser := createRandomUser(t)
	newFullName := util.RandomOwner()
	args := UpdateUserParams{
		HashedPassword: sql.NullString{
			String: "",
			Valid: false,
		},
		FullName: sql.NullString{
			String: newFullName,
			Valid: true,
		},
		Email: sql.NullString{
			String: "",
			Valid: false,
		},
		Username: oldUser.Username,
	}
	NewUser, err := testQueries.UpdateUser(context.Background(), args)

	require.NoError(t,err)
	require.NotEmpty(t,NewUser)
	require.NotEqual(t, oldUser.FullName, NewUser.FullName)
	require.Equal(t, oldUser.Email, NewUser.Email)
	require.Equal(t, oldUser.HashedPassword, NewUser.HashedPassword)
}