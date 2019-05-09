package graphql

type Collection struct {
	Name      string
	Arguments []string
	Fields    []*Field
}

type Field struct {
	Name   string
	Fields []*Field
}
