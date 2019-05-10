package graphql

type Request struct {
	Query     string `json:"query"`
	Variables string `json:"variables"`
}
