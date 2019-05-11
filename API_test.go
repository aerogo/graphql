package graphql_test

import (
	"fmt"
	"io/ioutil"
	"strings"
	"syscall"
	"testing"

	"github.com/aerogo/aero"
	"github.com/aerogo/graphql"
	"github.com/aerogo/http/client"
	"github.com/aerogo/nano"
	"github.com/stretchr/testify/assert"
)

// User as a type of sample data
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

	testUser := &User{
		ID:      "4J6qpK1ve",
		Nick:    "Akyoto",
		Website: "eduardurbach.com",
	}

	db.Set("User", testUser.ID, testUser)

	otherUser := &User{
		ID:      "VJOK1ckvx",
		Nick:    "Scott",
		Website: "github.com/soulcramer",
	}

	db.Set("User", otherUser.ID, otherUser)

	// Create new API
	api := graphql.New(db)

	// Create web app
	app := aero.New()
	app.Post("/", api.Handler())
	query, err := ioutil.ReadFile("testdata/simple.gql")
	assert.NoError(t, err)

	app.OnStart(func() {
		// Request
		request := client.Post(fmt.Sprintf("http://localhost:%d/", app.Config.Ports.HTTP))

		request = request.BodyJSON(&graphql.Request{
			Query: string(query),
			Variables: graphql.Map{
				"id": testUser.ID,
			},
		})

		response, err := request.End()

		// Error checks
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.True(t, response.Ok(), "Status %d", response.StatusCode())
		fmt.Println(response.String())
		assert.True(t, strings.Contains(response.String(), testUser.Nick))

		// Kill server
		err = syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
		assert.NoError(t, err)
	})

	app.Run()
}
