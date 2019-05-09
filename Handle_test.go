package graphql_test

import (
	"fmt"
	"io/ioutil"
	"syscall"
	"testing"

	"github.com/aerogo/aero"
	"github.com/aerogo/graphql"
	"github.com/aerogo/http/client"
	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	app := aero.New()
	app.Post("/", graphql.Handle)
	query, err := ioutil.ReadFile("testdata/simple.gql")
	assert.NoError(t, err)

	app.OnStart(func() {
		// Request
		request := client.Post(fmt.Sprintf("http://localhost:%d/", app.Config.Ports.HTTP))
		request = request.Body(query)
		response, err := request.End()

		// Error checks
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.True(t, response.Ok(), "Status %d", response.StatusCode())

		// Kill server
		err = syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
		assert.NoError(t, err)
	})

	app.Run()
}
