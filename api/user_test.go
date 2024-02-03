package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	mockdb "github.com/SunnyChugh99/banking_management_golang/db/mock"
	db "github.com/SunnyChugh99/banking_management_golang/db/sqlc"
	"github.com/SunnyChugh99/banking_management_golang/util"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)


func TestCreateUsertApi(t *testing.T){
	user, password := randomUser(t)
	
	testCases := []struct{
		name string
		body gin.H
		buildStubs func(store *mockdb.MockStore)
		checkResponse  func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{    
			"username": user.Username,
			"full_name": user.FullName,
			"email": user.Email,
			"password": password,
		},
			buildStubs: func(store *mockdb.MockStore){
				store.EXPECT().
				CreateUser(gomock.Any(), gomock.Any()).
				Times(1).
				Return(user, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder){
				require.Equal(t, http.StatusOK, recorder.Code)
				requiredBodyMatchAccount(t, recorder.Body, user)
			},
		},
	}

	for i:= range testCases{

	tc := testCases[i]
	t.Run(tc.name, func(t *testing.T){
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		store := mockdb.NewMockStore(ctrl)
		tc.buildStubs(store)
	
	
		// start test server and send request
		server := NewServer(store)
		recorder := httptest.NewRecorder()
	
		url := fmt.Sprintf("/accounts/%d", tc.accountId)
		request, err := http.NewRequest(http.MethodGet, url, nil)
		require.NoError(t, err)
		server.router.ServeHTTP(recorder, request)
		tc.checkResponse(t, recorder)
	})

	}

}

func randomUser(t *testing.T) (db.User, string){
	password := util.RandomString(6)
	hashedPassword, err := util.HashPassword(password)
	require.NoError(t,err)
	require.NotEmpty(t, hashedPassword) 

	return db.User{
		Username: util.RandomString(6),
		HashedPassword: hashedPassword,
		FullName: util.RandomString(6),
		Email: util.RandomString(6),
		PasswordChangedAt: time.Now(),
		CreatedAt: time.Now(),
	}, hashedPassword
}

// func requiredBodyMatchAccount(t *testing.T, body *bytes.Buffer, account db.Account){
// 	data , err := io.ReadAll(body)
// 	require.NoError(t, err)
	

// 	var gotAccount db.Account
// 	err = json.Unmarshal(data, &gotAccount)
// 	require.NoError(t, err)
// 	require.Equal(t, account, gotAccount)


// }