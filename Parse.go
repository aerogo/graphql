package graphql

import (
	"fmt"
	"io"
	"io/ioutil"
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
					field.arguments = parseArguments(blockPrefix[argumentsPos+1 : len(blockPrefix)-1])
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

func parseArguments(raw string) map[string]string {
	arguments := map[string]string{}

	// TODO: Use ignore.Reader
	lines := strings.Split(raw, ",")

	for _, line := range lines {
		parts := strings.Split(line, ":")
		name := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		arguments[name] = value
	}

	return arguments
}
