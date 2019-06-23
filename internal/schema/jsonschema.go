package schema

//go:generate go-bindata -o ./assets.go -ignore ".*\\.go" -pkg schema ./...

import (
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/qri-io/jsonschema"
)

type Registry struct {
	V1 struct {
		Accounts struct {
			PATCH jsonschema.RootSchema
			POST  jsonschema.RootSchema
			PUT   jsonschema.RootSchema
		}
		Transactions struct {
			PATCH jsonschema.RootSchema
			POST  jsonschema.RootSchema
			PUT   jsonschema.RootSchema
		}
	}
}

func NewSchemaRegistry() (*Registry, error) {
	registry := Registry{}
	mustLoadSchema(&registry.V1.Accounts.POST, "v1_accounts_POST.json5")
	mustLoadSchema(&registry.V1.Accounts.PATCH, "v1_accounts_id_PATCH.json5")
	mustLoadSchema(&registry.V1.Accounts.PUT, "v1_accounts_id_PUT.json5")
	mustLoadSchema(&registry.V1.Transactions.PATCH, "v1_transactions_id_PATCH.json5")
	mustLoadSchema(&registry.V1.Transactions.POST, "v1_transactions_id_POST.json5")
	mustLoadSchema(&registry.V1.Transactions.PUT, "v1_transactions_id_PUT.json5")
	return &registry, nil
}

func mustLoadSchema(target *jsonschema.RootSchema, assetPath string) {
	if err := json.Unmarshal(MustAsset(assetPath), target); err != nil {
		panic(errors.Wrapf(err, "failed parsing JSON schema '%s'", assetPath))
	}
}
