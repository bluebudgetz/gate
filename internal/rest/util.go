package rest

import (
	"net/http"

	"github.com/golangly/errors"
	"github.com/golangly/webutil"
	"github.com/neo4j/neo4j-go-driver/neo4j"

	"github.com/bluebudgetz/gate/internal/util"
)

func GetNode(w http.ResponseWriter, r *http.Request, getQuery string, params map[string]interface{}, nodeToResponseMapper func(neo4j.Record) (interface{}, error)) {
	if result, err := util.GetNeo4jSession(r.Context()).Run(getQuery, params); err != nil {
		webutil.RenderWithStatusCode(w, r, http.StatusInternalServerError, err)
	} else if !result.Next() {
		if err := result.Err(); err != nil {
			webutil.RenderWithStatusCode(w, r, http.StatusInternalServerError, err)
		} else {
			webutil.RenderWithStatusCode(w, r, http.StatusNotFound, nil)
		}
	} else if response, err := nodeToResponseMapper(result.Record()); err != nil {
		webutil.RenderWithStatusCode(w, r, http.StatusInternalServerError, err)
	} else {
		webutil.RenderWithStatusCode(w, r, http.StatusOK, response)
	}
}

func PostNode(w http.ResponseWriter, r *http.Request, postQuery string, params map[string]interface{}, nodeToResponseMapper func(neo4j.Record) (interface{}, error)) {
	if result, err := util.GetNeo4jSession(r.Context()).Run(postQuery, params); err != nil {
		webutil.RenderWithStatusCode(w, r, http.StatusInternalServerError, err)
	} else if !result.Next() {
		webutil.RenderWithStatusCode(w, r, http.StatusInternalServerError, result.Err())
	} else if summary, err := result.Summary(); err != nil {
		webutil.RenderWithStatusCode(w, r, http.StatusInternalServerError, errors.Wrap(err, "summary failed"))
	} else if response, err := nodeToResponseMapper(result.Record()); err != nil {
		webutil.RenderWithStatusCode(w, r, http.StatusInternalServerError, err)
	} else if summary.Counters().NodesCreated() > 0 {
		webutil.RenderWithStatusCode(w, r, http.StatusCreated, response)
	} else {
		webutil.RenderWithStatusCode(w, r, http.StatusOK, response)
	}
}

func PatchNode(w http.ResponseWriter, r *http.Request, patchQuery string, params map[string]interface{}, nodeToResponseMapper func(neo4j.Record) (interface{}, error)) {
	if result, err := util.GetNeo4jSession(r.Context()).Run(patchQuery, params); err != nil {
		webutil.RenderWithStatusCode(w, r, http.StatusInternalServerError, err)
	} else if !result.Next() {
		if err := result.Err(); err != nil {
			webutil.RenderWithStatusCode(w, r, http.StatusInternalServerError, err)
		} else {
			webutil.RenderWithStatusCode(w, r, http.StatusNotFound, nil)
		}
	} else if summary, err := result.Summary(); err != nil {
		webutil.RenderWithStatusCode(w, r, http.StatusInternalServerError, errors.Wrap(err, "summary failed"))
	} else if response, err := nodeToResponseMapper(result.Record()); err != nil {
		webutil.RenderWithStatusCode(w, r, http.StatusInternalServerError, err)
	} else if summary.Counters().NodesCreated() > 0 {
		webutil.RenderWithStatusCode(w, r, http.StatusCreated, response)
	} else {
		webutil.RenderWithStatusCode(w, r, http.StatusOK, response)
	}
}

func PutNode(w http.ResponseWriter, r *http.Request, patchQuery string, params map[string]interface{}, nodeToResponseMapper func(neo4j.Record) (interface{}, error)) {
	if result, err := util.GetNeo4jSession(r.Context()).Run(patchQuery, params); err != nil {
		webutil.RenderWithStatusCode(w, r, http.StatusInternalServerError, err)
	} else if !result.Next() {
		if err := result.Err(); err != nil {
			webutil.RenderWithStatusCode(w, r, http.StatusInternalServerError, err)
		} else {
			webutil.RenderWithStatusCode(w, r, http.StatusNotFound, nil)
		}
	} else if response, err := nodeToResponseMapper(result.Record()); err != nil {
		webutil.RenderWithStatusCode(w, r, http.StatusInternalServerError, err)
	} else {
		webutil.RenderWithStatusCode(w, r, http.StatusOK, response)
	}
}

func DeleteNode(w http.ResponseWriter, r *http.Request, deleteQuery string, params map[string]interface{}) {
	if result, err := util.GetNeo4jSession(r.Context()).Run(deleteQuery, params); err != nil {
		webutil.RenderWithStatusCode(w, r, http.StatusInternalServerError, err)
	} else if summary, err := result.Summary(); err != nil {
		webutil.RenderWithStatusCode(w, r, http.StatusInternalServerError, err)
	} else if summary.Counters().NodesDeleted() == 0 {
		webutil.RenderWithStatusCode(w, r, http.StatusNotFound, nil)
	} else if summary.Counters().NodesDeleted() == 1 {
		webutil.RenderWithStatusCode(w, r, http.StatusNoContent, nil)
	} else {
		webutil.RenderWithStatusCode(w, r, http.StatusInternalServerError,
			errors.New("invalid amount of transactions deleted").
				AddTag("deleted", summary.Counters().NodesDeleted()),
		)
	}
}
