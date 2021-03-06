package graphql

import "reflect"

// Database is an interface for any kind of database.
type Database interface {
	Get(collection string, id string) (interface{}, error)
	Set(collection string, id string, obj interface{})
	Delete(collection string, id string) bool
	All(collection string) chan interface{}
	Types() map[string]reflect.Type
	HasType(typeName string) bool
}
