package graphql

// Query represents a query.
type Query struct {
	fields []*Field
}

func (query *Query) AddField(field *Field) {
	query.fields = append(query.fields, field)
}

func (query *Query) Fields() []*Field {
	return query.fields
}

func (query *Query) Parent() FieldContainer {
	return nil
}
