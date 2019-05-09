package graphql

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
