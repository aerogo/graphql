package graphql

// Operation represents an operation insides a GraphQL request.
type Operation struct {
	typ    string
	fields []*Field
}

// AddField adds a field to the operation.
func (operation *Operation) AddField(field *Field) {
	operation.fields = append(operation.fields, field)
}

// Fields returns the list of fields inside the operation.
func (operation *Operation) Fields() []*Field {
	return operation.fields
}

// Parent is always nil for queries.
func (operation *Operation) Parent() FieldContainer {
	return nil
}
