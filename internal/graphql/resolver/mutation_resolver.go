package resolver

import (
	"context"
	"github.com/bluebudgetz/gate/internal/graphql/model"
	"github.com/bluebudgetz/gate/internal/middleware"
	"github.com/pkg/errors"
)

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) CreateAccount(ctx context.Context, name string, parentID *int) (*model.Account, error) {
	var id int
	err := middleware.GetDB(ctx).QueryRowContext(
		ctx,
		"INSERT INTO bb.accounts (name, parent_id) VALUES ($1, $2) RETURNING id",
		name, parentID).Scan(&id)
	if err != nil {
		return nil, errors.Wrapf(err, "failed creating account '%s'", name)
	}
	return r.Query().Account(ctx, int(id))
}

func (r *mutationResolver) CreateTransaction(ctx context.Context, origin string, sourceAccountId int, targetAccountId int, amount float64, comments *string) (*model.Transaction, error) {

	db := middleware.GetDB(ctx)
	tx, err := db.Begin()
	if err != nil {
		return nil, errors.Wrap(err, "failed starting transaction")
	}
	defer tx.Rollback()

	var id int
	err = tx.QueryRowContext(
		ctx,
		"INSERT INTO bb.transactions (origin, source_account_id, target_account_id, amount, comments) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		origin, sourceAccountId, targetAccountId, amount, comments).Scan(&id)
	if err != nil {
		return nil, errors.Wrap(err, "failed creating transaction")
	}

	err = tx.Commit()
	if err != nil {
		return nil, errors.Wrap(err, "failed committing transaction")
	}
	return r.Query().Transaction(ctx, int(id))
}
