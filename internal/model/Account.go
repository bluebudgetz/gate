package model

import (
	"context"
	"database/sql"
	"github.com/pkg/errors"
	"time"
)

const accountsDataManagerKey = "accountsDataManagerKey"

type Account struct {
	ID        int        `json:"id"`
	CreatedOn time.Time  `json:"createdOn"`
	UpdatedOn *time.Time `json:"updatedOn"`
	DeletedOn *time.Time `json:"deletedOn"`
	Name      string     `json:"name"`
	ParentID  *int       `json:"parentId"`
}

func GetAccountsDataManager(ctx context.Context) AccountsDataManager {
	raw, ok := ctx.Value(accountsDataManagerKey).(AccountsDataManager)
	if ok {
		return raw
	} else {
		panic(errors.New("could not find AccountsDataManager in context"))
	}
}

type AccountsDataManager interface {
	LoadAccountByID(ctx context.Context, id int) (*Account, error)
	LoadAllAccounts(ctx context.Context) ([]Account, error)
	LoadRootAccounts(ctx context.Context) ([]Account, error)
	LoadChildAccounts(ctx context.Context, parentId int) ([]Account, error)
	PutInContext(ctx context.Context) context.Context
}

type accountsDataManager struct {
	db *sql.DB
}

func NewAccountsDataManager(db *sql.DB) (AccountsDataManager, error) {
	return &accountsDataManager{db}, nil
}

func (m *accountsDataManager) PutInContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, accountsDataManagerKey, m)
}

func (m *accountsDataManager) loadAccounts(ctx context.Context, sql string, args ...interface{}) ([]Account, error) {
	rows, err := m.db.QueryContext(ctx, sql, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed fetching accounts")
	}
	defer rows.Close()

	accounts := make([]Account, 0, 100)
	for rows.Next() {
		var account Account
		err := m.scanAccount(&account, rows.Scan)
		if err != nil {
			return nil, errors.Wrap(err, "failed scanning account row")
		}
		accounts = append(accounts, account)
	}
	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "failed iterating account rows")
	}
	return accounts, nil
}

func (m *accountsDataManager) scanAccount(account *Account, scanner func(dest ...interface{}) error) error {
	err := scanner(
		&account.ID,
		&account.CreatedOn,
		&account.UpdatedOn,
		&account.DeletedOn,
		&account.Name,
		&account.ParentID)
	if err != nil {
		return errors.Wrap(err, "failed scanning row")
	}
	return nil
}

func (m *accountsDataManager) LoadAccountByID(ctx context.Context, id int) (*Account, error) {
	row := m.db.QueryRowContext(ctx, "SELECT * FROM bb.accounts WHERE id = ?", id)
	var account Account
	err := m.scanAccount(&account, row.Scan)
	if err != nil {
		return nil, errors.Wrapf(err, "failed scanning account %d row", id)
	}
	return &account, nil
}

func (m *accountsDataManager) LoadAllAccounts(ctx context.Context) ([]Account, error) {
	return m.loadAccounts(ctx, "SELECT * FROM accounts WHERE deleted_on IS NULL")
}

func (m *accountsDataManager) LoadRootAccounts(ctx context.Context) ([]Account, error) {
	return m.loadAccounts(ctx, "SELECT * FROM accounts WHERE deleted_on IS NULL AND parent_id IS NULL")
}

func (m *accountsDataManager) LoadChildAccounts(ctx context.Context, parentId int) ([]Account, error) {
	return m.loadAccounts(ctx, "SELECT * FROM accounts WHERE deleted_on IS NULL AND parent_id = ?", parentId)
}
