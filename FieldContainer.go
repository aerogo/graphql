package graphql

type FieldContainer interface {
	AddField(*Field)
	Parent() FieldContainer
}
