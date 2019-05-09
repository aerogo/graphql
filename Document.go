package graphql

// Document represents a GraphQL request.
type Document struct {
	Query *Query
}

// Execute executes the operations defined in the GraphQL document.
func (document *Document) Execute() *Response {
	return &Response{
		Data: "Hello",
	}
}
