package graphql

import "reflect"

type AliasMap = map[string]string

// registerJSONAliases checks all fields for json tags and adds the tags as field aliases.
func registerJSONAliases(typ reflect.Type, aliases map[string]AliasMap) {
	typeAliases := AliasMap{}
	aliases[typ.Name()] = typeAliases

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		jsonTag := field.Tag.Get("json")

		if jsonTag == "" || jsonTag == "-" {
			continue
		}

		typeAliases[jsonTag] = field.Name
	}
}
