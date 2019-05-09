package graphql

import "github.com/aerogo/aero"

// Handle deals with a GraphQL request and responds to it.
func Handle(ctx *aero.Context) string {
	return ctx.JSON("abc")
}
