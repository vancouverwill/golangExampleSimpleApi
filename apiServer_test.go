package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAPIGet(t *testing.T) {
	req, err := http.NewRequest("GET", "http://localhost:4000/v1/", nil)
	if err != nil {
		t.Error("GET did not work as expected.")
		t.Log(err)
	}

	w := httptest.NewRecorder()
	APIHandler(w, req)

	if w.Code != 200 && w.Code != 202 {
		t.Error("GET did not work as expected. the status was not ", http.StatusOK, ", it was ", w.Code)
	}

	t.Log("status:", w.Code, "body:", w.Body.String())
}

func TestAPIGetSearch(t *testing.T) {
	req, err := http.NewRequest("GET", "http://localhost:4000/v1/?userId=14", nil)
	if err != nil {
		t.Error("GET did not work as expected.")
		t.Log(err)
	}

	w := httptest.NewRecorder()
	APIHandler(w, req)

	if w.Code != 200 && w.Code != 202 {
		t.Error("GET did not work as expected. the status was not ", http.StatusOK, ", it was ", w.Code)
	}

	t.Log("status:", w.Code, "body:", w.Body.String())
}



func TestAPIPOST(t *testing.T) {
	var userRestObject UserRestObject
	userRestObject.Name = "test Name"
	userRestObject.Age = 50
	data, err := json.Marshal(userRestObject)
	if err != nil {
		t.Error("error:", err)
	}
	req, err := http.NewRequest("POST", "http://localhost:4000/v1/", bytes.NewBufferString(string(data)))
	if err != nil {
		t.Error("POST did not work as expected.")
		t.Log(err)
	}

	w := httptest.NewRecorder()
	APIHandler(w, req)

	if w.Code != 200 && w.Code != 202 {
		t.Error("POST did not work as expected. the status was not ", http.StatusOK, ", it was ", w.Code)
	}

	t.Log("status:", w.Code, "body:", w.Body.String())
}

func TestAPIPUT(t *testing.T) {
	var userRestObject UserRestObject
	userRestObject.Name = "test Name 3"
	userRestObject.Age = 50
	data, err := json.Marshal(userRestObject)
	if err != nil {
		t.Error("error:", err)
	}
	req, err := http.NewRequest("PUT", "http://localhost:4000/v1/14", bytes.NewBufferString(string(data)))
	if err != nil {
		t.Error("POST did not work as expected.")
		t.Log(err)
	}

	w := httptest.NewRecorder()
	APIHandler(w, req)

	if w.Code != 200 && w.Code != 202 {
		t.Error("POST did not work as expected. the status was not ", http.StatusOK, ", it was ", w.Code)
	}

	t.Log("status:", w.Code, "body:", w.Body.String())
}
