package apply

import (
	"database/sql"
	"errors"
	"testing"
)

func Test_Apply_NilTransaction(t *testing.T) {
	sqlerr, txerr := Apply(nil, nil)

	if sqlerr != nil {
		t.Error("unexpected error:", sqlerr)
	}

	if txerr == nil {
		t.Error("didn't failed on nil transaction error")
		return
	}

	if txerr.Error() != "nil transaction" {
		t.Error("didn't returned expected error:", txerr)
	}
}

func Test_Apply_SqlError(t *testing.T) {
	tx := MockTx{}
	tx.exec = func(_ string, _ ...interface{}) (sql.Result, error) {
		return nil, errors.New("sql error")
	}
	tx.rollback = func() error {
		return nil
	}

	sqlerr, txerr := Apply([]string{"select 1;"}, tx)

	if txerr != nil {
		t.Error("unexpected error:", txerr)
	}

	if sqlerr == nil {
		t.Error("didn't failed on sql error")
		return
	}

	if sqlerr.Error() != "sql error" {
		t.Error("didn't returned expected error:", sqlerr)
	}
}

func Test_Apply_RollbackError(t *testing.T) {
	tx := MockTx{}
	tx.exec = func(_ string, _ ...interface{}) (sql.Result, error) {
		return nil, errors.New("sql error")
	}
	tx.rollback = func() error {
		return errors.New("rollback error")
	}

	sqlerr, txerr := Apply([]string{"select 1;"}, tx)

	if sqlerr == nil {
		t.Error("didn't failed on sql error")
		return
	}

	if sqlerr.Error() != "sql error" {
		t.Error("didn't returned expected error:", sqlerr)
	}

	if txerr == nil {
		t.Error("didn't failed on rollback error")
		return
	}

	if txerr.Error() != "rollback error" {
		t.Error("didn't returned expected error:", sqlerr)
	}
}

func Test_Apply_CommitError(t *testing.T) {
	tx := MockTx{}
	tx.exec = func(_ string, _ ...interface{}) (sql.Result, error) {
		return nil, nil
	}
	tx.commit = func() error {
		return errors.New("commit error")
	}

	sqlerr, txerr := Apply([]string{"select 1;"}, tx)

	if sqlerr != nil {
		t.Error("unexpected error:", sqlerr)
	}

	if txerr == nil {
		t.Error("didn't failed on commit error")
		return
	}

	if txerr.Error() != "commit error" {
		t.Error("didn't returned expected error:", sqlerr)
	}
}
