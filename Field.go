package graphql

import (
	"errors"
	"fmt"
	"strings"

	"github.com/aerogo/mirror"
)

// Field represents a queryable field.
type Field struct {
	name      string
	arguments Map
	fields    []*Field
	parent    FieldContainer
}

// AddField adds a field to the query.
func (field *Field) AddField(newField *Field) {
	field.fields = append(field.fields, newField)
}

// Fields returns the list of fields inside the query.
func (field *Field) Fields() []*Field {
	return field.fields
}

// Parent is always nil for queries.
func (field *Field) Parent() FieldContainer {
	return field.parent
}

// Resolve resolves the field value for the given parent in the given database.
func (field *Field) Resolve(parent interface{}, api *API) (interface{}, error) {
	if parent == nil {
		return field.ResolveRootQuery(api)
	}

	structField, _, value, err := mirror.GetChildField(parent, field.name)

	if err != nil {
		return nil, err
	}

	if structField.Tag.Get("private") == "true" {
		return nil, fmt.Errorf("'%s' is a private field", field.name)
	}

	return value.Interface(), nil
}

// ResolveRootQuery resolves a root query.
func (field *Field) ResolveRootQuery(api *API) (interface{}, error) {
	// "All" queries
	if strings.HasPrefix(field.name, "All") {
		return field.ResolveAll(api)
	}

	// Single object queries
	if len(field.arguments) != 1 || field.arguments["ID"] == nil {
		return nil, errors.New("Can only query objects by 'ID'")
	}

	if !api.db.HasType(field.name) {
		return nil, fmt.Errorf("Type '%s' does not exist", field.name)
	}

	return api.db.Get(field.name, field.arguments["ID"].(string))
}

// ResolveAll returns a list of objects that matches the filter arguments.
func (field *Field) ResolveAll(api *API) (interface{}, error) {
	records := []interface{}{}
	typeName := strings.TrimPrefix(field.name, "All")

	for argName, argValue := range field.arguments {
		if !strings.Contains(argName, "_") {
			continue
		}

		delete(field.arguments, argName)
		argName = strings.ReplaceAll(argName, "_", ".")
		field.arguments[argName] = argValue
	}

	for record := range api.db.All(typeName) {
		matchingFields := 0

		for argName, argValue := range field.arguments {
			_, _, value, err := mirror.GetPublicField(record, argName)

			if err != nil {
				return nil, err
			}

			switch argValue.(type) {
			case string:
				if value.String() == argValue {
					matchingFields++
				}

			case int64:
				if value.Int() == argValue {
					matchingFields++
				}

			case float64:
				if value.Float() == argValue {
					matchingFields++
				}

			case bool:
				if value.Bool() == argValue {
					matchingFields++
					continue
				}
			}
		}

		if matchingFields == len(field.arguments) {
			records = append(records, record)
		}
	}

	return records, nil
}
