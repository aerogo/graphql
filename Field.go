package graphql

import "errors"

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

func (field *Field) Parent() FieldContainer {
	return field.parent
}

func (field *Field) Resolve(db Database) (interface{}, error) {
	if len(field.arguments) == 1 && field.arguments["id"] != nil {
		return db.Get(field.name, field.arguments["id"].(string))
	}

	return nil, errors.New("This parameter is not available in the resolver")
}
