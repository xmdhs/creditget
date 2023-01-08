package db

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	"github.com/jmoiron/sqlx"
	"github.com/xmdhs/creditget/model"
)

type DB interface {
	BatchInsterCreditInfo(cxt context.Context, tx *sqlx.Tx, c []model.CreditInfo) error
	GetCreditInfo(cxt context.Context, uid int) (*model.CreditInfo, error)
	Begin(ctx context.Context, opts *sql.TxOptions) (*Tx, error)
	GetRank(cxt context.Context, uid int, field string) (int, error)
	InsterConfig(cxt context.Context, tx *sqlx.Tx, c *model.Confing) error
	SelectConfig(cxt context.Context, id int) (*model.Confing, error)
}

type Tx struct {
	tx          *sqlx.Tx
	afterCommit []func()
	m           *sync.Mutex
}

func NewTx(tx *sqlx.Tx) *Tx {
	return &Tx{
		tx:          tx,
		afterCommit: []func(){},
		m:           &sync.Mutex{},
	}
}

func GetCxtTx(cxt context.Context) *Tx {
	return cxt.Value(txCxtKey{}).(*Tx)
}

func (t *Tx) AddAfterCommit(f func()) {
	t.m.Lock()
	defer t.m.Unlock()
	t.afterCommit = append(t.afterCommit, f)
}

type txCxtKey struct{}

func (t *Tx) Transaction(cxt context.Context, f func(cxt context.Context, tx *sqlx.Tx) error) (err error) {
	defer func() {
		if err == nil {
			return
		}
		if rollbackErr := t.tx.Rollback(); rollbackErr != nil {
			err = fmt.Errorf("Transaction: failed to rollback (%s) %w", rollbackErr.Error(), err)
		}
	}()
	cxt = context.WithValue(cxt, txCxtKey{}, t)
	err = f(cxt, t.tx)
	if err != nil {
		return fmt.Errorf("Transaction: failed to execute transaction: %w", err)
	}
	err = t.tx.Commit()
	if err != nil {
		return fmt.Errorf("Transaction: failed to commit transaction: %w", err)
	}
	t.m.Lock()
	defer t.m.Unlock()
	for _, v := range t.afterCommit {
		v()
	}
	return nil
}
