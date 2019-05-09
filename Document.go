package graphql

import (
	"github.com/akyoto/color"
)

// Document represents a GraphQL request.
type Document struct {
	Query *Query
}

// Execute executes the operations defined in the GraphQL document.
func (document *Document) Execute(db Database) *Response {
	color.Yellow("EXECUTE")
	data := map[string]interface{}{}
	var err error
	var allErrors []string

	if document.Query != nil {
		for _, field := range document.Query.Fields {
			data[field.name], err = field.Resolve(db)

			if err != nil {
				allErrors = append(allErrors, err.Error())
			}
		}
	}

	return &Response{
		Data:   data,
		Errors: allErrors,
	}
}
