package integration

import (
	"testing"

	"github.com/bluebudgetz/gate/test"
)

func TestDeleteAccount(t *testing.T) {
	app, cleanup := test.Run(t)
	defer cleanup()

	test.MustAccount(t, app, "company")
	test.DeleteAccount(t, app, "company")
	if account := test.GetAccount(t, app, "company"); account != nil {
		t.Fatalf("account was not deleted properly: %+v", account)
	}
}
