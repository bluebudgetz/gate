package integration

import (
	"testing"

	"github.com/bluebudgetz/gate/test"
)

func TestDeleteTx(t *testing.T) {
	app, cleanup := test.Run(t)
	defer cleanup()

	test.MustTx(t, app, "salary-2020-01-01")
	test.DeleteTx(t, app, "salary-2020-01-01")
	if tx := test.GetTx(t, app, "salary-2020-01-01"); tx != nil {
		t.Fatalf("transaction was not deleted properly: %+v", tx)
	}
}
