package graphql

import (
	"net/http"

	"github.com/aerogo/aero"
)

// Handler returns a function that deals with a GraphQL request and responds to it.
func Handler(db Database) aero.Handle {
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

		response := ctx.JSON(document.Execute(db))
		return response
	}
}
