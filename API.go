package graphql

import (
	"net/http"
	"reflect"

	"github.com/aerogo/aero"
)

// API represents the API configuration for GraphQL.
type API struct {
	// The database interface
	db Database

	// Custom root query resolvers
	rootResolvers []Resolver

	// The schema we compile on creation
	schema Schema

	// A map of type names to field aliases
	aliases map[string]AliasMap
}

// New creates a new GraphQL API.
func New(db Database) *API {
	aliases := map[string]AliasMap{}

	// Introspection
	schemaTypes := []SchemaType{}

	for _, typ := range db.Types() {
		schemaTypes = append(schemaTypes, SchemaType{
			Name: typ.Name(),
		})

		if typ.Kind() != reflect.Struct {
			continue
		}

		registerJSONAliases(typ, aliases)
	}

	schemaType := SchemaType{
		Name: "__Schema",
	}

	schemaTypes = append(schemaTypes)

	queryType := QueryType{
		Name: "Query",
	}

	mutationType := MutationType{
		Name: "Mutation",
	}

	subscriptionType := SubscriptionType{
		Name: "Subscription",
	}

	schema := Schema{
		Types:            schemaTypes,
		QueryType:        queryType,
		MutationType:     mutationType,
		SubscriptionType: subscriptionType,
	}

	registerJSONAliases(reflect.TypeOf(schema), aliases)
	registerJSONAliases(reflect.TypeOf(schemaType), aliases)
	registerJSONAliases(reflect.TypeOf(queryType), aliases)
	registerJSONAliases(reflect.TypeOf(mutationType), aliases)
	registerJSONAliases(reflect.TypeOf(subscriptionType), aliases)

	return &API{
		db:      db,
		schema:  schema,
		aliases: aliases,
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
			return ctx.Error(
				http.StatusBadRequest,
				ctx.JSON(&Response{
					Errors: []string{
						err.Error(),
					},
				}),
			)
		}

		return ctx.JSON(document.Execute(api))
	}
}
