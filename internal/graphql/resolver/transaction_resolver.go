package resolver

import (
	"context"
	"github.com/bluebudgetz/gate/internal/middleware"
	"github.com/bluebudgetz/gate/internal/model"
	"github.com/pkg/errors"
)

type transactionResolver struct{ *Resolver }

func (r *transactionResolver) transaction(ctx context.Context, id int) (*model.Transaction, error) {
	_sql := `
	SELECT
       	tx.id,
	   	tx.created_on, 
	   	tx.updated_on, 
	   	tx.deleted_on, 
	   	tx.origin, 
	   	tx.source_account_id,
	   	tx.target_account_id,
	   	tx.amount,
	   	tx.comments
	FROM
    	bb.transactions AS tx 
	WHERE
        tx.deleted_on IS NULL AND tx.id = ? `
	row := middleware.GetDB(ctx).QueryRowContext(ctx, _sql, id)

	var tx model.Transaction
	err := row.Scan(&tx.ID, &tx.CreatedOn, &tx.UpdatedOn, &tx.DeletedOn, &tx.Origin, &tx.SourceID, &tx.TargetID, &tx.Amount, &tx.Comments)
	if err != nil {
		return nil, errors.Wrap(err, "failed scanning transaction row")
	}
	return &tx, nil
}

func (r *transactionResolver) transactions(ctx context.Context, where string, args ...interface{}) ([]model.Transaction, error) {
	_sql := `
	SELECT
       	tx.id,
	   	tx.created_on, 
	   	tx.updated_on, 
	   	tx.deleted_on, 
	   	tx.origin, 
	   	tx.source_account_id,
	   	tx.target_account_id,
	   	tx.amount,
	   	tx.comments
	FROM
    	bb.transactions AS tx 
	WHERE
        tx.deleted_on IS NULL `
	if where != "" {
		_sql = _sql + " AND " + where
	}
	rows, err := middleware.GetDB(ctx).QueryContext(ctx, _sql, args...)
	if err != nil {
		return nil, errors.Wrapf(err, "failed fetching transactions")
	}
	defer rows.Close()

	transactions := make([]model.Transaction, 0, 100)
	for rows.Next() {
		var tx model.Transaction
		err := rows.Scan(&tx.ID, &tx.CreatedOn, &tx.UpdatedOn, &tx.DeletedOn, &tx.Origin, &tx.SourceID, &tx.TargetID, &tx.Amount, &tx.Comments)
		if err != nil {
			return nil, errors.Wrap(err, "failed scanning transaction row")
		}
		transactions = append(transactions, tx)
	}
	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "failed iterating transaction rows")
	}
	return transactions, nil
}

func (r *transactionResolver) Source(ctx context.Context, obj *model.Transaction) (*model.Account, error) {
	return r.Query().Account(ctx, obj.SourceID)
}

func (r *transactionResolver) Target(ctx context.Context, obj *model.Transaction) (*model.Account, error) {
	return r.Query().Account(ctx, obj.TargetID)
}
