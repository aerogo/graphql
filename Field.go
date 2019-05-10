package graphql

import (
	"errors"

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

func (field *Field) Resolve(obj interface{}, db Database) (interface{}, error) {
	if obj == nil {
		if len(field.arguments) == 1 && field.arguments["id"] != nil {
			return db.Get(field.name, field.arguments["id"].(string))
		}

		return nil, errors.New("This parameter is not available in the resolver")
	}

	_, _, v, err := mirror.GetChildField(obj, field.name)

	if err != nil {
		return nil, err
	}

	return v.Interface(), nil
}
