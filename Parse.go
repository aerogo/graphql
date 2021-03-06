package graphql

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/aerogo/aero"
	jsoniter "github.com/json-iterator/go"
)

// Parse parses the request from the body reader and returns a GraphQL document.
func Parse(ctx aero.Context) (*Document, error) {
	httpRequest := ctx.Request()
	reader := httpRequest.Body().Reader()
	request := Request{}
	var err error

	if httpRequest.Header("Content-Type") == "application/graphql" {
		// Body contains only the query
		body, err := ioutil.ReadAll(reader)

		if err != nil {
			return nil, err
		}

		request.Query = string(body)
	} else {
		// Body contains full GraphQL request
		decoder := jsoniter.NewDecoder(reader)
		err := decoder.Decode(&request)

		if err != nil {
			return nil, err
		}
	}

	document := &Document{}
	var currentContainer FieldContainer

	// State
	processedUntil := 0

	// Loop over the characters
	for i := 0; i < len(request.Query); i++ {
		switch request.Query[i] {
		case '{':
			blockPrefix := request.Query[processedUntil:i]
			blockPrefix = strings.TrimSpace(blockPrefix)

			if currentContainer != nil {
				field := &Field{
					name:   blockPrefix,
					parent: currentContainer,
				}

				argumentsPos := strings.Index(blockPrefix, "(")

				if argumentsPos != -1 {
					field.name = strings.TrimSpace(blockPrefix[:argumentsPos])
					field.arguments, err = parseArguments(blockPrefix[argumentsPos+1:len(blockPrefix)-1], request.Variables)

					if err != nil {
						return nil, err
					}
				}

				currentContainer.AddField(field)
				currentContainer = field
			}

			switch blockPrefix {
			// Query
			case "", "query":
				document.Operation = &Operation{
					typ: "query",
				}

				currentContainer = document.Operation

			// Mutation
			case "mutation":
				document.Operation = &Operation{
					typ: "mutation",
				}

				currentContainer = document.Operation
			}

			processedUntil = i + 1

		case '}':
			processedUntil = i + 1

			if currentContainer == nil {
				return nil, fmt.Errorf("Encountered '}' without an opening '{'")
			}

			currentContainer = currentContainer.Parent()

		case '\n':
			if currentContainer != nil {
				blockPrefix := request.Query[processedUntil:i]
				blockPrefix = strings.TrimSpace(blockPrefix)

				if len(blockPrefix) > 0 {
					field := &Field{
						name:   blockPrefix,
						parent: currentContainer,
					}

					argumentsPos := strings.Index(blockPrefix, "(")

					if argumentsPos != -1 {
						field.name = strings.TrimSpace(blockPrefix[:argumentsPos])
						field.arguments, err = parseArguments(blockPrefix[argumentsPos+1:len(blockPrefix)-1], request.Variables)

						if err != nil {
							return nil, err
						}
					}

					currentContainer.AddField(field)
				}
			}

			processedUntil = i + 1
		}
	}

	return document, nil
}

func parseArguments(raw string, variables Map) (Map, error) {
	arguments := Map{}

	// TODO: Use ignore.Reader
	lines := strings.Split(raw, ",")

	for _, line := range lines {
		parts := strings.Split(line, ":")
		name := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch {
		// Variable
		case strings.HasPrefix(value, "$"):
			varName := strings.TrimPrefix(value, "$")
			value, found := variables[varName]

			if !found {
				return nil, fmt.Errorf("Variable %s doesn't exist", varName)
			}

			arguments[name] = value

		// String
		case strings.HasPrefix(value, `"`) && strings.HasSuffix(value, `"`):
			value = value[1 : len(value)-1]
			arguments[name] = value

		// Float
		case strings.Contains(value, "."):
			floatValue, err := strconv.ParseFloat(value, 64)

			if err != nil {
				return nil, err
			}

			arguments[name] = floatValue

		// Int
		default:
			intValue, err := strconv.Atoi(value)

			if err != nil {
				return nil, err
			}

			arguments[name] = intValue
		}
	}

	return arguments, nil
}
