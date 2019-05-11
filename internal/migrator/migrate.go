package migrator

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/99designs/gqlgen/graphql"
	"github.com/bluebudgetz/common/pkg/logging"
	"github.com/bluebudgetz/gate/internal/assets"
	"github.com/bluebudgetz/gate/internal/graphql/impl"
	"github.com/bluebudgetz/gate/internal/graphql/resolver"
	"github.com/bluebudgetz/gate/internal/middleware"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/pkg/errors"
	"github.com/vektah/gqlparser/ast"
	"github.com/vektah/gqlparser/parser"
	"github.com/vektah/gqlparser/validator"
	"io/ioutil"
	"os"
)

type createAccountResponse struct {
	Result struct {
		ID int `json:"id"`
	} `json:"createAccount"`
}

type Migrator struct {
	db     *sql.DB
	cfg    *impl.Config
	schema graphql.ExecutableSchema
	data   struct {
		accounts struct {
			myEmployer             createAccountResponse
			acmeBank               createAccountResponse
			aig                    createAccountResponse
			myBankAccount          createAccountResponse
			loansAccount           createAccountResponse
			insuranceAccount       createAccountResponse
			lifeInsuranceAccount   createAccountResponse
			healthInsuranceAccount createAccountResponse
		}
	}
}

func New(db *sql.DB) (*Migrator, error) {
	config := impl.Config{Resolvers: &resolver.Resolver{}}
	return &Migrator{
		db:     db,
		cfg:    &config,
		schema: impl.NewExecutableSchema(config),
	}, nil
}

func (m *Migrator) Migrate() error {

	// Extract all migrations to a temporary directory
	logging.Log.Info("Extracting migrations")
	tempMigrationsPath, err := ioutil.TempDir("", "bluebudgetzMigrations")
	if err != nil {
		return errors.Wrap(err, "failed extracting database migration files")
	}
	defer os.RemoveAll(tempMigrationsPath)
	if err = assets.RestoreAssets(tempMigrationsPath, "deployments/rdbms/migrations"); err != nil {
		return errors.Wrap(err, "failed extracting database migration files")
	}

	logging.Log.Info("Migrating schema")
	driver, err := postgres.WithInstance(m.db, &postgres.Config{})
	if err != nil {
		return errors.Wrap(err, "failed creating database migration driver")
	}

	// Migrate all the way down, and then all the way up
	rdbmsUrl := "file://" + tempMigrationsPath + "/deployments/rdbms/migrations"
	migrator, err := migrate.NewWithDatabaseInstance(rdbmsUrl, "postgres", driver)
	if err != nil {
		return errors.Wrap(err, "failed creating database migrator")
	}
	if err = migrator.Down(); err != nil && err != migrate.ErrNoChange {
		return errors.Wrap(err, "failed to drop database")
	}
	if err = migrator.Up(); err != nil && err != migrate.ErrNoChange {
		return errors.Wrap(err, "failed to migrate database")
	}

	return nil
}

func (m *Migrator) Populate() error {
	var err error

	err = m.populateAccounts()
	if err != nil {
		return err
	}

	err = m.populateTransactions()
	if err != nil {
		return err
	}

	return nil
}

