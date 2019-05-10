package graphql

// Request describes the GraphQL input.
type Request struct {
	Query     string `json:"query"`
	Variables Map    `json:"variables"`
}
