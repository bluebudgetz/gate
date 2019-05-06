package resolver

import (
	"context"
	"github.com/bluebudgetz/gate/internal/middleware"
	"github.com/bluebudgetz/gate/internal/model"
	"github.com/pkg/errors"
)

type accountResolver struct{ *Resolver }

func (r *accountResolver) account(ctx context.Context, id int) (*model.Account, error) {
	_sql := `
	SELECT
       acc.id,
	   acc.created_on, 
	   acc.updated_on, 
	   acc.deleted_on, 
	   acc.name, 
	   acc.parent_id,
	   IFNULL((SELECT SUM(tx.amount) FROM bb.transactions AS tx WHERE tx.target_account_id = acc.id), 0) -
	   IFNULL((SELECT SUM(tx.amount) FROM bb.transactions AS tx WHERE tx.source_account_id = acc.id), 0) AS balance
	FROM
    	bb.accounts AS acc
	WHERE
        acc.deleted_on IS NULL AND acc.id = ? `
	row := middleware.GetDB(ctx).QueryRowContext(ctx, _sql, id)

	var account model.Account
	err := row.Scan(&account.ID, &account.CreatedOn, &account.UpdatedOn, &account.DeletedOn, &account.Name, &account.ParentID, &account.Balance)
	if err != nil {
		return nil, errors.Wrapf(err, "failed scanning account %d", id)
	}
	return &account, nil
}

func (r *accountResolver) accounts(ctx context.Context, where string, args ...interface{}) ([]model.Account, error) {
	_sql := `
	SELECT
       acc.id,
	   acc.created_on, 
	   acc.updated_on, 
	   acc.deleted_on, 
	   acc.name, 
	   acc.parent_id,
	   IFNULL((SELECT SUM(tx.amount) FROM bb.transactions AS tx WHERE tx.target_account_id = acc.id), 0) -
	   IFNULL((SELECT SUM(tx.amount) FROM bb.transactions AS tx WHERE tx.source_account_id = acc.id), 0) AS balance
	FROM
    	bb.accounts AS acc
	WHERE
        acc.deleted_on IS NULL `
	if where != "" {
		_sql = _sql + " AND " + where
	}
	rows, err := middleware.GetDB(ctx).QueryContext(ctx, _sql, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed fetching accounts")
	}
	defer rows.Close()

	accounts := make([]model.Account, 0, 100)
	for rows.Next() {
		var account model.Account
		err := rows.Scan(&account.ID, &account.CreatedOn, &account.UpdatedOn, &account.DeletedOn, &account.Name, &account.ParentID, &account.Balance)
		if err != nil {
			return nil, errors.Wrapf(err, "failed scanning account")
		}
		accounts = append(accounts, account)
	}
	if err := rows.Err(); err != nil {
		return nil, errors.Wrapf(err, "failed iterating account rows")
	}
	return accounts, nil
}

func (r *accountResolver) Parent(ctx context.Context, obj *model.Account) (*model.Account, error) {
	if obj.ParentID == nil {
		return nil, nil
	} else {
		return r.account(ctx, *obj.ParentID)
	}
}

func (r *accountResolver) ChildAccounts(ctx context.Context, obj *model.Account) ([]model.Account, error) {
	return r.accounts(ctx, "acc.parent_id = ?", obj.ID)
}
