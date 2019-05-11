package graphql

import (
	"net/http"

	"github.com/aerogo/aero"
)

// API represents the API configuration for GraphQL.
type API struct {
	db            Database
	rootResolvers []Resolver
}

// New creates a new GraphQL API.
func New(db Database) *API {
	return &API{
		db: db,
	}
}

// AddRootResolver adds a new resolver for root queries.
// The resolver can return the resolved object and a bool flag
// that determines whether the request has been dealt with or not.
func (api *API) AddRootResolver(resolver Resolver) {
	api.rootResolvers = append(api.rootResolvers, resolver)
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
