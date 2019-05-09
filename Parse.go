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
					field.arguments = []string{blockPrefix[argumentsPos:]}
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
