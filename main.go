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
	lobby1 := lobby.NewLobby("KSSV Lobby", 4)
	lobby1.NumPlayers = 0
	lobby1.Server = "pylon1.usc.edu:7777"
	lobbies = append(lobbies, lobby1)

	router := mux.NewRouter()
	router.HandleFunc("/lobby", GetLobby).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", router))
}
