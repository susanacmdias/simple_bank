package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/techschool/simplebank/utils"
)

func createRandomEntry(t *testing.T, account Account) Entry {
	arg := createEntryParams{
		AccountID: account.ID,
		Ammount:   utils.RandomMoney(),
	}

	entry, err := testQueries.createEntry(context.Background(), arg)
	require.NoError(t, err)

	require.NotEmpty(t, entry)

	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Ammount, entry.Ammount)

	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

	return entry
}

func TestCreateEntry(t *testing.T) {
	account := createRandomAccount(t)
	createRandomEntry(t, account)
}

func TestGetEntry(t *testing.T) {
	account := createRandomAccount(t)
	entry1 := createRandomEntry(t, account)

	entry2, err := testQueries.getEntry(context.Background(), entry1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, entry2)

	require.Equal(t, entry2.ID, entry1.ID)
	require.Equal(t, entry2.AccountID, entry1.AccountID)
	require.Equal(t, entry2.Ammount, entry1.Ammount)
	require.WithinDuration(t, entry2.CreatedAt, entry1.CreatedAt, time.Second)

}

func TestListEntry(t *testing.T) {
	account := createRandomAccount(t)
	for i := 0; i < 10; i++ {
		createRandomEntry(t, account)
	}

	arg := listEntriesParams{
		AccountID: account.ID,
		Limit:     5,
		Offset:    5,
	}

	entries, err := testQueries.listEntries(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, entries, 5)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
		require.Equal(t, arg.AccountID, entry.AccountID)
	}

}
