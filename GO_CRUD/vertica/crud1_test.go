package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	wrongStatusCodeMessage = "Handler returned wrong status code: got %v, want %v"
	clientPath             = "/clients/1"
	updatedClientName      = "Updated Client 1"
	nonExistingClientPath  = "/clients/99"
)

func TestHandleGetClients(t *testing.T) {
	req, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleGetClients)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf(wrongStatusCodeMessage, rr.Code, http.StatusOK)
	}

	var clients []Client
	err = json.Unmarshal(rr.Body.Bytes(), &clients)
	if err != nil {
		t.Errorf("Error decoding response body: %v", err)
	}

	if len(clients) != len(initialClients) {
		t.Errorf("Expected clients count: %d, got: %d", len(initialClients), len(clients))
	}
}

func TestHandleGetSingleClient(t *testing.T) {
	req, err := http.NewRequest("GET", clientPath, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleGetSingleClient)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf(wrongStatusCodeMessage, rr.Code, http.StatusOK)
	}

	var client Client
	err = json.Unmarshal(rr.Body.Bytes(), &client)
	if err != nil {
		t.Errorf("Error decoding response body: %v", err)
	}

	expectedClient := Client{ID: 1, Name: "Client 1"}
	if client != expectedClient {
		t.Errorf("Expected client %+v, got: %+v", expectedClient, client)
	}
}

func TestHandleGetNonExistingClient(t *testing.T) {
	req, err := http.NewRequest("GET", nonExistingClientPath, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleGetSingleClient)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf(wrongStatusCodeMessage, rr.Code, http.StatusNotFound)
	}
}

func TestHandleCreateClient(t *testing.T) {
	client := Client{ID: 4, Name: "Client 4"}
	clientJSON, err := json.Marshal(client)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", baseURL, bytes.NewBuffer(clientJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set(contentTypeHeader, applicationJSON)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleCreateClient)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf(wrongStatusCodeMessage, rr.Code, http.StatusCreated)
	}

	// Check if client was added
	if len(clients) != 4 {
		t.Errorf("Expected clients count: %d, got: %d", 4, len(clients))
	}
}

func TestHandleCreateClientInvalidJSON(t *testing.T) {
	invalidJSON := []byte(`{"ID": 4, "Name":`)

	req, err := http.NewRequest("POST", baseURL, bytes.NewBuffer(invalidJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set(contentTypeHeader, applicationJSON)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleCreateClient)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf(wrongStatusCodeMessage, rr.Code, http.StatusBadRequest)
	}
}

func TestHandleUpdateClient(t *testing.T) {
	client := Client{ID: 1, Name: updatedClientName}
	clientJSON, err := json.Marshal(client)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("PUT", clientPath, bytes.NewBuffer(clientJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set(contentTypeHeader, applicationJSON)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleUpdateClient)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf(wrongStatusCodeMessage, rr.Code, http.StatusOK)
	}

	// Check if client was updated
	updatedClient, _ := findClientByID(1)
	if updatedClient.Name != updatedClientName {
		t.Errorf("Expected client name: %s, got: %s", updatedClientName, updatedClient.Name)
	}
}

func TestHandleUpdateNonExistingClient(t *testing.T) {
	client := Client{ID: 99, Name: "Non-existing Client"}
	clientJSON, err := json.Marshal(client)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("PUT", nonExistingClientPath, bytes.NewBuffer(clientJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set(contentTypeHeader, applicationJSON)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleUpdateClient)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf(wrongStatusCodeMessage, rr.Code, http.StatusNotFound)
	}
}

func TestHandleDeleteClient(t *testing.T) {
	req, err := http.NewRequest("DELETE", clientPath, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleDeleteClient)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf(wrongStatusCodeMessage, rr.Code, http.StatusOK)
	}

	// Check if client was deleted
	_, found := findClientByID(1)
	if found {
		t.Error("Client was not deleted successfully")
	}
}

func TestHandleDeleteNonExistingClient(t *testing.T) {
	req, err := http.NewRequest("DELETE", nonExistingClientPath, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleDeleteClient)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf(wrongStatusCodeMessage, rr.Code, http.StatusNotFound)
	}
}

func TestHandleInvalidMethod(t *testing.T) {
	req, err := http.NewRequest("POST", clientPath, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleGetSingleClient)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf(wrongStatusCodeMessage, rr.Code, http.StatusMethodNotAllowed)
	}
}

var initialClients = []Client{
	{ID: 1, Name: "Client 1"},
	{ID: 2, Name: "Client 2"},
	{ID: 3, Name: "Client 3"},
}

func setup() {
	// Reset clients to initial state before each test
	clients = make([]Client, len(initialClients))
	copy(clients, initialClients)
}

func TestMain(m *testing.M) {
	setup()
	m.Run()
}
