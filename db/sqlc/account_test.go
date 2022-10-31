package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/techschool/simplebank/util"
)

func createRandomAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	account := createRandomAccount(t)
	account_stored, err := testQueries.GetAccount(context.Background(), account.ID)

	require.NoError(t, err)
	require.NotEmpty(t, account_stored)
	require.Equal(t, account_stored.ID, account.ID)
	require.Equal(t, account_stored.Owner, account.Owner)
	require.Equal(t, account_stored.Currency, account.Currency)
	require.Equal(t, account_stored.Balance, account.Balance)

	require.WithinDuration(t, account.CreatedAt, account_stored.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	account := createRandomAccount(t)

	args := UpdateAccountParams{
		ID:      account.ID,
		Balance: util.RandomMoney(),
	}

	account_stored, err := testQueries.UpdateAccount(context.Background(), args)

	require.NoError(t, err)
	require.NotEmpty(t, account_stored)
	require.Equal(t, account_stored.ID, account.ID)
	require.Equal(t, account_stored.Owner, account.Owner)
	require.Equal(t, account_stored.Currency, account.Currency)
	require.Equal(t, account_stored.Balance, args.Balance)
	require.Equal(t, account.CreatedAt, account_stored.CreatedAt)
}

func TestDeleteAccount(t *testing.T) {
	account := createRandomAccount(t)

	err := testQueries.DeleteAccount(context.Background(), account.ID)
	require.NoError(t, err)

	_, err = testQueries.GetAccount(context.Background(), account.ID)
	require.Error(t, err, sql.ErrNoRows)
}

func TestListAccounts(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}

	args := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), args)
	require.NoError(t, err)
	require.Len(t, accounts, 5)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}
