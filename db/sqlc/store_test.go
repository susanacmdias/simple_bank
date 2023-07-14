package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	fmt.Println(">> before:", account1.Balance, account2.Balance)

	//run n concurrent transfer transactions
	n := 5
	ammount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {

		go func() {
			result, err := store.TranferTx(context.Background(), TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Ammount:       ammount,
			})

			errs <- err
			results <- result

		}()
	}

	//check results
	existed := make(map[int]bool)
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		//check transfer
		transfer := result.Tranfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, ammount, transfer.Ammount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.getTranfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// check entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -ammount, fromEntry.Ammount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.getEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, ammount, toEntry.Ammount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.getEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		//check accounts
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		ToAccount := result.ToAccount
		require.NotEmpty(t, ToAccount)
		require.Equal(t, account2.ID, ToAccount.ID)

		//check baccounts ballance
		fmt.Println(">> tx:", fromAccount.Balance, ToAccount.Balance)
		diff1 := account1.Balance - fromAccount.Balance
		diff2 := ToAccount.Balance - account2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%ammount == 0)

		k := int(diff1 / ammount)
		require.True(t, k >= 1 && k <= n)

		require.NotContains(t, existed, k)
		existed[k] = true
	}
	//check the final updated balances
	updateAccount1, err := testQueries.getAccount(context.Background(), account1.ID)
	require.NotEmpty(t, updateAccount1)
	require.NoError(t, err)

	updateAccount2, err := testQueries.getAccount(context.Background(), account2.ID)
	require.NotEmpty(t, updateAccount2)
	require.NoError(t, err)

	fmt.Println(">> after:", updateAccount1.Balance, updateAccount2.Balance)
	require.Equal(t, account1.Balance-int64(n)*ammount, updateAccount1.Balance)
	require.Equal(t, account2.Balance+int64(n)*ammount, updateAccount2.Balance)
}

func TestTransferTxDeadlock(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	fmt.Println(">> before:", account1.Balance, account2.Balance)

	//run n concurrent transfer transactions
	n := 10
	ammount := int64(10)

	errs := make(chan error)
	for i := 0; i < 10; i++ {
		fromAccountID := account1.ID
		toAccountID := account2.ID

		if i%2 == 1 {
			fromAccountID = account2.ID
			toAccountID = account1.ID
		}

		go func() {
			_, err := store.TranferTx(context.Background(), TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Ammount:       ammount,
			})

			errs <- err

		}()
	}

	//check results
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}
	//check the final updated balances
	updateAccount1, err := testQueries.getAccount(context.Background(), account1.ID)
	require.NotEmpty(t, updateAccount1)
	require.NoError(t, err)

	updateAccount2, err := testQueries.getAccount(context.Background(), account2.ID)
	require.NotEmpty(t, updateAccount2)
	require.NoError(t, err)

	fmt.Println(">> after:", updateAccount1.Balance, updateAccount2.Balance)
	require.Equal(t, account1.Balance, updateAccount1.Balance)
	require.Equal(t, account2.Balance, updateAccount2.Balance)
}
