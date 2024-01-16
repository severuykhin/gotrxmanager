package gotrxmanager

import (
	"context"
	"database/sql"
	"fmt"
)

type trxManagerKey string

const trxKey trxManagerKey = "trxKey"

type transactionManager struct {
	db *sql.DB
}

func NewTransactionManager(db *sql.DB) *transactionManager {
	return &transactionManager{
		db: db,
	}
}

func (trm *transactionManager) Do(ctx context.Context, f func(ctx context.Context) (any, error)) (any, error) {

	trx, err := trm.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	ctx = context.WithValue(ctx, trxKey, trx)

	res, err := f(ctx)
	if err != nil {
		if rbErr := trx.Rollback(); rbErr != nil {
			err = fmt.Errorf("cannot rollback transaction with err: %s prev error: %s", rbErr, err)
		}
		return nil, err
	}

	if err := trx.Commit(); err != nil {
		return nil, fmt.Errorf("cannot commit transaction with error: %s", err)
	}

	return res, nil
}

func TxFromContext(ctx context.Context) (*sql.Tx, error) {

	t := ctx.Value(trxKey)
	if t == nil {
		return nil, fmt.Errorf("cannot find transaction")
	}

	tx, ok := t.(*sql.Tx)
	if !ok {
		return nil, fmt.Errorf("received value is not a *sql.Tx")
	}

	return tx, nil
}
