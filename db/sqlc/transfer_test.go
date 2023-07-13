package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/techschool/simplebank/utils"
)

func createRandomTransfer(t *testing.T, account1 Account, account2 Account) Tranfer {
	arg := createTranferParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Ammount:       utils.RandomMoney(),
	}

	tranfer, err := testQueries.createTranfer(context.Background(), arg)
	require.NoError(t, err)

	require.Equal(t, arg.FromAccountID, tranfer.FromAccountID)
	require.Equal(t, arg.ToAccountID, tranfer.ToAccountID)
	require.Equal(t, arg.Ammount, tranfer.Ammount)

	require.NotZero(t, tranfer.ID)
	require.NotZero(t, tranfer.CreatedAt)

	return tranfer
}

func TestCreateTransfer(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	createRandomTransfer(t, account1, account2)
}

func TestGetTransfer(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	tranfer1 := createRandomTransfer(t, account1, account2)

	tranfer, err := testQueries.getTranfer(context.Background(), tranfer1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, tranfer)

	require.Equal(t, tranfer.ID, tranfer1.ID)
	require.Equal(t, tranfer.FromAccountID, tranfer1.FromAccountID)
	require.Equal(t, tranfer.ToAccountID, tranfer1.ToAccountID)
	require.Equal(t, tranfer.FromAccountID, tranfer1.FromAccountID)
	require.Equal(t, tranfer.Ammount, tranfer1.Ammount)
	require.Equal(t, tranfer.CreatedAt, tranfer1.CreatedAt)

}

func TestListTransfer(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	arg := listTranferParams{
		FromAccountID: account1.ID,
		ToAccountID:   account1.ID,
		Offset:        5,
		Limit:         5,
	}

	for i := 0; i < 10; i++ {
		createRandomTransfer(t, account1, account2)
		createRandomTransfer(t, account2, account1)
	}

	tranfer, err := testQueries.listTranfer(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, tranfer, 5)

	for _, tranfers := range tranfer {
		require.NotEmpty(t, tranfers)
		require.True(t, tranfers.FromAccountID == account1.ID || tranfers.ToAccountID == account1.ID)
	}

}
