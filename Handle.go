package graphql

import (
	"net/http"

	"github.com/aerogo/aero"
)

// Handle deals with a GraphQL request and responds to it.
func Handle(ctx *aero.Context) string {
	document, err := Parse(ctx.Request().Body().Reader())

	if err != nil {
		return ctx.Error(http.StatusBadRequest, err)
	}

	return ctx.JSON(document)
}
