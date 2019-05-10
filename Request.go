package graphql

type Request struct {
	Query     string        `json:"query"`
	Variables KeyValueStore `json:"variables"`
}
