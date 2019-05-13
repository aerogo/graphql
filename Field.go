package graphql

import (
	"errors"
	"fmt"
	"reflect"
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

// Returns the field's parent.
func (field *Field) Parent() FieldContainer {
	return field.parent
}

// Resolve resolves the field value for the given parent in the given database.
func (field *Field) Resolve(parent interface{}, api *API) (interface{}, error) {
	// If we have no parent object, treat it as a root query
	if parent == nil {
		return field.ResolveRootQuery(api)
	}

	// Allow querying the current type name
	if field.name == "__typename" {
		t := reflect.TypeOf(parent)

		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}

		return t.Name(), nil
	}

	// Aliases
	name := field.name
	typeName := reflect.TypeOf(parent).Name()
	aliasMap := api.aliases[typeName]

	if aliasMap != nil {
		aliasedName, hasAlias := aliasMap[field.name]

		if hasAlias {
			name = aliasedName
		}
	}

	// Fields that are direct descendants
	structField, _, value, err := mirror.GetChildField(parent, name)

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
	// Custom resolvers
	for _, resolve := range api.rootResolvers {
		obj, err, ok := resolve(field.name, field.arguments)

		if ok {
			return obj, err
		}
	}

	// Schema query
	if field.name == "__schema" {
		return api.schema, nil
	}

	// "All" queries
	if strings.HasPrefix(field.name, "all") {
		return field.ResolveAll(api)
	}

	// Return an error if the type doesn't exist
	if !api.db.HasType(field.name) {
		return nil, fmt.Errorf("Type '%s' does not exist", field.name)
	}

	// Single object queries
	if len(field.arguments) != 1 {
		return nil, errors.New("Single object queries require must specify an ID and nothing else")
	}

	for _, id := range field.arguments {
		return api.db.Get(field.name, id.(string))
	}

	// This code is actually unreachable,
	// but the linter is too dumb to realize that
	// so we're going to end it with a return statement.
	return nil, nil
}

// ResolveAll returns a list of objects that matches the filter arguments.
func (field *Field) ResolveAll(api *API) (interface{}, error) {
	records := []interface{}{}
	typeName := strings.TrimPrefix(field.name, "all")

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
