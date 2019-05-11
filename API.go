package graphql

import (
	"net/http"

	"github.com/aerogo/aero"
)

// API represents the API configuration for GraphQL.
type API struct {
	db Database
}

// New creates a new GraphQL API.
func New(db Database) *API {
	return &API{
		db: db,
	}
}

// Handler returns a function that deals with a GraphQL request and responds to it.
func (api *API) Handler() aero.Handle {
	return func(ctx *aero.Context) string {
		document, err := Parse(ctx)

		if err != nil {
			ctx.StatusCode = http.StatusBadRequest

			return ctx.JSON(&Response{
				Errors: []string{
					err.Error(),
				},
			})
		}

		return ctx.JSON(document.Execute(api))
	}
}
