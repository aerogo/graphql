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

	document := &Document{
		Definitions: []*Definition{},
	}

	// State
	lineStart := 0

	// Loop over the characters
	for i := 0; i < len(body); i++ {
		switch body[i] {
		case '{':
			blockPrefix := string(body[lineStart:i])
			blockPrefix = strings.TrimSpace(blockPrefix)
			fmt.Println("BLOCK", blockPrefix)
		case '}':

		case '\n':
			lineStart = i
		}
	}

	return document, nil
}
