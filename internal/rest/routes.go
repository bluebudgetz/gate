package rest

import (
	"github.com/go-chi/chi"

	"github.com/bluebudgetz/gate/internal/rest/accounts"
	"github.com/bluebudgetz/gate/internal/rest/transactions"
)

func NewRoutes(accountsMgr accounts.Manager, txMgr transactions.Manager) func(chi.Router) {
	return func(r chi.Router) {
		r.Get("/routes", routesDoc(r))
		r.Route("/accounts", func(r chi.Router) {
			r.Get("/", accounts.List(accountsMgr))
			r.Post("/", accounts.Post(accountsMgr))
			r.Route("/{ID}", func(r chi.Router) {
				r.Get("/", accounts.Get(accountsMgr))
				r.Put("/", accounts.Put(accountsMgr))
				r.Patch("/", accounts.Patch(accountsMgr))
				r.Delete("/", accounts.Delete(accountsMgr))
			})
		})
		r.Route("/transactions", func(r chi.Router) {
			r.Get("/", transactions.List(txMgr))
			r.Post("/", transactions.Post(txMgr))
			r.Route("/{ID}", func(r chi.Router) {
				r.Get("/", transactions.Get(txMgr))
				r.Put("/", transactions.Put(txMgr))
				r.Patch("/", transactions.Patch(txMgr))
				r.Delete("/", transactions.Delete(txMgr))
			})
		})
	}
}
