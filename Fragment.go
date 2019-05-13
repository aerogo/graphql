package graphql

// Fragment is a reusable specification of fields.
type Fragment struct {
	fields []*Field
}

// AddField adds a field to the fragment.
func (fragment *Fragment) AddField(newField *Field) {
	fragment.fields = append(fragment.fields, newField)
}

// Fields returns the list of fields inside the fragment.
func (fragment *Fragment) Fields() []*Field {
	return fragment.fields
}

// Parent is always nil for fragments.
func (fragment *Fragment) Parent() FieldContainer {
	return nil
}
