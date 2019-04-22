package internal

import (
	"context"
	"database/sql"
	"github.com/bluebudgetz/gate/internal/graphql"
	"github.com/bluebudgetz/gate/internal/model"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

func NewResolverRoot(db *sql.DB) (graphql.ResolverRoot, error) {
	return &resolver{db}, nil
}

type resolver struct {
	db *sql.DB
}

func (r *resolver) Account() graphql.AccountResolver {
	return &accountResolver{r}
}

func (r *resolver) Query() graphql.QueryResolver {
	return &queryResolver{r}
}

type queryResolver struct{ *resolver }

func (r *queryResolver) Accounts(ctx context.Context, limit *int, offset *int) ([]model.Account, error) {
	accounts, err := model.GetAccountsDataManager(ctx).LoadRootAccounts(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed fetching accounts")
	}
	return accounts, nil
}

type accountResolver struct{ *resolver }

func (r *accountResolver) Parent(ctx context.Context, obj *model.Account) (*model.Account, error) {
	if obj.ParentID == nil {
		return nil, nil
	}

	account, err := model.GetAccountsDataManager(ctx).LoadAccountByID(ctx, *obj.ParentID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed fetching parent account %d", *obj.ParentID)
	}
	return account, nil
}

func (r *accountResolver) ChildAccounts(ctx context.Context, obj *model.Account) ([]model.Account, error) {
	accounts, err := model.GetAccountsDataManager(ctx).LoadChildAccounts(ctx, obj.ID)
	if err != nil {
		return nil, errors.Wrap(err, "failed fetching accounts")
	}
	return accounts, nil
}
