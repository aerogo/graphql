package graphql

import (
	"net/http"

	"github.com/aerogo/aero"
)

// API represents the API configuration for GraphQL.
type API struct {
	db            Database
	rootResolvers []Resolver
	schema        *Schema
}

// New creates a new GraphQL API.
func New(db Database) *API {
	schemaTypes := []SchemaType{}

	for _, typ := range db.Types() {
		schemaTypes = append(schemaTypes, SchemaType{
			Name: typ.Name(),
		})
	}

	schema := &Schema{
		Types: schemaTypes,
	}

	return &API{
		db:     db,
		schema: schema,
	}
}

// AddRootResolver adds a new resolver for root queries.
// The resolver can return the resolved object, an error and a bool
// flag that determines whether the request has been dealt with or not.
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
