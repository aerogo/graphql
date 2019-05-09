package graphql_test

import (
	"fmt"
	"io/ioutil"
	"syscall"
	"testing"

	"github.com/aerogo/aero"
	"github.com/aerogo/graphql"
	"github.com/aerogo/http/client"
	"github.com/aerogo/nano"
	"github.com/stretchr/testify/assert"
)

type User struct {
	ID      string `json:"id"`
	Nick    string `json:"nick"`
	Website string `json:"website"`
}

func Test(t *testing.T) {
	// Fill database with sample data
	db := nano.New(5000).Namespace("test").RegisterTypes((*User)(nil))
	defer db.Close()
	defer db.Clear("User")

	db.Set("User", "4J6qpK1ve", &User{
		ID:      "4J6qpK1ve",
		Nick:    "Akyoto",
		Website: "eduardurbach.com",
	})

	// Create web app
	app := aero.New()
	app.Post("/", graphql.Handler(db))
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
