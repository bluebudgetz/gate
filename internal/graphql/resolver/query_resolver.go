package resolver

import (
	"context"
	"github.com/bluebudgetz/gate/internal/graphql/model"
	"github.com/bluebudgetz/gate/internal/middleware"
	"github.com/pkg/errors"
	"time"
)

type queryResolver struct{ *Resolver }

func (r *queryResolver) Account(ctx context.Context, id int) (*model.Account, error) {
	_sql := `
		WITH RECURSIVE cte (id, child_id, created_on, updated_on, name) AS (
	        SELECT acc.id, acc.id, acc.created_on, acc.updated_on, acc.name
	        FROM bb.accounts AS acc
	        WHERE acc.id = $1
	        UNION ALL
	        SELECT h.id, acc.id, h.created_on, h.updated_on, h.name
	        FROM bb.accounts AS acc
	            INNER JOIN cte AS h ON acc.parent_id = h.child_id)
		SELECT
		    cte.id                                                                                          AS id,
       		cte.created_on                                                                                  AS created_on,
       		cte.updated_on                                                                                  AS updated_on,
       		cte.name                                                                                        AS name,
        	COALESCE((SELECT COUNT(*) FROM bb.accounts AS children WHERE children.parent_id = cte.id), 0)   AS child_count,
        	COALESCE(SUM(outgoing_tx.amount), 0)                                                            AS outgoing,
        	COALESCE(SUM(incoming_tx.amount), 0)                                                            AS incoming
		FROM cte
         	LEFT JOIN bb.transactions AS outgoing_tx ON outgoing_tx.source_account_id = cte.child_id
         	LEFT JOIN bb.transactions AS incoming_tx ON incoming_tx.target_account_id = cte.child_id
		GROUP BY 
			cte.id, 
		    cte.created_on, 
	        cte.updated_on,
	        cte.name`
	var acc model.Account
	err := middleware.GetDB(ctx).QueryRowContext(ctx, _sql, id).Scan(&acc.ID, &acc.CreatedOn, &acc.UpdatedOn, &acc.Name, &acc.ChildCount, &acc.Outgoing, &acc.Incoming)
	if err != nil {
		return nil, errors.Wrapf(err, "failed fetching account %d", id)
	}
	return &acc, nil
}

func (r *queryResolver) RootAccounts(ctx context.Context) ([]model.Account, error) {
	_sql := `
		WITH RECURSIVE cte (id, child_id, created_on, updated_on, name) AS (
	        SELECT acc.id, acc.id, acc.created_on, acc.updated_on, acc.name
	        FROM bb.accounts AS acc
	        WHERE acc.parent_id IS NULL
	        UNION ALL
	        SELECT h.id, acc.id, h.created_on, h.updated_on, h.name
	        FROM bb.accounts AS acc
	            INNER JOIN cte AS h ON acc.parent_id = h.child_id)
		SELECT
		    cte.id                                                                                          AS id,
       		cte.created_on                                                                                  AS created_on,
       		cte.updated_on                                                                                  AS updated_on,
       		cte.name                                                                                        AS name,
        	COALESCE((SELECT COUNT(*) FROM bb.accounts AS children WHERE children.parent_id = cte.id), 0)   AS child_count,
        	COALESCE(SUM(outgoing_tx.amount), 0)                                                            AS outgoing,
        	COALESCE(SUM(incoming_tx.amount), 0)                                                            AS incoming
		FROM cte
         	LEFT JOIN bb.transactions AS outgoing_tx ON outgoing_tx.source_account_id = cte.child_id
         	LEFT JOIN bb.transactions AS incoming_tx ON incoming_tx.target_account_id = cte.child_id
		GROUP BY 
			cte.id, 
		    cte.created_on, 
	        cte.updated_on,
	        cte.name`
	rows, err := middleware.GetDB(ctx).QueryContext(ctx, _sql)
	if err != nil {
		return nil, errors.Wrap(err, "failed fetching root accounts")
	}
	defer rows.Close()

	accounts := make([]model.Account, 0, 100)
	for rows.Next() {
		var acc model.Account
		err := rows.Scan(&acc.ID, &acc.CreatedOn, &acc.UpdatedOn, &acc.Name, &acc.ChildCount, &acc.Outgoing, &acc.Incoming)
		if err != nil {
			return nil, errors.Wrapf(err, "failed scanning account")
		}
		accounts = append(accounts, acc)
	}
	if err := rows.Err(); err != nil {
		return nil, errors.Wrapf(err, "failed iterating root account rows")
	}
	return accounts, nil
}

