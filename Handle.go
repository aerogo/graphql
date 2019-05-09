package graphql

import (
	"fmt"
	"net/http"

	"github.com/aerogo/aero"
)

// Handler returns a function that deals with a GraphQL request and responds to it.
func Handler(db Database) aero.Handle {
	return func(ctx *aero.Context) string {
		document, err := Parse(ctx.Request().Body().Reader())

		if err != nil {
			return ctx.Error(http.StatusBadRequest, err)
		}

		response := ctx.JSON(document.Execute(db))
		fmt.Println(response)
		return response
	}
}
