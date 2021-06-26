package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"testing"

	"github.com/speedmancs/vmmanager/app"
	"github.com/speedmancs/vmmanager/middleware"
	"github.com/speedmancs/vmmanager/model"
	"github.com/speedmancs/vmmanager/util"
	"github.com/stretchr/testify/assert"
)

var a app.App
var logFile string

func TestMain(m *testing.M) {
	logFile = "main_test_logs.txt"
	if _, err := os.Stat(logFile); err == nil {
		os.Remove(logFile)
	}
	a.Initialize(logFile)
	code := m.Run()
	os.Exit(code)
}

func executeRequestWithToken(req *http.Request, token string) *httptest.ResponseRecorder {
	var bearer = "Bearer " + token
	req.Header.Add("Authorization", bearer)
	return executeRequest(req)
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

func TestLogin(t *testing.T) {
	token := login(t)
	assert.NotEmpty(t, token)
}

func TestAuthFailure(t *testing.T) {
	req, _ := http.NewRequest("GET", "/vm/all", nil)
	response := executeRequest(req)
	assert.Equal(t, http.StatusUnauthorized, response.Code)
}

func login(t *testing.T) string {
	var jsonStr = []byte(`{"username":"harry", "password": "password"}`)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)
	assert.Equal(t, http.StatusOK, response.Code)
	var tokenObj middleware.Token
	err := json.NewDecoder(response.Body).Decode(&tokenObj)
	assert.Nil(t, err)
	return tokenObj.Token
}
func TestAuthGetAllVM(t *testing.T) {
	Setup()
	token := login(t)
	assert.NotEmpty(t, token)

	req, _ := http.NewRequest("GET", "/vm/all", nil)
	response := executeRequestWithToken(req, token)
	assert.Equal(t, http.StatusOK, response.Code)
	var vms []model.VM
	err := json.NewDecoder(response.Body).Decode(&vms)
	assert.Nil(t, err)
	assert.Equal(t, 3, len(vms))
	assert.Equal(t, "vm1", vms[1].Name)
	assert.Equal(t, "running", vms[1].Status)
}

func TestRequestID(t *testing.T) {
	token := login(t)
	assert.NotEmpty(t, token)
	req, _ := http.NewRequest("GET", "/vm/all", nil)
	response := executeRequestWithToken(req, token)

	assert.Equal(t, http.StatusOK, response.Code)
	assert.NotEmpty(t, response.Header().Get(middleware.ContextKeyRequestID))
}

func TestLogging(t *testing.T) {
	token := login(t)
	assert.NotEmpty(t, token)
	req, _ := http.NewRequest("GET", "/vm/all", nil)
	response := executeRequestWithToken(req, token)

	assert.Equal(t, http.StatusOK, response.Code)
	reqID := response.Header().Get(middleware.ContextKeyRequestID)
	assert.NotEmpty(t, reqID)

	lines := util.ReadAllLines(logFile)
	n := len(lines)
	assert.Contains(t, lines[n-4], fmt.Sprintf("Generated id %s", reqID))
	assert.Contains(t, lines[n-3], fmt.Sprintf("request %s started at", reqID))
	assert.Contains(t, lines[n-2], fmt.Sprintf("request %s url:/vm/all method:GET", reqID))

	r, _ := regexp.Compile(fmt.Sprintf("request %s duration: .+, with status code: 200", reqID))
	assert.True(t, r.MatchString(lines[n-1]))
}
