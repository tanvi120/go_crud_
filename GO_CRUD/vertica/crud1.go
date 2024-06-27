package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

const (
	baseURL             = "/clients"
	contentTypeHeader   = "Content-Type"
	applicationJSON     = "application/json"
	methodNotAllowedMsg = "Method not allowed"
	invalidClientIDMsg  = "Invalid client ID"
	clientNotFoundMsg   = "Client not found"
)

type Client struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var clients []Client

func handleGetClients(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, methodNotAllowedMsg, http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set(contentTypeHeader, applicationJSON)
	json.NewEncoder(w).Encode(clients)
}

func handleGetSingleClient(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, methodNotAllowedMsg, http.StatusMethodNotAllowed)
		return
	}
	idStr := strings.TrimPrefix(r.URL.Path, baseURL+"/")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, invalidClientIDMsg, http.StatusBadRequest)
		return
	}
	client, found := findClientByID(id)
	if !found {
		http.Error(w, clientNotFoundMsg, http.StatusNotFound)
		return
	}
	w.Header().Set(contentTypeHeader, applicationJSON)
	json.NewEncoder(w).Encode(client)
}

func handleCreateClient(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, methodNotAllowedMsg, http.StatusMethodNotAllowed)
		return
	}
	var client Client
	if err := json.NewDecoder(r.Body).Decode(&client); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	clients = append(clients, client)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(client)
}

func handleUpdateClient(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, methodNotAllowedMsg, http.StatusMethodNotAllowed)
		return
	}
	idStr := strings.TrimPrefix(r.URL.Path, baseURL+"/")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, invalidClientIDMsg, http.StatusBadRequest)
		return
	}
	var updatedClient Client
	if err := json.NewDecoder(r.Body).Decode(&updatedClient); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	client, found := findClientByID(id)
	if !found {
		http.Error(w, clientNotFoundMsg, http.StatusNotFound)
		return
	}
	client.Name = updatedClient.Name
	w.Header().Set(contentTypeHeader, applicationJSON)
	json.NewEncoder(w).Encode(client)
}

func handleDeleteClient(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, methodNotAllowedMsg, http.StatusMethodNotAllowed)
		return
	}
	idStr := strings.TrimPrefix(r.URL.Path, baseURL+"/")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, invalidClientIDMsg, http.StatusBadRequest)
		return
	}
	index, found := findClientIndexByID(id)
	if !found {
		http.Error(w, clientNotFoundMsg, http.StatusNotFound)
		return
	}
	clients = append(clients[:index], clients[index+1:]...)
	w.WriteHeader(http.StatusOK)
}

func findClientByID(id int) (Client, bool) {
	for _, client := range clients {
		if client.ID == id {
			return client, true
		}
	}
	return Client{}, false
}

func findClientIndexByID(id int) (int, bool) {
	for i, client := range clients {
		if client.ID == id {
			return i, true
		}
	}
	return -1, false
}
