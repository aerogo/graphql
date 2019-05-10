package graphql

// Response describes the GraphQL output.
type Response struct {
	Data   interface{} `json:"data"`
	Errors []string    `json:"errors,omitempty"`
}
