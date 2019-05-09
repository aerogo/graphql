package graphql

type Field struct {
	name      string
	arguments []string
	fields    []*Field
	parent    FieldContainer
}

func (field *Field) AddField(newField *Field) {
	field.fields = append(field.fields, newField)
}

func (field *Field) Parent() FieldContainer {
	return field.parent
}