func (r *queryResolver) ChildAccounts(ctx context.Context, parentID int) ([]model.Account, error) {
	_sql := `
		WITH RECURSIVE cte (id, child_id, created_on, updated_on, name) AS (
	        SELECT acc.id, acc.id, acc.created_on, acc.updated_on, acc.name
	        FROM bb.accounts AS acc
	        WHERE acc.parent_id = $1
	        UNION ALL
	        SELECT h.id, acc.id, h.created_on, h.updated_on, h.name
	        FROM bb.accounts AS acc
	            INNER JOIN cte AS h ON acc.parent_id = h.child_id)
		SELECT
		    cte.id                                                                                          AS id,
       		cte.created_on                                                                                  AS created_on,
       		cte.updated_on                                                                                  AS updated_on,
       		cte.name                                                                                        AS name,
        	COALESCE((SELECT COUNT(*) FROM bb.accounts AS children WHERE children.parent_id = cte.id), 0)   AS child_count,
        	COALESCE(SUM(outgoing_tx.amount), 0)                                                            AS outgoing,
        	COALESCE(SUM(incoming_tx.amount), 0)                                                            AS incoming
		FROM cte
         	LEFT JOIN bb.transactions AS outgoing_tx ON outgoing_tx.source_account_id = cte.child_id
         	LEFT JOIN bb.transactions AS incoming_tx ON incoming_tx.target_account_id = cte.child_id
		GROUP BY 
			cte.id, 
		    cte.created_on, 
	        cte.updated_on,
	        cte.name`
	rows, err := middleware.GetDB(ctx).QueryContext(ctx, _sql, parentID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed fetching child accounts of parent %d", parentID)
	}
	defer rows.Close()

	accounts := make([]model.Account, 0, 100)
	for rows.Next() {
		var acc model.Account
		err := rows.Scan(&acc.ID, &acc.CreatedOn, &acc.UpdatedOn, &acc.Name, &acc.ChildCount, &acc.Outgoing, &acc.Incoming)
		if err != nil {
			return nil, errors.Wrapf(err, "failed scanning a child account of parent %d", parentID)
		}
		accounts = append(accounts, acc)
	}
	if err := rows.Err(); err != nil {
		return nil, errors.Wrapf(err, "failed iterating child account rows of parent %d", parentID)
	}
	return accounts, nil
}

func (r *queryResolver) Transaction(ctx context.Context, id int) (*model.Transaction, error) {
	_sql := `
	SELECT
       	tx.id,
	   	tx.created_on, 
	   	tx.updated_on, 
	   	tx.origin, 
	   	tx.source_account_id,
	   	tx.target_account_id,
	   	tx.amount,
	   	tx.comments
	FROM
    	bb.transactions AS tx 
	WHERE
        tx.id = $1 `
	var tx model.Transaction
	err := middleware.GetDB(ctx).QueryRowContext(ctx, _sql, id).
		Scan(&tx.ID, &tx.CreatedOn, &tx.UpdatedOn, &tx.Origin, &tx.SourceID, &tx.TargetID, &tx.Amount, &tx.Comments)
	if err != nil {
		return nil, errors.Wrap(err, "failed scanning transaction row")
	}
	return &tx, nil
}

func (r *queryResolver) Transactions(ctx context.Context, from time.Time, until time.Time) ([]model.Transaction, error) {
	// TODO: fetch source/target accounts via JOIN to save two additional queries
	_sql := `
	SELECT
       	tx.id,
	   	tx.created_on, 
	   	tx.updated_on, 
	   	tx.origin, 
	   	tx.source_account_id,
	   	tx.target_account_id,
	   	tx.amount,
	   	tx.comments
	FROM
    	bb.transactions AS tx 
    WHERE 
        created_on >= $1 AND created_on <= $2 `
	rows, err := middleware.GetDB(ctx).QueryContext(ctx, _sql, from, until)
	if err != nil {
		return nil, errors.Wrapf(err, "failed fetching transactions")
	}
	defer rows.Close()

	transactions := make([]model.Transaction, 0, 100)
	for rows.Next() {
		var tx model.Transaction
		err := rows.Scan(&tx.ID, &tx.CreatedOn, &tx.UpdatedOn, &tx.Origin, &tx.SourceID, &tx.TargetID, &tx.Amount, &tx.Comments)
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
