package graphql

// Database is an interface for any kind of database.
type Database interface {
	Get(collection string, id string) (interface{}, error)
	All(collection string) chan interface{}
	HasType(typeName string) bool
}

// These could be useful in the future:
// Set(collection string, id string, obj interface{})
// Delete(collection string, id string) bool
// Types() map[string]reflect.Type
