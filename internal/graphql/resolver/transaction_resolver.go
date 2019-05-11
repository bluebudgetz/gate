package resolver

import (
	"context"
	"github.com/bluebudgetz/gate/internal/graphql/model"
)

type transactionResolver struct{ *Resolver }

func (r *transactionResolver) Source(ctx context.Context, obj *model.Transaction) (*model.Account, error) {
	return r.Query().Account(ctx, obj.SourceID)
}

func (r *transactionResolver) Target(ctx context.Context, obj *model.Transaction) (*model.Account, error) {
	return r.Query().Account(ctx, obj.TargetID)
}
