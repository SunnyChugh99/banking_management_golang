package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	mockdb "github.com/SunnyChugh99/banking_management_golang/db/mock"
	db "github.com/SunnyChugh99/banking_management_golang/db/sqlc"
	"github.com/SunnyChugh99/banking_management_golang/util"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)


func TestGetAccountApi(t *testing.T){
	account := randomAccount()
	
	testCases := []struct{
		name string
		accountId int64
		buildStubs func(store *mockdb.MockStore)
		checkResponse  func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			accountId: account.ID,
			buildStubs: func(store *mockdb.MockStore){
				store.EXPECT().
				GetAccount(gomock.Any(), gomock.Eq(account.ID)).
				Times(1).
				Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder){
				require.Equal(t, http.StatusOK, recorder.Code)
				requiredBodyMatchAccount(t, recorder.Body, account)
			},
		},
		{
			name: "NotFound",
			accountId: account.ID,
			buildStubs: func(store *mockdb.MockStore){
				store.EXPECT().
				GetAccount(gomock.Any(), gomock.Eq(account.ID)).
				Times(1).
				Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder){
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "InternalError",
			accountId: account.ID,
			buildStubs: func(store *mockdb.MockStore){
				store.EXPECT().
				GetAccount(gomock.Any(), gomock.Eq(account.ID)).
				Times(1).
				Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder){
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidId",
			accountId: 0,
			buildStubs: func(store *mockdb.MockStore){
				store.EXPECT().
				GetAccount(gomock.Any(), gomock.Any()).
				Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder){
				require.Equal(t, http.StatusBadRequest, recorder.Code)
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

func randomAccount() db.Account{
	return db.Account{
		ID: util.RandomInt(1, 1000),
		Owner: util.RandomOwner(),
		Balance:util.RandomMoney(),
		Currency:util.RandomCurrency(),
	}
}

func requiredBodyMatchAccount(t *testing.T, body *bytes.Buffer, user db.User){
	data , err := io.ReadAll(body)
	require.NoError(t, err)
	
	var gotUser db.User
	err = json.Unmarshal(data, &gotUser)
	require.NoError(t, err)
	require.Equal(t, user, gotUser)


}