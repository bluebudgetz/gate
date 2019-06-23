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
			POST jsonschema.RootSchema
			PUT  jsonschema.RootSchema
		}
	}
}

func NewSchemaRegistry() (*Registry, error) {
	registry := Registry{}

	if err := json.Unmarshal(MustAsset("v1_accounts_POST.json5"), &registry.V1.Accounts.POST); err != nil {
		return nil, errors.Wrapf(err, "failed parsing JSON schema")
	}
	if err := json.Unmarshal(MustAsset("v1_accounts_id_PATCH.json5"), &registry.V1.Accounts.PATCH); err != nil {
		return nil, errors.Wrapf(err, "failed parsing JSON schema")
	}
	if err := json.Unmarshal(MustAsset("v1_accounts_id_PUT.json5"), &registry.V1.Accounts.PUT); err != nil {
		return nil, errors.Wrapf(err, "failed parsing JSON schema")
	}
	return &registry, nil
}
