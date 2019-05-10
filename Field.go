package graphql

import (
	"errors"
	"fmt"
	"strings"

	"github.com/aerogo/mirror"
)

type ArgumentsList = map[string]interface{}

type Field struct {
	name      string
	arguments ArgumentsList
	fields    []*Field
	parent    FieldContainer
}

func (field *Field) AddField(newField *Field) {
	field.fields = append(field.fields, newField)
}

func (field *Field) Fields() []*Field {
	return field.fields
}

func (field *Field) Parent() FieldContainer {
	return field.parent
}

func (field *Field) Resolve(parent interface{}, db Database) (interface{}, error) {
	if parent == nil {
		return field.ResolveRootQuery(db)
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

func (field *Field) ResolveRootQuery(db Database) (interface{}, error) {
	if strings.HasPrefix(field.name, "All") {
		return field.ResolveAll(db)
	}

	if len(field.arguments) != 1 || field.arguments["ID"] == nil {
		return nil, errors.New("Can only query objects by 'ID'")
	}

	return db.Get(field.name, field.arguments["ID"].(string))
}

func (field *Field) ResolveAll(db Database) (interface{}, error) {
	records := []interface{}{}
	typeName := strings.TrimPrefix(field.name, "All")

	for record := range db.All(typeName) {
		matchingFields := 0

		for argName, argValue := range field.arguments {
			structField, _, value, err := mirror.GetChildField(record, argName)

			if err != nil {
				return nil, err
			}

			if structField.Tag.Get("private") == "true" {
				return nil, fmt.Errorf("'%s' is a private field", structField.Name)
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
