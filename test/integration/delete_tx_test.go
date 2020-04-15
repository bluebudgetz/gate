package integration

import (
	"testing"

	"github.com/bluebudgetz/gate/test"
)

func TestDeleteTx(t *testing.T) {
	app, cleanup := test.Run(t)
	defer cleanup()

	test.MustTx(t, app, "salary_20190109")
	test.DeleteTx(t, app, "salary_20190109")
	if account := test.GetTx(t, app, "salary_20190109"); account != nil {
		t.Fatalf("transaction was not deleted properly: %+v", account)
	}
}
