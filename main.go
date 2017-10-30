package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"load_balancer/lobby"
	"log"
	"net/http"
)

// lobbies holds all the currently available lobbies.
var lobbies []lobby.Lobby

// GetLobby writes the list of lobbies.
func GetLobby(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(lobbies)
}

// main function bootstraps the load balancer.
func main() {
	lobby1 := lobby.Lobby{
		Name:       "KSSV Lobby",
		MaxPlayers: 4,
		NumPlayers: 0,
		Host:       "pylon1.usc.edu",
		Port:       8888,
	}
	lobbies = append(lobbies, lobby1)

	router := mux.NewRouter()
	router.HandleFunc("/lobby", GetLobby).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", router))
}
