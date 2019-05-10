package graphql

type Request struct {
	Query     string    `json:"query"`
	Variables Variables `json:"variables"`
}
