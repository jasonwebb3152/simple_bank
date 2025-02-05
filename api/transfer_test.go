package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	mockdb "github.com/jasonwebb3152/simplebank/db/mock"
	db "github.com/jasonwebb3152/simplebank/db/sqlc"
	"github.com/jasonwebb3152/simplebank/token"
	"github.com/stretchr/testify/require"
)

func randomAccountWithCurrency(owner string, currency string) db.Account {
	account := randomAccount(owner)
	account.Currency = currency
	return account
}

func TestCreateTransfer(t *testing.T) {
	user1, _ := randomUser(t)
	user2, _ := randomUser(t)
	// user3, _ := randomUser(t)

	account1 := randomAccountWithCurrency(user1.Username, "USD")
	account2 := randomAccountWithCurrency(user2.Username, "USD")
	// account3 := randomAccount(user3.Username)

	testCases := []struct {
		name          string
		fromAccount   db.Account
		toAccount     db.Account
		currency      string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore, fromAccount, toAccount db.Account, amount int64)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:        "OK",
			fromAccount: account1,
			toAccount:   account2,
			currency:    "USD",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, "bearer", user1.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore, fromAccount, toAccount db.Account, amount int64) {
				store.EXPECT().
					TransferMoneyTx(gomock.Any(), gomock.Eq(db.TransferTxParams{
						FromAccountID: fromAccount.ID,
						ToAccountID:   toAccount.ID,
						Amount:        amount,
					})).
					Times(1).
					Return(db.TransferTxResult{}, nil)

				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(fromAccount.ID)).
					Times(1).
					Return(fromAccount, nil)

				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(toAccount.ID)).
					Times(1).
					Return(toAccount, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:        "Unauthorized",
			fromAccount: account1,
			toAccount:   account2,
			currency:    "USD",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore, fromAccount, toAccount db.Account, amount int64) {
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		// TODO: add test cases for errors (needed for wrong request inputs, invalid data...)
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockStore := mockdb.NewMockStore(ctrl)
			tc.buildStubs(mockStore, tc.fromAccount, tc.toAccount, 10)

			server := newTestServer(t, mockStore)
			recorder := httptest.NewRecorder()

			body := gin.H{
				"from_account_id": tc.fromAccount.ID,
				"to_account_id":   tc.toAccount.ID,
				"amount":          10,
				"currency":        tc.currency,
			}
			data, err := json.Marshal(body)
			require.NoError(t, err)

			url := "/transfers"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}
