package resolver

import (
	"context"
	"github.com/bluebudgetz/gate/internal/middleware"
	"github.com/bluebudgetz/gate/internal/model"
	"github.com/pkg/errors"
)

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) CreateAccount(ctx context.Context, name string, parentID *int) (*model.Account, error) {
	result, err := middleware.GetDB(ctx).ExecContext(ctx, "INSERT INTO bb.accounts (name, parent_id) VALUES (?, ?)", name, parentID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed creating account '%s' (parent %d)", name, parentID)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, errors.Wrap(err, "failed fetching newly-created account's ID")
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

	result, err := tx.ExecContext(
		ctx,
		"INSERT INTO bb.transactions (origin, source_account_id, target_account_id, amount, comments) VALUES (?, ?, ?, ?, ?)",
		origin, sourceAccountId, targetAccountId, amount, comments)
	if err != nil {
		return nil, errors.Wrap(err, "failed creating transaction")
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, errors.Wrap(err, "failed fetching newly-created transaction's ID")
	}

	accountId := &targetAccountId
	for accountId != nil {
		result, err = tx.ExecContext(
			ctx,
			"UPDATE bb.accounts SET balance = balance + ? WHERE id = ?",
			amount, *accountId)
		if err != nil {
			return nil, errors.Wrapf(err, "failed adding %d to balance of account %d (tx %d)", amount, *accountId, id)
		}

		account, err := r.Account().(*accountResolver).account(ctx, *accountId)
		if err != nil {
			return nil, errors.Wrapf(err, "failed looking up account %d (tx %d)", *accountId, id)
		}

		accountId = account.ParentID
	}

	accountId = &sourceAccountId
	for accountId != nil {
		result, err = tx.ExecContext(
			ctx,
			"UPDATE bb.accounts SET balance = balance - ? WHERE id = ?",
			amount, *accountId)
		if err != nil {
			return nil, errors.Wrapf(err, "failed subtracting %d from balance of account %d (tx %d)", amount, accountId, id)
		}

		account, err := r.Account().(*accountResolver).account(ctx, *accountId)
		if err != nil {
			return nil, errors.Wrapf(err, "failed looking up account %d (tx %d)", *accountId, id)
		}

		accountId = account.ParentID
	}

	err = tx.Commit()
	if err != nil {
		return nil, errors.Wrap(err, "failed committing transaction")
	}
	return r.Query().Transaction(ctx, int(id))
}
