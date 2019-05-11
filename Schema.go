package graphql

type Schema struct {
	Types            []SchemaType     `json:"types"`
	QueryType        QueryType        `json:"queryType"`
	MutationType     MutationType     `json:"mutationType"`
	SubscriptionType SubscriptionType `json:"subscriptionType"`
}

type SchemaType struct {
	Name string `json:"name"`
}

type QueryType struct {
	Name string `json:"name"`
}

type MutationType struct {
	Name string `json:"name"`
}

type SubscriptionType struct {
	Name string `json:"name"`
}
