package graphql

import (
	"reflect"
)

// Document represents a GraphQL request.
type Document struct {
	Query *Query
}

// Execute executes the operations defined in the GraphQL document.
func (document *Document) Execute(db Database) *Response {
	var data interface{}
	var allErrors []string

	if document.Query != nil {
		data, allErrors = resolve(document.Query, nil, db)
	}

	return &Response{
		Data:   data,
		Errors: allErrors,
	}
}

func resolve(container FieldContainer, parent interface{}, db Database) (Map, []string) {
	var allErrors []string
	var errors []string
	data := Map{}

	for _, field := range container.Fields() {
		obj, err := field.Resolve(parent, db)

		if err != nil {
			allErrors = append(allErrors, err.Error())
		}

		value := reflect.ValueOf(obj)
		kind := reflect.Indirect(value).Kind()

		switch kind {
		case reflect.Slice:
			if len(field.fields) == 0 {
				continue
			}

			slice := make([]Map, value.Len())

			for i := 0; i < value.Len(); i++ {
				element := value.Index(i).Interface()
				slice[i], errors = resolve(field, element, db)

				if errors != nil {
					allErrors = append(allErrors, errors...)
				}
			}

			data[field.name] = slice

		case reflect.Struct:
			if len(field.fields) == 0 {
				continue
			}

			data[field.name], errors = resolve(field, obj, db)

			if errors != nil {
				allErrors = append(allErrors, errors...)
			}

		default:
			data[field.name] = obj
		}
	}

	return data, allErrors
}
