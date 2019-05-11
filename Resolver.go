package graphql

// Resolver represents the function signature of resolver functions.
type Resolver = func(name string, arguments Map) (interface{}, error, bool)
