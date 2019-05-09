package graphql

// Query represents a query.
type Query struct {
	Fields []*Field
}

func (query *Query) AddField(field *Field) {
	query.Fields = append(query.Fields, field)
}

func (query *Query) Parent() FieldContainer {
	return nil
}
