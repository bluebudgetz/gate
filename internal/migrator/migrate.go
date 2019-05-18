package migrator

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/bluebudgetz/common/pkg/logging"
	"github.com/bluebudgetz/gate/internal/assets"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type account struct {
	ID     *int
	Name   string
	Parent *account
}

type Migrator struct {
	db       *sql.DB
	http     *http.Client
	port     int
	accounts struct {
		myEmployer      account
		acmeBank        account
		aig             account
		myBankAccount   account
		loans           account
		houseMortgage   account
		insurances      account
		lifeInsurance   account
		healthInsurance account
	}
}

func New(db *sql.DB, port int) (*Migrator, error) {
	return &Migrator{db: db, http: &http.Client{Timeout: 30 * time.Second}, port: port}, nil
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

	logging.Log.Info("Migration successful!")
	return nil
}

func (m *Migrator) populateAccounts() error {
	var err error

	acc := &m.accounts
	acc.myEmployer = account{Name: "Big Company"}
	acc.acmeBank = account{Name: "A.C.M.E Bank"}
	acc.aig = account{Name: "A.I.G"}
	acc.myBankAccount = account{Name: "My Account"}
	acc.loans = account{Name: "Loans", Parent: &acc.myBankAccount}
	acc.insurances = account{Name: "Insurances", Parent: &acc.myBankAccount}
	acc.houseMortgage = account{Name: "Mortgage", Parent: &acc.loans}
	acc.lifeInsurance = account{Name: "Life Insurance", Parent: &acc.insurances}
	acc.healthInsurance = account{Name: "Health Insurance", Parent: &acc.insurances}

	accounts := [...]*account{
		&acc.myEmployer,
		&acc.acmeBank,
		&acc.aig,
		&acc.myBankAccount,
		&acc.loans,
		&acc.insurances,
		&acc.houseMortgage,
		&acc.lifeInsurance,
		&acc.healthInsurance,
	}
	for _, account := range accounts {
		err = m.populateAccount(account)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *Migrator) populateAccount(account *account) error {
	var parentID *int = nil
	if account.Parent != nil {
		parentID = account.Parent.ID
	}

	requestBody := map[string]interface{}{"name": account.Name, "parentId": parentID}
	responseBody := struct{ ID int `json:"id"` }{}
	err := m.sendRequest(http.MethodPost, fmt.Sprintf("http://localhost:%d%s", m.port, "/v1/accounts"), requestBody, &responseBody)
	if err != nil {
		return errors.Wrapf(err, "failed creating account %+v", account)
	}

	account.ID = &responseBody.ID
	return nil
}

func (m *Migrator) populateTransactions() error {
	var err error
	var acc = &m.accounts

	err = m.populateTransaction("Initialization", *acc.myEmployer.ID, *acc.myBankAccount.ID, 5000, "June salary")
	if err != nil {
		return err
	}

	err = m.populateTransaction("Initialization", *acc.loans.ID, *acc.acmeBank.ID, 590, "Loan payment")
	if err != nil {
		return err
	}

	err = m.populateTransaction("Initialization", *acc.lifeInsurance.ID, *acc.aig.ID, 199, "Life insurance")
	if err != nil {
		return err
	}

	err = m.populateTransaction("Initialization", *acc.healthInsurance.ID, *acc.aig.ID, 98, "Health insurance")
	if err != nil {
		return err
	}

	return nil
}

func (m *Migrator) populateTransaction(origin string, sourceAccountId int, targetAccountId int, amount float64, comments string) error {
	requestBody := map[string]interface{}{
		"origin":          origin,
		"sourceAccountId": sourceAccountId,
		"targetAccountId": targetAccountId,
		"amount":          amount,
		"comments":        comments,
	}
	responseBody := struct{ ID int `json:"id"` }{}
	err := m.sendRequest(http.MethodPost, fmt.Sprintf("http://localhost:%d%s", m.port, "/v1/transactions"), requestBody, &responseBody)
	if err != nil {
		return errors.Wrapf(err, "failed creating transaction %+v", requestBody)
	}
	return nil
}

func (m *Migrator) sendRequest(method, url string, requestBody interface{}, responseBody interface{}) error {
	var body io.Reader = nil

	if requestBody != nil {
		jsonBody, err := json.Marshal(requestBody)
		if err != nil {
			return errors.Wrapf(err, "failed marshalling JSON: %+v", requestBody)
		}
		body = bytes.NewBuffer(jsonBody)
	}

	request, err := http.NewRequest(method, url, body)
	response, err := m.http.Do(request)
	if err != nil {
		return errors.Wrapf(err, "failed invoking HTTP request for: %+v", requestBody)
	}
	defer response.Body.Close()

	if response.StatusCode >= http.StatusBadRequest {
		return errors.Wrapf(err, "received HTTP status code %d for: ", response.StatusCode, requestBody)
	}

	location, err := response.Location()
	if err != nil {
		if err == http.ErrNoLocation {
			location = nil
		} else {
			return errors.Wrapf(err, "failed extracting 'Location' header from response")
		}
	}

	if location != nil {
		return m.sendRequest(http.MethodGet, location.String(), nil, responseBody)
	}

	responseBodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return errors.Wrapf(err, "failed reading response body: %+v", requestBody)
	}

	if len(responseBodyBytes) > 0 {
		err = json.Unmarshal(responseBodyBytes, responseBody)
		if err != nil {
			return errors.Wrapf(err, "failed unmarshalling response body: %+v", string(responseBodyBytes))
		}
	}
	return nil
}
