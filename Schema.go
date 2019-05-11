package graphql

type __Schema struct {
	Types            []__SchemaType     `json:"types"`
	QueryType        __QueryType        `json:"queryType"`
	MutationType     __MutationType     `json:"mutationType"`
	SubscriptionType __SubscriptionType `json:"subscriptionType"`
}

type __SchemaType struct {
	Name string `json:"name"`
}

type __QueryType struct {
	Name string `json:"name"`
}

type __MutationType struct {
	Name string `json:"name"`
}

type __SubscriptionType struct {
	Name string `json:"name"`
}
