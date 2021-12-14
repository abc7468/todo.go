package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"testing"

	"github.com/abc7468/todo.go/model"
	"github.com/stretchr/testify/assert"
)

func TestTodo(t *testing.T) {
	getSessionID = func(r *http.Request) string {
		return "test"
	}
	os.Remove("./test.db")
	assert := assert.New(t)
	ah := MakeHandler("./test.db")
	defer ah.Close()
	ts := httptest.NewServer(ah)
	defer ts.Close()
	res, err := http.PostForm(ts.URL+"/todos", url.Values{"name": {"Test Todo"}})
	assert.NoError(err)
	assert.Equal(http.StatusCreated, res.StatusCode)

	var todo model.Todo
	err = json.NewDecoder(res.Body).Decode(&todo)
	assert.NoError(err)

	assert.Equal(todo.Name, "Test Todo")

	id1 := todo.ID

	res, err = http.PostForm(ts.URL+"/todos", url.Values{"name": {"Test Todo2"}})
	assert.NoError(err)
	assert.Equal(http.StatusCreated, res.StatusCode)

	err = json.NewDecoder(res.Body).Decode(&todo)
	assert.NoError(err)

	assert.Equal(todo.Name, "Test Todo2")
	id2 := todo.ID

	res, err = http.Get(ts.URL + "/todos")
	assert.NoError(err)
	assert.Equal(http.StatusOK, res.StatusCode)

	todos := []*model.Todo{}
	err = json.NewDecoder(res.Body).Decode(&todos)
	assert.NoError(err)
	assert.Equal(len(todos), 2)

	for _, t := range todos {
		if t.ID == id1 {
			assert.Equal("Test Todo", t.Name)
		} else if t.ID == id2 {
			assert.Equal("Test Todo2", t.Name)
		} else {
			assert.Error(fmt.Errorf("testID should e id1 or id2"))
		}
	}

	res, err = http.Get(ts.URL + "/complete-todo/" + strconv.Itoa(id1) + "?complete=true")
	assert.NoError(err)
	assert.Equal(http.StatusOK, res.StatusCode)

	res, err = http.Get(ts.URL + "/todos")
	assert.NoError(err)
	assert.Equal(http.StatusOK, res.StatusCode)

	todos = []*model.Todo{}
	err = json.NewDecoder(res.Body).Decode(&todos)
	assert.NoError(err)
	assert.Equal(len(todos), 2)

	for _, t := range todos {
		if t.ID == id1 {
			assert.True(t.Completed)
		}
	}

	req, _ := http.NewRequest("DELETE", ts.URL+"/todos/"+strconv.Itoa(id1), nil)
	res, err = http.DefaultClient.Do(req)
	assert.NoError(err)
	assert.Equal(http.StatusOK, res.StatusCode)

	res, err = http.Get(ts.URL + "/todos")
	assert.NoError(err)
	assert.Equal(http.StatusOK, res.StatusCode)

	todos = []*model.Todo{}
	err = json.NewDecoder(res.Body).Decode(&todos)
	assert.NoError(err)
	assert.Equal(len(todos), 1)

	for _, t := range todos {
		assert.Equal(t.ID, id2)
	}
}
