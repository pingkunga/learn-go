//go:build integration

package user

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Response is a wrapper around http.Response that provides a way to check for
type Response struct {
	*http.Response
	err error
}

func clientRequest(method, url string, body io.Reader) *Response {
	req, _ := http.NewRequest(method, url, body)
	req.Header.Add("Authorization", os.Getenv("AUTH_TOKEN"))
	req.Header.Add("Content-Type", "application/json")
	client := http.Client{}
	res, err := client.Do(req)
	return &Response{res, err}
}

func (r *Response) Decode(v interface{}) error {
	if r.err != nil {
		return r.err
	}
	defer r.Body.Close()

	//NewDecoder == json.unmarshal
	//เอา r.Body มา decode แล้วเก็บใน v >>
	// - v เป็น Struct ในที่นี้ user
	return json.NewDecoder(r.Body).Decode(v)
}

func uri(paths ...string) string {
	baseURL := "http://localhost:10170"
	if paths == nil {
		return baseURL
	}
	return baseURL + "/" + strings.Join(paths, "/")
}

func TestGetAuthTokenFromEnv(t *testing.T) {
	//Act
	token := os.Getenv("AUTH_TOKEN")
	//Assert
	assert.NotEmpty(t, token)
}

func TestGetAllUser(t *testing.T) {
	//Arrange
	var us []User

	//Act
	res := clientRequest(http.MethodGet, uri("api/users"), nil)
	err := res.Decode(&us)

	//Assert
	assert.Nil(t, err)
	assert.EqualValues(t, http.StatusOK, res.StatusCode)
	//ดูว่ามีของคืนมาไหม ไม่สนใจว่า Value ถูกไหม
	assert.Greater(t, len(us), 0)
}

func TestCreateUser(t *testing.T) {
	body := bytes.NewBufferString(`{
		"name": "Chatri Ng",
		"age": 33
	}`)
	var u User

	res := clientRequest(http.MethodPost, uri("api/users"), body)
	err := res.Decode(&u)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, res.StatusCode)
	assert.NotEqual(t, 0, u.Id)
	assert.Equal(t, "Chatri Ng", u.Name)
	assert.Equal(t, 33, u.Age)
}

func TestGetUserById(t *testing.T) {
	userEntry := seedUser(t)

	var latestUser User
	res := clientRequest(http.MethodGet, uri("api/users", strconv.Itoa(userEntry.Id)), nil)
	err := res.Decode(&latestUser)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, userEntry.Id, latestUser.Id)
	assert.NotEmpty(t, latestUser.Name)
	assert.NotEmpty(t, latestUser.Age)
}

func TestUpdateUserById(t *testing.T) {
	//t.Skip("TODO: implement me")

	userEntry := seedUser(t)
	log.Println("userEntry:", userEntry)

	var latestUser User
	res := clientRequest(http.MethodGet, uri("api/users", strconv.Itoa(userEntry.Id)), nil)
	err := res.Decode(&latestUser)

	if err != nil {
		t.Fatal("can't get user:", err)
	}

	latestUser.Age = 30
	latestUser.Name = "Chatri Nga"

	latestUserJSON, _ := json.Marshal(latestUser) // Convert latestUser to JSON format

	//Act
	resPut := clientRequest(http.MethodPut, uri("api/users", strconv.Itoa(userEntry.Id)), bytes.NewBuffer(latestUserJSON)) // Use bytes.NewBuffer to create a reader from the JSON data
	var updatedUser User
	updateerr := resPut.Decode(&updatedUser)

	//Assert
	assert.Nil(t, updateerr)
	assert.Equal(t, http.StatusOK, resPut.StatusCode)
	assert.Equal(t, latestUser.Id, updatedUser.Id)
	assert.Equal(t, latestUser.Name, updatedUser.Name)
	assert.Equal(t, latestUser.Age, updatedUser.Age)

}

func TestDeleteUserById(t *testing.T) {
	//t.Skip("TODO: implement me")
	userEntry := seedUser(t)

	//Act
	res := clientRequest(http.MethodDelete, uri("api/users", strconv.Itoa(userEntry.Id)), nil)

	//Assert
	assert.Equal(t, http.StatusNoContent, res.StatusCode)
}

func seedUser(t *testing.T) User {
	var userEntry User
	body := bytes.NewBufferString(`{
		"name": "pingkungA",
		"age": 28
	}`)
	err := clientRequest(http.MethodPost, uri("api/users"), body).Decode(&userEntry)
	if err != nil {
		t.Fatal("can't create uomer:", err)
	}
	return userEntry
}
