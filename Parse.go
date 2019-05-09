package graphql

import (
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
	"strings"
)

// Parse parses the request from the body reader and returns a GraphQL document.
func Parse(reader io.Reader) (*Document, error) {
	body, err := ioutil.ReadAll(reader)

	if err != nil {
		return nil, err
	}

	document := &Document{}
	var currentContainer FieldContainer

	// State
	processedUntil := 0

	// Loop over the characters
	for i := 0; i < len(body); i++ {
		switch body[i] {
		case '{':
			blockPrefix := string(body[processedUntil:i])
			blockPrefix = strings.TrimSpace(blockPrefix)

			if currentContainer != nil {
				field := &Field{
					name:   blockPrefix,
					parent: currentContainer,
				}

				argumentsPos := strings.Index(blockPrefix, "(")

				if argumentsPos != -1 {
					field.name = blockPrefix[:argumentsPos]
					field.arguments, err = parseArguments(blockPrefix[argumentsPos+1 : len(blockPrefix)-1])

					if err != nil {
						return nil, err
					}
				}

				fmt.Println(field.name)
				currentContainer.AddField(field)
				currentContainer = field
			}

			if blockPrefix == "query" {
				document.Query = &Query{}
				currentContainer = document.Query
			}

			processedUntil = i + 1
		case '}':
			processedUntil = i + 1
			currentContainer = currentContainer.Parent()

		case '\n':
			if currentContainer != nil {
				blockPrefix := string(body[processedUntil:i])
				blockPrefix = strings.TrimSpace(blockPrefix)

				if len(blockPrefix) > 0 {
					field := &Field{
						name:   blockPrefix,
						parent: currentContainer,
					}

					currentContainer.AddField(field)
					fmt.Println(field.name)
				}
			}

			processedUntil = i + 1
		}
	}

	return document, nil
}

func parseArguments(raw string) (ArgumentsList, error) {
	arguments := ArgumentsList{}

	// TODO: Use ignore.Reader
	lines := strings.Split(raw, ",")

	for _, line := range lines {
		parts := strings.Split(line, ":")
		name := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch {
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
