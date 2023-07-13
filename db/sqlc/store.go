package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Provides all functions to execute db queries and transactions
type Store struct {
	*Queries
	db *sql.DB
}

// New store creates a new Store
func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

// execTx executes a function within a database transaction
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx error %v, rb error %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()

}

type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Ammount       int64 `json:"ammount"`
}

type TransferTxResult struct {
	Tranfer     Tranfer `json:"transfer"`
	FromAccount Account `json:"from_account"`
	ToAccount   Account `json:"to_account"`
	FromEntry   Entry   `json:"from_entry"`
	ToEntry     Entry   `json:"to_entry"`
}

// TranferTx performs money tranfers from one account to other
// Creates a transfer record, add account entries, and update account balance within a single database
func (store *Store) TranferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult
	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.Tranfer, err = q.createTranfer(ctx, createTranferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Ammount:       arg.Ammount,
		})
		if err != nil {
			return err
		}

		result.FromEntry, err = q.createEntry(ctx, createEntryParams{
			AccountID: arg.FromAccountID,
			Ammount:   -arg.Ammount,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.createEntry(ctx, createEntryParams{
			AccountID: arg.ToAccountID,
			Ammount:   arg.Ammount,
		})
		if err != nil {
			return err
		}

		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, result.FromAccount, err = addMoney(ctx, q, arg.FromAccountID, -arg.Ammount, arg.ToAccountID, arg.Ammount)
		} else {
			result.ToAccount, result.FromAccount, err = addMoney(ctx, q, arg.ToAccountID, arg.Ammount, arg.FromAccountID, -arg.Ammount)
		}

		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}

func addMoney(
	ctx context.Context,
	q *Queries,
	accountID1 int64,
	ammount1 int64,
	accountID2 int64,
	ammount2 int64,
) (account1 Account, account2 Account, err error) {

	account1, err = q.addAccountBalance(ctx, addAccountBalanceParams{
		ID:      accountID1,
		Ammount: ammount1,
	})
	if err != nil {
		return
	}

	account2, err = q.addAccountBalance(ctx, addAccountBalanceParams{
		ID:      accountID2,
		Ammount: ammount2,
	})
	return
}
