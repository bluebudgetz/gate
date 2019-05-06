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
