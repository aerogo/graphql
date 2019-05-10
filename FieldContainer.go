package graphql

type FieldContainer interface {
	AddField(*Field)
	Fields() []*Field
	Parent() FieldContainer
}
