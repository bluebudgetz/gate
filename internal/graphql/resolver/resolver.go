package resolver

import (
	"context"
	"github.com/99designs/gqlgen/graphql"
	"github.com/bluebudgetz/gate/internal/graphql/impl"
	"github.com/vektah/gqlparser/ast"
)

type Resolver struct{}

func (r *Resolver) Account() impl.AccountResolver {
	return &accountResolver{r}
}

func (r *Resolver) Mutation() impl.MutationResolver {
	return &mutationResolver{r}
}

func (r *Resolver) Query() impl.QueryResolver {
	return &queryResolver{r}
}
func (r *Resolver) Transaction() impl.TransactionResolver {
	return &transactionResolver{r}
}

func hasProperty(ctx context.Context, selectionSet ast.SelectionSet, property string) bool {
	reqCtx := graphql.GetRequestContext(ctx)
	for _, sel := range selectionSet {
		switch sel := sel.(type) {
		case *ast.Field:
			if sel.Name == property {
				return true
			}
		case *ast.InlineFragment:
			return hasProperty(ctx, sel.SelectionSet, property)
		case *ast.FragmentSpread:
			fragment := reqCtx.Doc.Fragments.ForName(sel.Name)
			return hasProperty(ctx, fragment.SelectionSet, property)
		}
	}
	return false
}
