package resolver

import (
	"context"
	"github.com/bluebudgetz/gate/internal/model"
	"time"
)

type queryResolver struct{ *Resolver }

func (r *queryResolver) Account(ctx context.Context, id int) (*model.Account, error) {
	return r.Resolver.Account().(*accountResolver).account(ctx, id)
}

func (r *queryResolver) Accounts(ctx context.Context) ([]model.Account, error) {
	return r.Resolver.Account().(*accountResolver).accounts(ctx, "acc.parent_id IS NULL")
}

func (r *queryResolver) Transactions(ctx context.Context, from time.Time, until time.Time) ([]model.Transaction, error) {
	return r.Resolver.Transaction().(*transactionResolver).transactions(ctx, "created_on >= ? AND created_on <= ?", from, until)
}
