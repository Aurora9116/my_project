package dao

import (
	"errors"
	"test.com/project-common/errs"
	"test.com/project-project/internal/database"
	"test.com/project-project/internal/database/gorms"
)

type TransactionImpl struct {
	conn *gorms.GormConn
}

func (t *TransactionImpl) Action(f func(conn database.DbConn) error) error {
	t.conn.Begin()
	err := f(t.conn)
	var bErr *errs.BError
	if errors.Is(err, bErr) {
		bErr = err.(*errs.BError)
		if bErr != nil {
			t.conn.RollBack()
			return bErr
		} else {
			t.conn.Commit()
			return nil
		}
	}
	if err != nil {
		t.conn.RollBack()
		return err
	}
	t.conn.Commit()
	return nil
}

func NewTransactionImpl() *TransactionImpl {
	return &TransactionImpl{
		conn: gorms.NewTx(),
	}
}
