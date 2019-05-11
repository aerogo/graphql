package graphql

type Schema struct {
	Types []SchemaType `json:"types"`
}

type SchemaType struct {
	Name string `json:"name"`
}
