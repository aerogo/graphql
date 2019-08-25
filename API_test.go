package graphql_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/aerogo/aero"
	"github.com/aerogo/graphql"
	"github.com/aerogo/nano"
	"github.com/akyoto/assert"
	jsoniter "github.com/json-iterator/go"
)

// User as a type of sample data
type User struct {
	ID      string `json:"id"`
	Nick    string `json:"nick"`
	Website string `json:"website"`
}

func Test(t *testing.T) {
	// Fill database with sample data
	db := nano.New(nano.Configuration{Port: 5000}).Namespace("test").RegisterTypes((*User)(nil))
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
	assert.Nil(t, err)

	gqlRequest := &graphql.Request{
		Query: string(query),
		Variables: graphql.Map{
			"id": testUser.ID,
		},
	}

	gqlRequestBody, err := jsoniter.Marshal(gqlRequest)

	if err != nil {
		t.Fatal(err)
	}

	request := httptest.NewRequest("POST", "/", bytes.NewReader(gqlRequestBody))
	response := httptest.NewRecorder()
	app.ServeHTTP(response, request)

	// Error checks
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, http.StatusOK, response.Code)
	assert.True(t, strings.Contains(response.Body.String(), testUser.Nick))
}
