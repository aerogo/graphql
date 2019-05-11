package graphql

import (
	"reflect"
)

// Document represents a GraphQL request.
type Document struct {
	Operation *Operation
}

// Execute executes the operations defined in the GraphQL document.
func (document *Document) Execute(api *API) *Response {
	var data interface{}
	var allErrors []string

	if document.Operation != nil {
		data, allErrors = resolve(document.Operation, nil, api)
	}

	return &Response{
		Data:   data,
		Errors: allErrors,
	}
}

func resolve(container FieldContainer, parent interface{}, api *API) (Map, []string) {
	var allErrors []string
	var errors []string
	data := Map{}

	for _, field := range container.Fields() {
		obj, err := field.Resolve(parent, api)

		if err != nil {
			allErrors = append(allErrors, err.Error())
		}

		value := reflect.ValueOf(obj)
		kind := reflect.Indirect(value).Kind()

		switch kind {
		case reflect.Slice:
			// Simple types can be inserted instantly
			sliceElementKind := reflect.TypeOf(obj).Elem().Kind()

			if sliceElementKind != reflect.Ptr && sliceElementKind != reflect.Struct {
				data[field.name] = obj
				continue
			}

			// If we have complex types as elements and we didn't specify a field, skip it
			if len(field.fields) == 0 {
				continue
			}

			// Create a slice with the data we requested
			slice := make([]Map, value.Len())

			for i := 0; i < value.Len(); i++ {
				element := value.Index(i).Interface()
				slice[i], errors = resolve(field, element, api)

				if errors != nil {
					allErrors = append(allErrors, errors...)
				}
			}

			data[field.name] = slice

		case reflect.Struct:
			if len(field.fields) == 0 {
				continue
			}

			data[field.name], errors = resolve(field, obj, api)

			if errors != nil {
				allErrors = append(allErrors, errors...)
			}

		default:
			data[field.name] = obj
		}
	}

	return data, allErrors
}
