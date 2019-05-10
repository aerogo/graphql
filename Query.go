package graphql

// Query represents a query insides a GraphQL request.
type Query struct {
	fields []*Field
}

// AddField adds a field to the query.
func (query *Query) AddField(field *Field) {
	query.fields = append(query.fields, field)
}

// Fields returns the list of fields inside the query.
func (query *Query) Fields() []*Field {
	return query.fields
}

// Parent is always nil for queries.
func (query *Query) Parent() FieldContainer {
	return nil
}
