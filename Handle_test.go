package graphql_test

import (
	"fmt"
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

	app.OnStart(func() {
		response, err := client.Post(fmt.Sprintf("http://localhost:%d/", app.Config.Ports.HTTP)).End()
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.True(t, response.Ok(), "Status %d", response.StatusCode())

		err = syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
		assert.NoError(t, err)
	})

	app.Run()
}
