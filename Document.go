package graphql

import (
	"reflect"
)

type Variables map[string]interface{}

// Document represents a GraphQL request.
type Document struct {
	Query *Query
}

// Execute executes the operations defined in the GraphQL document.
func (document *Document) Execute(db Database) *Response {
	var data Variables
	var allErrors []string

	if document.Query != nil {
		data, allErrors = resolve(document.Query, nil, db)
	}

	return &Response{
		Data:   data,
		Errors: allErrors,
	}
}

func resolve(container FieldContainer, parent interface{}, db Database) (Variables, []string) {
	var allErrors []string
	var errors []string
	data := Variables{}

	for _, field := range container.Fields() {
		obj, err := field.Resolve(parent, db)

		if err != nil {
			allErrors = append(allErrors, err.Error())
		}

		kind := reflect.Indirect(reflect.ValueOf(obj)).Kind()

		if kind != reflect.Struct {
			data[field.name] = obj
			continue
		}

		if len(field.fields) == 0 {
			continue
		}

		data[field.name], errors = resolve(field, obj, db)

		if errors != nil {
			allErrors = append(allErrors, errors...)
		}
	}

	return data, allErrors
}