func (m *Migrator) populateAccounts() error {
	var err error

	createQuery := `mutation($name: String!) { 
		createAccount(name: $name) { 
			id 
		} 
	}`
	createChildQuery := `mutation($name: String!, $parentId: Int!) { 
		createAccount(name: $name, parentId: $parentId) { 
			id 
		} 
	}`
	acc := &m.data.accounts

	err = m.runMutation(createQuery, map[string]interface{}{"name": "MyEmployer"}, &acc.myEmployer)
	if err != nil {
		return err
	}

	err = m.runMutation(createQuery, map[string]interface{}{"name": "A.C.M.E Bank Inc."}, &acc.acmeBank)
	if err != nil {
		return err
	}

	err = m.runMutation(createQuery, map[string]interface{}{"name": "A.I.G"}, &acc.aig)
	if err != nil {
		return err
	}

	err = m.runMutation(createQuery, map[string]interface{}{"name": "Account @ MyBank"}, &acc.myBankAccount)
	if err != nil {
		return err
	}

	err = m.runMutation(createChildQuery, map[string]interface{}{"name": "Loans", "parentId": acc.myBankAccount.Result.ID}, &acc.loansAccount)
	if err != nil {
		return err
	}

	err = m.runMutation(createChildQuery, map[string]interface{}{"name": "Insurance", "parentId": acc.myBankAccount.Result.ID}, &acc.insuranceAccount)
	if err != nil {
		return err
	}

	err = m.runMutation(createChildQuery, map[string]interface{}{"name": "Life", "parentId": acc.insuranceAccount.Result.ID}, &acc.lifeInsuranceAccount)
	if err != nil {
		return err
	}

	err = m.runMutation(createChildQuery, map[string]interface{}{"name": "Health", "parentId": acc.insuranceAccount.Result.ID}, &acc.healthInsuranceAccount)
	if err != nil {
		return err
	}

	return nil
}

func (m *Migrator) populateTransactions() error {
	var err error

	query := `mutation(	$origin: String!, 
						$sourceAccountId: Int!, 
						$targetAccountId: Int!, 
						$amount: Float!, 
						$comments: String) { 
		createTransaction(	origin: $origin, 
							sourceAccountId: $sourceAccountId, 
							targetAccountId: $targetAccountId, 
							amount: $amount, 
							comments: $comments) { 
			id 
		} 
	}`

	acc := &m.data.accounts

	err = m.runMutation(query, createTxVarsMap(acc.myEmployer, acc.myBankAccount, 5000, "June salary"), nil)
	if err != nil {
		return err
	}

	err = m.runMutation(query, createTxVarsMap(acc.loansAccount, acc.acmeBank, 590, "Loan payment"), nil)
	if err != nil {
		return err
	}

	err = m.runMutation(query, createTxVarsMap(acc.lifeInsuranceAccount, acc.aig, 199, "Life insurance"), nil)
	if err != nil {
		return err
	}

	err = m.runMutation(query, createTxVarsMap(acc.healthInsuranceAccount, acc.aig, 98, "Health insurance"), nil)
	if err != nil {
		return err
	}

	return nil
}

func (m *Migrator) runMutation(query string, variables map[string]interface{}, response interface{}) error {
	doc, gerr := parser.ParseQuery(&ast.Source{Input: query})
	if gerr != nil {
		return errors.Wrapf(gerr, "failed parsing query '%s'", query)
	}

	errs := validator.Validate(m.schema.Schema(), doc)
	if len(errs) != 0 {
		return errors.Wrapf(errs, "failed validating query '%s'", query)
	}

	op := doc.Operations.ForName("")
	if op == nil {
		return errors.Errorf("operation not found")
	}

	vars, gerr := validator.VariableValues(m.schema.Schema(), op, variables)
	if gerr != nil {
		return errors.Wrap(gerr, "failed parsing variables")
	}

	reqCtx := graphql.NewRequestContext(doc, query, vars)
	ctx := middleware.EnrichContext(graphql.WithRequestContext(context.Background(), reqCtx), m.db)

	mutation := m.schema.Mutation(ctx, op)
	if len(mutation.Errors) > 0 {
		return errors.Wrap(mutation.Errors, "mutation failed")
	}

	if response != nil {
		err := json.Unmarshal(mutation.Data, response)
		if err != nil {
			return errors.Wrap(err, "failed unmarshalling response")
		}
	}

	return nil
}

func createTxVarsMap(sourceAccountId createAccountResponse, targetAccountId createAccountResponse, amount float64, comments string) map[string]interface{} {
	return map[string]interface{}{
		"origin":          "Initialization",
		"sourceAccountId": sourceAccountId.Result.ID,
		"targetAccountId": targetAccountId.Result.ID,
		"amount":          amount,
		"comments":        comments,
	}
}
